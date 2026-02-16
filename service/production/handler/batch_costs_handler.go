package handler

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

// BatchCostsStore defines the storage operations needed by the batch costs handler.
type BatchCostsStore interface {
	GetBatchByUUID(context.Context, string) (storage.Batch, error)
	GetBatchSummaryByUUID(context.Context, string) (storage.BatchSummary, error)
	ListAdditionsByBatchUUID(context.Context, string) ([]storage.Addition, error)
}

// IngredientLotFetcher abstracts the inter-service call to the Inventory service
// for retrieving ingredient lot data consumed by a batch.
type IngredientLotFetcher interface {
	GetBatchIngredientLots(ctx context.Context, authToken string, batchUUID string) ([]BatchIngredientLot, error)
}

// POLineFetcher abstracts the inter-service call to the Procurement service
// for looking up purchase order line costs.
type POLineFetcher interface {
	BatchLookupPOLines(ctx context.Context, authToken string, uuids []string) ([]PurchaseOrderLineCost, error)
}

// HandleBatchCosts handles [GET /batches/{uuid}/costs].
func HandleBatchCosts(db BatchCostsStore, invClient IngredientLotFetcher, procClient POLineFetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		batchUUID := r.PathValue("uuid")
		if batchUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		ctx := r.Context()

		// 1. Verify batch exists.
		_, err := db.GetBatchByUUID(ctx, batchUUID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "batch not found", http.StatusNotFound)
			return
		} else if err != nil {
			service.InternalError(w, "error getting batch", "error", err, "batch_uuid", batchUUID)
			return
		}

		// 2. Get batch summary for volume metrics.
		summary, err := db.GetBatchSummaryByUUID(ctx, batchUUID)
		if err != nil {
			service.InternalError(w, "error getting batch summary", "error", err, "batch_uuid", batchUUID)
			return
		}

		// 3. List all additions for the batch.
		additions, err := db.ListAdditionsByBatchUUID(ctx, batchUUID)
		if err != nil {
			service.InternalError(w, "error listing additions", "error", err, "batch_uuid", batchUUID)
			return
		}

		// 4. Early return if no additions.
		if len(additions) == 0 {
			service.JSON(w, dto.BatchCostsResponse{
				BatchUUID:         batchUUID,
				LineItems:         []dto.CostLineItem{},
				UncostedAdditions: []dto.UncostedAddition{},
				Totals: dto.CostTotals{
					CostComplete: true,
				},
			})
			return
		}

		// 5. Extract auth token for inter-service calls.
		authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		// 6. Fetch ingredient lot data from Inventory service.
		lots, err := invClient.GetBatchIngredientLots(ctx, authToken, batchUUID)
		if err != nil {
			service.InternalError(w, "error fetching ingredient lot data", "error", err, "batch_uuid", batchUUID)
			return
		}

		// 7. Build lot lookup map: lotUUID → *BatchIngredientLot.
		lotMap := make(map[string]*BatchIngredientLot, len(lots))
		for i := range lots {
			lotMap[lots[i].IngredientLotUUID] = &lots[i]
		}

		// 8. Collect unique PO line UUIDs from lot data.
		poLineUUIDSet := make(map[string]struct{})
		for _, lot := range lots {
			if lot.PurchaseOrderLineUUID != nil {
				poLineUUIDSet[*lot.PurchaseOrderLineUUID] = struct{}{}
			}
		}

		// 9. Fetch PO line costs from Procurement service (if any).
		poLineMap := make(map[string]*PurchaseOrderLineCost)
		if len(poLineUUIDSet) > 0 {
			poLineUUIDs := make([]string, 0, len(poLineUUIDSet))
			for uuid := range poLineUUIDSet {
				poLineUUIDs = append(poLineUUIDs, uuid)
			}

			poLines, err := procClient.BatchLookupPOLines(ctx, authToken, poLineUUIDs)
			if err != nil {
				service.InternalError(w, "error fetching purchase order line data", "error", err, "batch_uuid", batchUUID)
				return
			}

			for i := range poLines {
				poLineMap[poLines[i].UUID] = &poLines[i]
			}
		}

		// 10. Build response line items and uncosted additions.
		lineItems := make([]dto.CostLineItem, 0, len(additions))
		uncostedAdditions := make([]dto.UncostedAddition, 0)

		for _, addition := range additions {
			additionUUID := addition.UUID.String()

			// No inventory lot → uncosted.
			if addition.InventoryLotUUID == nil {
				uncostedAdditions = append(uncostedAdditions, dto.UncostedAddition{
					AdditionUUID: additionUUID,
					AdditionType: addition.AdditionType,
					AmountUsed:   addition.Amount,
					AmountUnit:   addition.AmountUnit,
					Reason:       "no_inventory_lot",
				})
				continue
			}

			lotUUID := addition.InventoryLotUUID.String()
			lot := lotMap[lotUUID]

			item := dto.CostLineItem{
				AdditionUUID:      additionUUID,
				IngredientLotUUID: lotUUID,
				AdditionType:      addition.AdditionType,
				AmountUsed:        addition.Amount,
				AmountUnit:        addition.AmountUnit,
				CostSource:        "unavailable",
			}

			if lot != nil {
				item.IngredientUUID = &lot.IngredientUUID
				item.IngredientName = &lot.IngredientName
				item.IngredientCategory = &lot.IngredientCategory
				item.LotCode = lot.BreweryLotCode

				if lot.PurchaseOrderLineUUID != nil {
					poLine := poLineMap[*lot.PurchaseOrderLineUUID]
					if poLine != nil && addition.AmountUnit == poLine.QuantityUnit {
						costCents := addition.Amount * poLine.UnitCostCents
						item.CostCents = &costCents
						item.UnitCostCents = &poLine.UnitCostCents
						item.UnitCostUnit = &poLine.QuantityUnit
						item.PurchaseOrderLineUUID = lot.PurchaseOrderLineUUID
						item.CostSource = "purchase_order"
					}
				}
			}

			lineItems = append(lineItems, item)
		}

		// 11. Compute totals.
		var totalCostCents int64
		var costedLineCount int
		var uncostedLineCount int

		for _, item := range lineItems {
			if item.CostCents != nil {
				totalCostCents += *item.CostCents
				costedLineCount++
			} else {
				uncostedLineCount++
			}
		}
		uncostedLineCount += len(uncostedAdditions)

		costComplete := uncostedLineCount == 0

		// 12. Determine currency from costed PO lines.
		var currency *string
		currencies := make(map[string]struct{})
		for _, item := range lineItems {
			if item.CostCents != nil && item.PurchaseOrderLineUUID != nil {
				poLine := poLineMap[*item.PurchaseOrderLineUUID]
				if poLine != nil {
					currencies[poLine.Currency] = struct{}{}
				}
			}
		}
		if len(currencies) == 1 {
			for c := range currencies {
				currency = &c
			}
		} else if len(currencies) > 1 {
			mixed := "MIXED"
			currency = &mixed
		}

		// 13. Compute batch volume in BBL from the starting volume.
		var batchVolumeBBL *float64
		var startingVolume *storage.BatchVolumeWithAmount
		for i := range summary.Volumes {
			v := &summary.Volumes[i]
			if startingVolume == nil || v.BatchVolume.PhaseAt.Before(startingVolume.BatchVolume.PhaseAt) {
				startingVolume = v
			}
		}
		if startingVolume != nil {
			batchVolumeBBL = dto.ConvertToBBL(startingVolume.Volume.Amount, startingVolume.Volume.AmountUnit)
		}

		// 14. Compute cost per BBL.
		var costPerBBLCents *int64
		if batchVolumeBBL != nil && *batchVolumeBBL > 0 && (currency == nil || *currency != "MIXED") {
			v := int64(math.Round(float64(totalCostCents) / *batchVolumeBBL))
			costPerBBLCents = &v
		}

		resp := dto.BatchCostsResponse{
			BatchUUID:         batchUUID,
			Currency:          currency,
			LineItems:         lineItems,
			UncostedAdditions: uncostedAdditions,
			Totals: dto.CostTotals{
				TotalCostCents:    totalCostCents,
				CostedLineCount:   costedLineCount,
				UncostedLineCount: uncostedLineCount,
				CostComplete:      costComplete,
				CostPerBBLCents:   costPerBBLCents,
				BatchVolumeBBL:    batchVolumeBBL,
			},
		}

		service.JSON(w, resp)
	}
}
