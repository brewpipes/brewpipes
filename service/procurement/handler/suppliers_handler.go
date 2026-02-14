package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/procurement/handler/dto"
	"github.com/brewpipes/brewpipes/service/procurement/storage"
)

type SupplierStore interface {
	ListSuppliers(context.Context) ([]storage.Supplier, error)
	GetSupplierByUUID(context.Context, string) (storage.Supplier, error)
	CreateSupplier(context.Context, storage.Supplier) (storage.Supplier, error)
	UpdateSupplierByUUID(context.Context, string, storage.SupplierUpdate) (storage.Supplier, error)
}

// HandleSuppliers handles [GET /suppliers] and [POST /suppliers].
func HandleSuppliers(db SupplierStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			suppliers, err := db.ListSuppliers(r.Context())
			if err != nil {
				service.InternalError(w, "error listing suppliers", "error", err)
				return
			}

			service.JSON(w, dto.NewSuppliersResponse(suppliers))
		case http.MethodPost:
			var req dto.CreateSupplierRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			supplier := storage.Supplier{
				Name:         req.Name,
				ContactName:  req.ContactName,
				Email:        req.Email,
				Phone:        req.Phone,
				AddressLine1: req.AddressLine1,
				AddressLine2: req.AddressLine2,
				City:         req.City,
				Region:       req.Region,
				PostalCode:   req.PostalCode,
				Country:      req.Country,
			}

			created, err := db.CreateSupplier(r.Context(), supplier)
			if err != nil {
				service.InternalError(w, "error creating supplier", "error", err)
				return
			}

			service.JSONCreated(w, dto.NewSupplierResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandleSupplierByUUID handles [GET /suppliers/{uuid}] and [PATCH /suppliers/{uuid}].
func HandleSupplierByUUID(db SupplierStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		supplierUUID := r.PathValue("uuid")
		if supplierUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			supplier, err := db.GetSupplierByUUID(r.Context(), supplierUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "supplier not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting supplier", "error", err)
				return
			}

			service.JSON(w, dto.NewSupplierResponse(supplier))
		case http.MethodPatch:
			var req dto.UpdateSupplierRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			update := storage.SupplierUpdate{
				Name:         req.Name,
				ContactName:  req.ContactName,
				Email:        req.Email,
				Phone:        req.Phone,
				AddressLine1: req.AddressLine1,
				AddressLine2: req.AddressLine2,
				City:         req.City,
				Region:       req.Region,
				PostalCode:   req.PostalCode,
				Country:      req.Country,
			}

			supplier, err := db.UpdateSupplierByUUID(r.Context(), supplierUUID, update)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "supplier not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error updating supplier", "error", err)
				return
			}

			slog.Info("supplier updated", "supplier_uuid", supplierUUID, "name", supplier.Name)
			service.JSON(w, dto.NewSupplierResponse(supplier))
		default:
			service.MethodNotAllowed(w)
		}
	}
}
