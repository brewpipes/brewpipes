package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateIngredientLotMaltDetail(ctx context.Context, detail IngredientLotMaltDetail) (IngredientLotMaltDetail, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO ingredient_lot_malt_detail (
			ingredient_lot_id,
			moisture_percent
		) VALUES ($1, $2)
		RETURNING id, uuid, ingredient_lot_id, moisture_percent, created_at, updated_at, deleted_at`,
		detail.IngredientLotID,
		detail.MoisturePercent,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.MoisturePercent,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientLotMaltDetail{}, fmt.Errorf("creating ingredient lot malt detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) GetIngredientLotMaltDetail(ctx context.Context, id int64) (IngredientLotMaltDetail, error) {
	var detail IngredientLotMaltDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.moisture_percent, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_malt_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.id = $1 AND d.deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.MoisturePercent,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotMaltDetail{}, service.ErrNotFound
		}
		return IngredientLotMaltDetail{}, fmt.Errorf("getting ingredient lot malt detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientLotMaltDetailByUUID(ctx context.Context, detailUUID string) (IngredientLotMaltDetail, error) {
	var detail IngredientLotMaltDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.moisture_percent, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_malt_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.uuid = $1 AND d.deleted_at IS NULL`,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.MoisturePercent,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotMaltDetail{}, service.ErrNotFound
		}
		return IngredientLotMaltDetail{}, fmt.Errorf("getting ingredient lot malt detail by uuid: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientLotMaltDetailByLot(ctx context.Context, lotUUID string) (IngredientLotMaltDetail, error) {
	var detail IngredientLotMaltDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.moisture_percent, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_malt_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE il.uuid = $1 AND d.deleted_at IS NULL`,
		lotUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.MoisturePercent,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotMaltDetail{}, service.ErrNotFound
		}
		return IngredientLotMaltDetail{}, fmt.Errorf("getting ingredient lot malt detail by lot: %w", err)
	}

	return detail, nil
}

func (c *Client) UpdateIngredientLotMaltDetail(ctx context.Context, detailUUID string, detail IngredientLotMaltDetail) (IngredientLotMaltDetail, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE ingredient_lot_malt_detail
		SET moisture_percent = $1, updated_at = timezone('utc', now())
		WHERE uuid = $2 AND deleted_at IS NULL
		RETURNING id, uuid, ingredient_lot_id, moisture_percent, created_at, updated_at, deleted_at`,
		detail.MoisturePercent,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.MoisturePercent,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotMaltDetail{}, service.ErrNotFound
		}
		return IngredientLotMaltDetail{}, fmt.Errorf("updating ingredient lot malt detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) CreateIngredientLotHopDetail(ctx context.Context, detail IngredientLotHopDetail) (IngredientLotHopDetail, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO ingredient_lot_hop_detail (
			ingredient_lot_id,
			alpha_acid,
			beta_acid
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, ingredient_lot_id, alpha_acid, beta_acid, created_at, updated_at, deleted_at`,
		detail.IngredientLotID,
		detail.AlphaAcid,
		detail.BetaAcid,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientLotHopDetail{}, fmt.Errorf("creating ingredient lot hop detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) GetIngredientLotHopDetail(ctx context.Context, id int64) (IngredientLotHopDetail, error) {
	var detail IngredientLotHopDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.alpha_acid, d.beta_acid, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_hop_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.id = $1 AND d.deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotHopDetail{}, service.ErrNotFound
		}
		return IngredientLotHopDetail{}, fmt.Errorf("getting ingredient lot hop detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientLotHopDetailByUUID(ctx context.Context, detailUUID string) (IngredientLotHopDetail, error) {
	var detail IngredientLotHopDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.alpha_acid, d.beta_acid, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_hop_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.uuid = $1 AND d.deleted_at IS NULL`,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotHopDetail{}, service.ErrNotFound
		}
		return IngredientLotHopDetail{}, fmt.Errorf("getting ingredient lot hop detail by uuid: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientLotHopDetailByLot(ctx context.Context, lotUUID string) (IngredientLotHopDetail, error) {
	var detail IngredientLotHopDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.alpha_acid, d.beta_acid, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_hop_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE il.uuid = $1 AND d.deleted_at IS NULL`,
		lotUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotHopDetail{}, service.ErrNotFound
		}
		return IngredientLotHopDetail{}, fmt.Errorf("getting ingredient lot hop detail by lot: %w", err)
	}

	return detail, nil
}

func (c *Client) UpdateIngredientLotHopDetail(ctx context.Context, detailUUID string, detail IngredientLotHopDetail) (IngredientLotHopDetail, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE ingredient_lot_hop_detail
		SET alpha_acid = $1, beta_acid = $2, updated_at = timezone('utc', now())
		WHERE uuid = $3 AND deleted_at IS NULL
		RETURNING id, uuid, ingredient_lot_id, alpha_acid, beta_acid, created_at, updated_at, deleted_at`,
		detail.AlphaAcid,
		detail.BetaAcid,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotHopDetail{}, service.ErrNotFound
		}
		return IngredientLotHopDetail{}, fmt.Errorf("updating ingredient lot hop detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) CreateIngredientLotYeastDetail(ctx context.Context, detail IngredientLotYeastDetail) (IngredientLotYeastDetail, error) {
	err := c.DB().QueryRow(ctx, `
		INSERT INTO ingredient_lot_yeast_detail (
			ingredient_lot_id,
			viability_percent,
			generation
		) VALUES ($1, $2, $3)
		RETURNING id, uuid, ingredient_lot_id, viability_percent, generation, created_at, updated_at, deleted_at`,
		detail.IngredientLotID,
		detail.ViabilityPercent,
		detail.Generation,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.ViabilityPercent,
		&detail.Generation,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientLotYeastDetail{}, fmt.Errorf("creating ingredient lot yeast detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) GetIngredientLotYeastDetail(ctx context.Context, id int64) (IngredientLotYeastDetail, error) {
	var detail IngredientLotYeastDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.viability_percent, d.generation, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_yeast_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.id = $1 AND d.deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.ViabilityPercent,
		&detail.Generation,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotYeastDetail{}, service.ErrNotFound
		}
		return IngredientLotYeastDetail{}, fmt.Errorf("getting ingredient lot yeast detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientLotYeastDetailByUUID(ctx context.Context, detailUUID string) (IngredientLotYeastDetail, error) {
	var detail IngredientLotYeastDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.viability_percent, d.generation, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_yeast_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE d.uuid = $1 AND d.deleted_at IS NULL`,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.ViabilityPercent,
		&detail.Generation,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotYeastDetail{}, service.ErrNotFound
		}
		return IngredientLotYeastDetail{}, fmt.Errorf("getting ingredient lot yeast detail by uuid: %w", err)
	}

	return detail, nil
}

func (c *Client) UpdateIngredientLotYeastDetail(ctx context.Context, detailUUID string, detail IngredientLotYeastDetail) (IngredientLotYeastDetail, error) {
	err := c.DB().QueryRow(ctx, `
		UPDATE ingredient_lot_yeast_detail
		SET viability_percent = $1, generation = $2, updated_at = timezone('utc', now())
		WHERE uuid = $3 AND deleted_at IS NULL
		RETURNING id, uuid, ingredient_lot_id, viability_percent, generation, created_at, updated_at, deleted_at`,
		detail.ViabilityPercent,
		detail.Generation,
		detailUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.ViabilityPercent,
		&detail.Generation,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotYeastDetail{}, service.ErrNotFound
		}
		return IngredientLotYeastDetail{}, fmt.Errorf("updating ingredient lot yeast detail: %w", err)
	}

	// Resolve ingredient lot UUID
	c.DB().QueryRow(ctx, `SELECT uuid FROM ingredient_lot WHERE id = $1`, detail.IngredientLotID).Scan(&detail.IngredientLotUUID)

	return detail, nil
}

func (c *Client) GetIngredientLotYeastDetailByLot(ctx context.Context, lotUUID string) (IngredientLotYeastDetail, error) {
	var detail IngredientLotYeastDetail
	err := c.DB().QueryRow(ctx, `
		SELECT d.id, d.uuid, d.ingredient_lot_id, il.uuid, d.viability_percent, d.generation, d.created_at, d.updated_at, d.deleted_at
		FROM ingredient_lot_yeast_detail d
		JOIN ingredient_lot il ON il.id = d.ingredient_lot_id
		WHERE il.uuid = $1 AND d.deleted_at IS NULL`,
		lotUUID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientLotID,
		&detail.IngredientLotUUID,
		&detail.ViabilityPercent,
		&detail.Generation,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientLotYeastDetail{}, service.ErrNotFound
		}
		return IngredientLotYeastDetail{}, fmt.Errorf("getting ingredient lot yeast detail by lot: %w", err)
	}

	return detail, nil
}
