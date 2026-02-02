package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/brewpipes/brewpipes/service"
	"github.com/jackc/pgx/v5"
)

func (c *Client) CreateSupplier(ctx context.Context, supplier Supplier) (Supplier, error) {
	err := c.db.QueryRow(ctx, `
		INSERT INTO supplier (
			name,
			contact_name,
			email,
			phone,
			address_line1,
			address_line2,
			city,
			region,
			postal_code,
			country
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, uuid, name, contact_name, email, phone, address_line1, address_line2, city, region, postal_code, country, created_at, updated_at, deleted_at`,
		supplier.Name,
		supplier.ContactName,
		supplier.Email,
		supplier.Phone,
		supplier.AddressLine1,
		supplier.AddressLine2,
		supplier.City,
		supplier.Region,
		supplier.PostalCode,
		supplier.Country,
	).Scan(
		&supplier.ID,
		&supplier.UUID,
		&supplier.Name,
		&supplier.ContactName,
		&supplier.Email,
		&supplier.Phone,
		&supplier.AddressLine1,
		&supplier.AddressLine2,
		&supplier.City,
		&supplier.Region,
		&supplier.PostalCode,
		&supplier.Country,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.DeletedAt,
	)
	if err != nil {
		return Supplier{}, fmt.Errorf("creating supplier: %w", err)
	}

	return supplier, nil
}

func (c *Client) GetSupplier(ctx context.Context, id int64) (Supplier, error) {
	var supplier Supplier
	err := c.db.QueryRow(ctx, `
		SELECT id, uuid, name, contact_name, email, phone, address_line1, address_line2, city, region, postal_code, country, created_at, updated_at, deleted_at
		FROM supplier
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&supplier.ID,
		&supplier.UUID,
		&supplier.Name,
		&supplier.ContactName,
		&supplier.Email,
		&supplier.Phone,
		&supplier.AddressLine1,
		&supplier.AddressLine2,
		&supplier.City,
		&supplier.Region,
		&supplier.PostalCode,
		&supplier.Country,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Supplier{}, service.ErrNotFound
		}
		return Supplier{}, fmt.Errorf("getting supplier: %w", err)
	}

	return supplier, nil
}

func (c *Client) ListSuppliers(ctx context.Context) ([]Supplier, error) {
	rows, err := c.db.Query(ctx, `
		SELECT id, uuid, name, contact_name, email, phone, address_line1, address_line2, city, region, postal_code, country, created_at, updated_at, deleted_at
		FROM supplier
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("listing suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []Supplier
	for rows.Next() {
		var supplier Supplier
		if err := rows.Scan(
			&supplier.ID,
			&supplier.UUID,
			&supplier.Name,
			&supplier.ContactName,
			&supplier.Email,
			&supplier.Phone,
			&supplier.AddressLine1,
			&supplier.AddressLine2,
			&supplier.City,
			&supplier.Region,
			&supplier.PostalCode,
			&supplier.Country,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
			&supplier.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning supplier: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listing suppliers: %w", err)
	}

	return suppliers, nil
}
