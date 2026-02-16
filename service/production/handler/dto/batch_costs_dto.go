package dto

// BatchCostsResponse is the response for GET /batches/{uuid}/costs.
type BatchCostsResponse struct {
	BatchUUID         string             `json:"batch_uuid"`
	Currency          *string            `json:"currency,omitempty"`
	LineItems         []CostLineItem     `json:"line_items"`
	UncostedAdditions []UncostedAddition `json:"uncosted_additions"`
	Totals            CostTotals         `json:"totals"`
}

// CostLineItem represents a single costed ingredient addition.
type CostLineItem struct {
	AdditionUUID          string  `json:"addition_uuid"`
	IngredientLotUUID     string  `json:"ingredient_lot_uuid"`
	IngredientUUID        *string `json:"ingredient_uuid,omitempty"`
	IngredientName        *string `json:"ingredient_name,omitempty"`
	IngredientCategory    *string `json:"ingredient_category,omitempty"`
	LotCode               *string `json:"lot_code,omitempty"`
	AdditionType          string  `json:"addition_type"`
	AmountUsed            int64   `json:"amount_used"`
	AmountUnit            string  `json:"amount_unit"`
	UnitCostCents         *int64  `json:"unit_cost_cents,omitempty"`
	UnitCostUnit          *string `json:"unit_cost_unit,omitempty"`
	CostCents             *int64  `json:"cost_cents,omitempty"`
	CostSource            string  `json:"cost_source"`
	PurchaseOrderLineUUID *string `json:"purchase_order_line_uuid,omitempty"`
}

// UncostedAddition represents an addition that cannot be costed.
type UncostedAddition struct {
	AdditionUUID string `json:"addition_uuid"`
	AdditionType string `json:"addition_type"`
	AmountUsed   int64  `json:"amount_used"`
	AmountUnit   string `json:"amount_unit"`
	Reason       string `json:"reason"`
}

// CostTotals holds aggregated cost metrics for a batch.
type CostTotals struct {
	TotalCostCents    int64    `json:"total_cost_cents"`
	CostedLineCount   int      `json:"costed_line_count"`
	UncostedLineCount int      `json:"uncosted_line_count"`
	CostComplete      bool     `json:"cost_complete"`
	CostPerBBLCents   *int64   `json:"cost_per_bbl_cents,omitempty"`
	BatchVolumeBBL    *float64 `json:"batch_volume_bbl,omitempty"`
}
