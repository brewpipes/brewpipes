package storage

import (
	"context"
	"fmt"
)

type Volume struct {
	ID          int
	Name        string
	Description string
	Amount      int64
	AmountUnit  string
}

func (v Volume) Validate() error {
	if v.AmountUnit != "ml" && v.AmountUnit != "usfloz" && v.AmountUnit != "ukfloz" {
		return fmt.Errorf("invalid unit: %s", v.AmountUnit)
	}

	return nil
}

func (c *Client) GetVolumes(ctx context.Context) ([]Volume, error) {
	rows, err := c.DB.Query(ctx, `
		SELECT id, name, description, amount, amount_unit
		FROM volume
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []Volume
	for rows.Next() {
		var v Volume
		if err := rows.Scan(&v.ID, &v.Name, &v.Description, &v.Amount, &v.AmountUnit); err != nil {
			return nil, err
		}
		volumes = append(volumes, v)
	}

	return volumes, nil
}

func (c *Client) CreateVolume(ctx context.Context, volume Volume) (Volume, error) {
	var id int
	err := c.DB.QueryRow(ctx, `
		INSERT INTO volume (
			name,
			description,
			amount,
			amount_unit
		) VALUES ($1, $2, $3, $4) RETURNING id`,
		volume.Name, volume.Description, volume.Amount, volume.AmountUnit).Scan(&id)
	if err != nil {
		return Volume{}, err
	}

	volume.ID = id
	return volume, nil
}
