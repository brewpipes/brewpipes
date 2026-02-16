package dto

import (
	"fmt"
	"time"

	"github.com/brewpipes/brewpipes/internal/validate"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type CreatePackagingRunLineRequest struct {
	PackageFormatUUID string `json:"package_format_uuid"`
	Quantity          int    `json:"quantity"`
}

type CreatePackagingRunRequest struct {
	BatchUUID         string                          `json:"batch_uuid"`
	OccupancyUUID     string                          `json:"occupancy_uuid"`
	StartedAt         *time.Time                      `json:"started_at"`
	EndedAt           *time.Time                      `json:"ended_at"`
	LossAmount        *int64                          `json:"loss_amount"`
	LossUnit          *string                         `json:"loss_unit"`
	Notes             *string                         `json:"notes"`
	Lines             []CreatePackagingRunLineRequest `json:"lines"`
	CloseSource       *bool                           `json:"close_source"`
	StockLocationUUID *string                         `json:"stock_location_uuid"`
	LotCodePrefix     *string                         `json:"lot_code_prefix"`
}

func (r CreatePackagingRunRequest) Validate() error {
	if err := validate.Required(r.BatchUUID, "batch_uuid"); err != nil {
		return err
	}
	if err := validate.Required(r.OccupancyUUID, "occupancy_uuid"); err != nil {
		return err
	}
	if len(r.Lines) == 0 {
		return fmt.Errorf("lines is required and must contain at least 1 line")
	}
	for i, line := range r.Lines {
		if err := validate.Required(line.PackageFormatUUID, fmt.Sprintf("lines[%d].package_format_uuid", i)); err != nil {
			return err
		}
		if line.Quantity <= 0 {
			return fmt.Errorf("lines[%d].quantity must be greater than zero", i)
		}
	}
	if (r.LossAmount == nil) != (r.LossUnit == nil) {
		return fmt.Errorf("loss_amount and loss_unit must be provided together")
	}
	if r.LossAmount != nil && *r.LossAmount <= 0 {
		return fmt.Errorf("loss_amount must be greater than zero")
	}
	if r.LossUnit != nil {
		if err := validateVolumeUnit(*r.LossUnit); err != nil {
			return err
		}
	}
	if r.StartedAt != nil && r.EndedAt != nil {
		if r.EndedAt.Before(*r.StartedAt) {
			return fmt.Errorf("ended_at must be after started_at")
		}
	}

	return nil
}

type PackagingRunLineResponse struct {
	UUID                           string    `json:"uuid"`
	PackagingRunUUID               string    `json:"packaging_run_uuid"`
	PackageFormatUUID              string    `json:"package_format_uuid"`
	PackageFormatName              string    `json:"package_format_name"`
	PackageFormatVolumePerUnit     int64     `json:"package_format_volume_per_unit"`
	PackageFormatVolumePerUnitUnit string    `json:"package_format_volume_per_unit_unit"`
	Quantity                       int       `json:"quantity"`
	CreatedAt                      time.Time `json:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at"`
}

type PackagingRunResponse struct {
	UUID          string                     `json:"uuid"`
	BatchUUID     string                     `json:"batch_uuid"`
	OccupancyUUID string                     `json:"occupancy_uuid"`
	StartedAt     time.Time                  `json:"started_at"`
	EndedAt       *time.Time                 `json:"ended_at,omitempty"`
	LossAmount    *int64                     `json:"loss_amount,omitempty"`
	LossUnit      *string                    `json:"loss_unit,omitempty"`
	Notes         *string                    `json:"notes,omitempty"`
	Lines         []PackagingRunLineResponse `json:"lines"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
}

func NewPackagingRunLineResponse(line storage.PackagingRunLine) PackagingRunLineResponse {
	return PackagingRunLineResponse{
		UUID:                           line.UUID.String(),
		PackagingRunUUID:               line.PackagingRunUUID,
		PackageFormatUUID:              line.PackageFormatUUID,
		PackageFormatName:              line.PackageFormatName,
		PackageFormatVolumePerUnit:     line.PackageFormatVolumePerUnit,
		PackageFormatVolumePerUnitUnit: line.PackageFormatVolumePerUnitUnit,
		Quantity:                       line.Quantity,
		CreatedAt:                      line.CreatedAt,
		UpdatedAt:                      line.UpdatedAt,
	}
}

func NewPackagingRunLinesResponse(lines []storage.PackagingRunLine) []PackagingRunLineResponse {
	resp := make([]PackagingRunLineResponse, 0, len(lines))
	for _, line := range lines {
		resp = append(resp, NewPackagingRunLineResponse(line))
	}
	return resp
}

func NewPackagingRunResponse(run storage.PackagingRun, lines []storage.PackagingRunLine) PackagingRunResponse {
	return PackagingRunResponse{
		UUID:          run.UUID.String(),
		BatchUUID:     run.BatchUUID,
		OccupancyUUID: run.OccupancyUUID,
		StartedAt:     run.StartedAt,
		EndedAt:       run.EndedAt,
		LossAmount:    run.LossAmount,
		LossUnit:      run.LossUnit,
		Notes:         run.Notes,
		Lines:         NewPackagingRunLinesResponse(lines),
		CreatedAt:     run.CreatedAt,
		UpdatedAt:     run.UpdatedAt,
	}
}

func NewPackagingRunsResponse(runs []storage.PackagingRun, linesByRunID map[int64][]storage.PackagingRunLine) []PackagingRunResponse {
	resp := make([]PackagingRunResponse, 0, len(runs))
	for _, run := range runs {
		resp = append(resp, NewPackagingRunResponse(run, linesByRunID[run.ID]))
	}
	return resp
}
