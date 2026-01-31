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
	GetSupplier(context.Context, int64) (storage.Supplier, error)
	CreateSupplier(context.Context, storage.Supplier) (storage.Supplier, error)
}

// HandleSuppliers handles [GET /suppliers] and [POST /suppliers].
func HandleSuppliers(db SupplierStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			suppliers, err := db.ListSuppliers(r.Context())
			if err != nil {
				slog.Error("error listing suppliers", "error", err)
				service.InternalError(w, err.Error())
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
				slog.Error("error creating supplier", "error", err)
				service.InternalError(w, err.Error())
				return
			}

			service.JSON(w, dto.NewSupplierResponse(created))
		default:
			methodNotAllowed(w)
		}
	}
}

// HandleSupplierByID handles [GET /suppliers/{id}].
func HandleSupplierByID(db SupplierStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		idValue := r.PathValue("id")
		if idValue == "" {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		supplierID, err := parseInt64Param(idValue)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		supplier, err := db.GetSupplier(r.Context(), supplierID)
		if errors.Is(err, service.ErrNotFound) {
			http.Error(w, "supplier not found", http.StatusNotFound)
			return
		} else if err != nil {
			slog.Error("error getting supplier", "error", err)
			service.InternalError(w, err.Error())
			return
		}

		service.JSON(w, dto.NewSupplierResponse(supplier))
	}
}
