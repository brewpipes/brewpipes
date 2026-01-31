package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateIngredientMaltDetail(ctx context.Context, detail IngredientMaltDetail) (IngredientMaltDetail, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO ingredient_malt_detail (
			ingredient_id,
			maltster_name,
			variety,
			lovibond,
			srm,
			diastatic_power
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, uuid, ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power, created_at, updated_at, deleted_at`,
		detail.IngredientID,
		detail.MaltsterName,
		detail.Variety,
		detail.Lovibond,
		detail.SRM,
		detail.DiastaticPower,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.MaltsterName,
		&detail.Variety,
		&detail.Lovibond,
		&detail.SRM,
		&detail.DiastaticPower,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientMaltDetail{}, fmt.Errorf("creating ingredient malt detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientMaltDetail(ctx context.Context, id int64) (IngredientMaltDetail, error) {
	var detail IngredientMaltDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power, created_at, updated_at, deleted_at
		FROM ingredient_malt_detail
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.MaltsterName,
		&detail.Variety,
		&detail.Lovibond,
		&detail.SRM,
		&detail.DiastaticPower,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientMaltDetail{}, service.ErrNotFound
		}
		return IngredientMaltDetail{}, fmt.Errorf("getting ingredient malt detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientMaltDetailByIngredient(ctx context.Context, ingredientID int64) (IngredientMaltDetail, error) {
	var detail IngredientMaltDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, maltster_name, variety, lovibond, srm, diastatic_power, created_at, updated_at, deleted_at
		FROM ingredient_malt_detail
		WHERE ingredient_id = $1 AND deleted_at IS NULL`,
		ingredientID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.MaltsterName,
		&detail.Variety,
		&detail.Lovibond,
		&detail.SRM,
		&detail.DiastaticPower,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientMaltDetail{}, service.ErrNotFound
		}
		return IngredientMaltDetail{}, fmt.Errorf("getting ingredient malt detail by ingredient: %w", err)
	}

	return detail, nil
}

func (c *Client) CreateIngredientHopDetail(ctx context.Context, detail IngredientHopDetail) (IngredientHopDetail, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO ingredient_hop_detail (
			ingredient_id,
			producer_name,
			variety,
			crop_year,
			form,
			alpha_acid,
			beta_acid
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, uuid, ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid, created_at, updated_at, deleted_at`,
		detail.IngredientID,
		detail.ProducerName,
		detail.Variety,
		detail.CropYear,
		detail.Form,
		detail.AlphaAcid,
		detail.BetaAcid,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.ProducerName,
		&detail.Variety,
		&detail.CropYear,
		&detail.Form,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientHopDetail{}, fmt.Errorf("creating ingredient hop detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientHopDetail(ctx context.Context, id int64) (IngredientHopDetail, error) {
	var detail IngredientHopDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid, created_at, updated_at, deleted_at
		FROM ingredient_hop_detail
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.ProducerName,
		&detail.Variety,
		&detail.CropYear,
		&detail.Form,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientHopDetail{}, service.ErrNotFound
		}
		return IngredientHopDetail{}, fmt.Errorf("getting ingredient hop detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientHopDetailByIngredient(ctx context.Context, ingredientID int64) (IngredientHopDetail, error) {
	var detail IngredientHopDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, producer_name, variety, crop_year, form, alpha_acid, beta_acid, created_at, updated_at, deleted_at
		FROM ingredient_hop_detail
		WHERE ingredient_id = $1 AND deleted_at IS NULL`,
		ingredientID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.ProducerName,
		&detail.Variety,
		&detail.CropYear,
		&detail.Form,
		&detail.AlphaAcid,
		&detail.BetaAcid,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientHopDetail{}, service.ErrNotFound
		}
		return IngredientHopDetail{}, fmt.Errorf("getting ingredient hop detail by ingredient: %w", err)
	}

	return detail, nil
}

func (c *Client) CreateIngredientYeastDetail(ctx context.Context, detail IngredientYeastDetail) (IngredientYeastDetail, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO ingredient_yeast_detail (
			ingredient_id,
			lab_name,
			strain,
			form
		) VALUES ($1, $2, $3, $4)
		RETURNING id, uuid, ingredient_id, lab_name, strain, form, created_at, updated_at, deleted_at`,
		detail.IngredientID,
		detail.LabName,
		detail.Strain,
		detail.Form,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.LabName,
		&detail.Strain,
		&detail.Form,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		return IngredientYeastDetail{}, fmt.Errorf("creating ingredient yeast detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientYeastDetail(ctx context.Context, id int64) (IngredientYeastDetail, error) {
	var detail IngredientYeastDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, lab_name, strain, form, created_at, updated_at, deleted_at
		FROM ingredient_yeast_detail
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.LabName,
		&detail.Strain,
		&detail.Form,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientYeastDetail{}, service.ErrNotFound
		}
		return IngredientYeastDetail{}, fmt.Errorf("getting ingredient yeast detail: %w", err)
	}

	return detail, nil
}

func (c *Client) GetIngredientYeastDetailByIngredient(ctx context.Context, ingredientID int64) (IngredientYeastDetail, error) {
	var detail IngredientYeastDetail
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, ingredient_id, lab_name, strain, form, created_at, updated_at, deleted_at
		FROM ingredient_yeast_detail
		WHERE ingredient_id = $1 AND deleted_at IS NULL`,
		ingredientID,
	).Scan(
		&detail.ID,
		&detail.UUID,
		&detail.IngredientID,
		&detail.LabName,
		&detail.Strain,
		&detail.Form,
		&detail.CreatedAt,
		&detail.UpdatedAt,
		&detail.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return IngredientYeastDetail{}, service.ErrNotFound
		}
		return IngredientYeastDetail{}, fmt.Errorf("getting ingredient yeast detail by ingredient: %w", err)
	}

	return detail, nil
}
