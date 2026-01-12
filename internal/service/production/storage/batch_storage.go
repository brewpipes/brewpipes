package storage

import (
	"context"
)

type Batch struct {
	ID          int
	Name        string
	Description string
}

func (v Batch) Validate() error {
	return nil
}

func (c *Client) CreateBatch(ctx context.Context, batch Batch) (Batch, error) {
	var id int
	err := c.DB.QueryRow(ctx, `
		INSERT INTO batch (
			name,
			description
		) VALUES ($1, $2) RETURNING id`,
		batch.Name,
		batch.Description,
	).Scan(&id)
	if err != nil {
		return Batch{}, err
	}

	batch.ID = id
	return batch, nil
}
