package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

const (
	batchImportMaxUploadSize = 5 << 20
	batchImportMaxRows       = 1000
)

var batchImportAllowedHeaders = map[string]struct{}{
	"short_name": {},
	"brew_date":  {},
	"notes":      {},
}

// HandleBatchImport handles [POST /batches/import].
func HandleBatchImport(db BatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, batchImportMaxUploadSize)
		if err := r.ParseMultipartForm(batchImportMaxUploadSize); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				http.Error(w, "file too large", http.StatusRequestEntityTooLarge)
				return
			}
			slog.Warn("invalid batch import form", "error", err)
			http.Error(w, "invalid multipart form", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1

		headers, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				http.Error(w, "missing header row", http.StatusBadRequest)
				return
			}
			slog.Warn("unable to read batch import header", "error", err)
			http.Error(w, "invalid csv", http.StatusBadRequest)
			return
		}

		headerIndex := make(map[string]int, len(headers))
		for idx, rawHeader := range headers {
			header := strings.TrimSpace(rawHeader)
			if header == "" {
				http.Error(w, "header value is required", http.StatusBadRequest)
				return
			}
			if _, ok := batchImportAllowedHeaders[header]; !ok {
				http.Error(w, fmt.Sprintf("unknown header: %s", header), http.StatusBadRequest)
				return
			}
			if _, exists := headerIndex[header]; exists {
				http.Error(w, fmt.Sprintf("duplicate header: %s", header), http.StatusBadRequest)
				return
			}
			headerIndex[header] = idx
		}

		if _, ok := headerIndex["short_name"]; !ok {
			http.Error(w, "short_name header is required", http.StatusBadRequest)
			return
		}

		var records [][]string
		for {
			record, err := reader.Read()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				slog.Warn("unable to read batch import row", "error", err)
				http.Error(w, "invalid csv", http.StatusBadRequest)
				return
			}
			records = append(records, record)
			if len(records) > batchImportMaxRows {
				http.Error(w, "row limit exceeded", http.StatusBadRequest)
				return
			}
		}

		results := make([]dto.BatchImportRowResult, 0, len(records))
		createdCount := 0
		failedCount := 0

		for idx, record := range records {
			rowNumber := idx + 2
			if len(record) != len(headers) {
				msg := "invalid column count"
				results = append(results, dto.BatchImportRowResult{
					Row:    rowNumber,
					Status: "error",
					Error:  &msg,
				})
				failedCount++
				continue
			}

			shortName := strings.TrimSpace(record[headerIndex["short_name"]])
			if shortName == "" {
				msg := "short_name is required"
				results = append(results, dto.BatchImportRowResult{
					Row:    rowNumber,
					Status: "error",
					Error:  &msg,
				})
				failedCount++
				continue
			}

			var brewDate *time.Time
			if idx, ok := headerIndex["brew_date"]; ok {
				value := strings.TrimSpace(record[idx])
				if value != "" {
					parsed, err := time.Parse("2006-01-02", value)
					if err != nil {
						msg := "invalid brew_date"
						results = append(results, dto.BatchImportRowResult{
							Row:    rowNumber,
							Status: "error",
							Error:  &msg,
						})
						failedCount++
						continue
					}
					brewDate = &parsed
				}
			}

			var notes *string
			if idx, ok := headerIndex["notes"]; ok {
				value := strings.TrimSpace(record[idx])
				if value != "" {
					notes = &value
				}
			}

			batch := storage.Batch{
				ShortName: shortName,
				BrewDate:  brewDate,
				Notes:     notes,
			}

			created, err := db.CreateBatch(r.Context(), batch)
			if err != nil {
				slog.Error("error creating batch from import", "error", err, "row", rowNumber)
				msg := "create batch failed"
				results = append(results, dto.BatchImportRowResult{
					Row:    rowNumber,
					Status: "error",
					Error:  &msg,
				})
				failedCount++
				continue
			}

			createdCount++
			batchResponse := dto.NewBatchResponse(created)
			results = append(results, dto.BatchImportRowResult{
				Row:    rowNumber,
				Status: "created",
				Batch:  &batchResponse,
			})
		}

		resp := dto.BatchImportResponse{
			Totals: dto.BatchImportTotals{
				TotalRows: len(records),
				Created:   createdCount,
				Failed:    failedCount,
			},
			Results: results,
		}

		service.JSON(w, resp)
	}
}
