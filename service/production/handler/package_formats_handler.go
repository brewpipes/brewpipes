package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type PackageFormatStore interface {
	CreatePackageFormat(context.Context, storage.PackageFormat) (storage.PackageFormat, error)
	GetPackageFormatByUUID(context.Context, string) (storage.PackageFormat, error)
	ListPackageFormats(context.Context) ([]storage.PackageFormat, error)
	UpdatePackageFormat(context.Context, int64, storage.PackageFormat) (storage.PackageFormat, error)
	DeletePackageFormat(context.Context, int64) error
}

// HandlePackageFormats handles [GET /package-formats] and [POST /package-formats].
func HandlePackageFormats(db PackageFormatStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			formats, err := db.ListPackageFormats(r.Context())
			if err != nil {
				service.InternalError(w, "error listing package formats", "error", err)
				return
			}

			service.JSON(w, dto.NewPackageFormatsResponse(formats))
		case http.MethodPost:
			var req dto.CreatePackageFormatRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			format := storage.PackageFormat{
				Name:              req.Name,
				Container:         req.Container,
				VolumePerUnit:     req.VolumePerUnit,
				VolumePerUnitUnit: req.VolumePerUnitUnit,
				IsActive:          true,
			}

			created, err := db.CreatePackageFormat(r.Context(), format)
			if err != nil {
				service.InternalError(w, "error creating package format", "error", err)
				return
			}

			slog.Info("package format created", "package_format_uuid", created.UUID, "name", created.Name)

			service.JSONCreated(w, dto.NewPackageFormatResponse(created))
		default:
			service.MethodNotAllowed(w)
		}
	}
}

// HandlePackageFormatByUUID handles [GET /package-formats/{uuid}], [PATCH /package-formats/{uuid}], and [DELETE /package-formats/{uuid}].
func HandlePackageFormatByUUID(db PackageFormatStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formatUUID := r.PathValue("uuid")
		if formatUUID == "" {
			http.Error(w, "invalid uuid", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			format, err := db.GetPackageFormatByUUID(r.Context(), formatUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "package format not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting package format", "error", err, "package_format_uuid", formatUUID)
				return
			}

			service.JSON(w, dto.NewPackageFormatResponse(format))
		case http.MethodPatch:
			var req dto.UpdatePackageFormatRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid request", http.StatusBadRequest)
				return
			}
			if err := req.Validate(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Resolve UUID to get current record
			current, err := db.GetPackageFormatByUUID(r.Context(), formatUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "package format not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting package format for update", "error", err, "package_format_uuid", formatUUID)
				return
			}

			// Apply partial update
			if req.Name != nil {
				current.Name = *req.Name
			}
			if req.Container != nil {
				current.Container = *req.Container
			}
			if req.VolumePerUnit != nil {
				current.VolumePerUnit = *req.VolumePerUnit
			}
			if req.VolumePerUnitUnit != nil {
				current.VolumePerUnitUnit = *req.VolumePerUnitUnit
			}
			if req.IsActive != nil {
				current.IsActive = *req.IsActive
			}

			updated, err := db.UpdatePackageFormat(r.Context(), current.ID, current)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "package format not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error updating package format", "error", err, "package_format_uuid", formatUUID)
				return
			}

			slog.Info("package format updated", "package_format_uuid", formatUUID, "name", updated.Name)

			service.JSON(w, dto.NewPackageFormatResponse(updated))
		case http.MethodDelete:
			// Resolve UUID to get internal ID
			format, err := db.GetPackageFormatByUUID(r.Context(), formatUUID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "package format not found", http.StatusNotFound)
				return
			} else if err != nil {
				service.InternalError(w, "error getting package format for delete", "error", err, "package_format_uuid", formatUUID)
				return
			}

			err = db.DeletePackageFormat(r.Context(), format.ID)
			if errors.Is(err, service.ErrNotFound) {
				http.Error(w, "package format not found", http.StatusNotFound)
				return
			} else if errors.Is(err, storage.ErrPackageFormatInUse) {
				http.Error(w, "package format is in use and cannot be deleted", http.StatusConflict)
				return
			} else if err != nil {
				service.InternalError(w, "error deleting package format", "error", err, "package_format_uuid", formatUUID)
				return
			}

			slog.Info("package format deleted", "package_format_uuid", formatUUID)

			w.WriteHeader(http.StatusNoContent)
		default:
			service.MethodNotAllowed(w)
		}
	}
}
