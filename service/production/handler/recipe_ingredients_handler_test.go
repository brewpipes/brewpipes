package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brewpipes/brewpipes/service"
	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/handler/dto"
	"github.com/brewpipes/brewpipes/service/production/storage"
	"github.com/gofrs/uuid/v5"
)

// mockRecipeIngredientStore implements RecipeIngredientStore for testing.
type mockRecipeIngredientStore struct {
	ListRecipeIngredientsFunc  func(ctx context.Context, recipeUUID string) ([]storage.RecipeIngredient, error)
	GetRecipeIngredientFunc    func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error)
	CreateRecipeIngredientFunc func(ctx context.Context, ri storage.RecipeIngredient) (storage.RecipeIngredient, error)
	UpdateRecipeIngredientFunc func(ctx context.Context, ingredientUUID string, ri storage.RecipeIngredient) (storage.RecipeIngredient, error)
	DeleteRecipeIngredientFunc func(ctx context.Context, ingredientUUID string) error
}

func (m *mockRecipeIngredientStore) ListRecipeIngredients(ctx context.Context, recipeUUID string) ([]storage.RecipeIngredient, error) {
	if m.ListRecipeIngredientsFunc != nil {
		return m.ListRecipeIngredientsFunc(ctx, recipeUUID)
	}
	return nil, nil
}

func (m *mockRecipeIngredientStore) GetRecipeIngredient(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
	if m.GetRecipeIngredientFunc != nil {
		return m.GetRecipeIngredientFunc(ctx, ingredientUUID)
	}
	return storage.RecipeIngredient{}, service.ErrNotFound
}

func (m *mockRecipeIngredientStore) CreateRecipeIngredient(ctx context.Context, ri storage.RecipeIngredient) (storage.RecipeIngredient, error) {
	if m.CreateRecipeIngredientFunc != nil {
		return m.CreateRecipeIngredientFunc(ctx, ri)
	}
	ri.ID = 1
	ri.UUID = uuid.Must(uuid.NewV4())
	return ri, nil
}

func (m *mockRecipeIngredientStore) UpdateRecipeIngredient(ctx context.Context, ingredientUUID string, ri storage.RecipeIngredient) (storage.RecipeIngredient, error) {
	if m.UpdateRecipeIngredientFunc != nil {
		return m.UpdateRecipeIngredientFunc(ctx, ingredientUUID, ri)
	}
	return ri, nil
}

func (m *mockRecipeIngredientStore) DeleteRecipeIngredient(ctx context.Context, ingredientUUID string) error {
	if m.DeleteRecipeIngredientFunc != nil {
		return m.DeleteRecipeIngredientFunc(ctx, ingredientUUID)
	}
	return nil
}

// mockRecipeChecker implements RecipeExistenceChecker for testing.
type mockRecipeChecker struct {
	GetRecipeFunc func(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error)
}

func (m *mockRecipeChecker) GetRecipe(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error) {
	if m.GetRecipeFunc != nil {
		return m.GetRecipeFunc(ctx, recipeUUID, opts)
	}
	// Return a default recipe with ID=1 for FK comparison
	recipe := storage.Recipe{}
	recipe.ID = 1
	recipe.UUID = uuid.Must(uuid.FromString(recipeUUID))
	return recipe, nil
}

func TestHandleRecipeIngredients_List(t *testing.T) {
	recipeUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		recipeUUID     string
		setupStore     func(*mockRecipeIngredientStore)
		setupRecipes   func(*mockRecipeChecker)
		expectedStatus int
	}{
		{
			name:       "success - empty list",
			recipeUUID: recipeUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.ListRecipeIngredientsFunc = func(ctx context.Context, recipeUUID string) ([]storage.RecipeIngredient, error) {
					return []storage.RecipeIngredient{}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "success - with ingredients",
			recipeUUID: recipeUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.ListRecipeIngredientsFunc = func(ctx context.Context, recipeUUID string) ([]storage.RecipeIngredient, error) {
					return []storage.RecipeIngredient{
						{
							RecipeID:       1,
							IngredientType: "fermentable",
							Amount:         10.0,
							AmountUnit:     "kg",
							UseStage:       "mash",
							ScalingFactor:  1.0,
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:       "recipe not found",
			recipeUUID: "550e8400-e29b-41d4-a716-446655440999",
			setupRecipes: func(m *mockRecipeChecker) {
				m.GetRecipeFunc = func(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error) {
					return storage.Recipe{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "empty recipe uuid",
			recipeUUID:     "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecipeIngredientStore{}
			recipes := &mockRecipeChecker{}

			if tt.setupStore != nil {
				tt.setupStore(store)
			}
			if tt.setupRecipes != nil {
				tt.setupRecipes(recipes)
			}

			h := handler.HandleRecipeIngredients(store, recipes)

			req := httptest.NewRequest(http.MethodGet, "/recipes/"+tt.recipeUUID+"/ingredients", nil)
			req.SetPathValue("uuid", tt.recipeUUID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandleRecipeIngredients_Create(t *testing.T) {
	recipeUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name           string
		recipeUUID     string
		body           dto.RecipeIngredientRequest
		setupStore     func(*mockRecipeIngredientStore)
		setupRecipes   func(*mockRecipeChecker)
		expectedStatus int
	}{
		{
			name:       "success - minimal fields",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType: "fermentable",
				Amount:         10.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:       "success - all fields",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType:        "hop",
				Amount:                2.5,
				AmountUnit:            "kg",
				UseStage:              "boil",
				UseType:               strPtr("bittering"),
				TimingDurationMinutes: intPtr(60),
				AlphaAcidAssumed:      float64Ptr(12.5),
				ScalingFactor:         float64Ptr(0.9),
				SortOrder:             intPtr(10),
				Notes:                 strPtr("Columbus for bittering"),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:       "validation error - invalid ingredient_type",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType: "invalid",
				Amount:         10.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "validation error - amount <= 0",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType: "fermentable",
				Amount:         0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "validation error - alpha_acid on non-hop",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType:   "fermentable",
				Amount:           10.0,
				AmountUnit:       "kg",
				UseStage:         "mash",
				AlphaAcidAssumed: float64Ptr(12.5),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:       "recipe not found",
			recipeUUID: "550e8400-e29b-41d4-a716-446655440999",
			body: dto.RecipeIngredientRequest{
				IngredientType: "fermentable",
				Amount:         10.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			setupRecipes: func(m *mockRecipeChecker) {
				m.GetRecipeFunc = func(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error) {
					return storage.Recipe{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:       "validation error - invalid ingredient_uuid format",
			recipeUUID: recipeUUID,
			body: dto.RecipeIngredientRequest{
				IngredientUUID: strPtr("not-a-valid-uuid"),
				IngredientType: "fermentable",
				Amount:         10.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecipeIngredientStore{}
			recipes := &mockRecipeChecker{}

			if tt.setupStore != nil {
				tt.setupStore(store)
			}
			if tt.setupRecipes != nil {
				tt.setupRecipes(recipes)
			}

			h := handler.HandleRecipeIngredients(store, recipes)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/recipes/"+tt.recipeUUID+"/ingredients", bytes.NewReader(body))
			req.SetPathValue("uuid", tt.recipeUUID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleRecipeIngredient_Get(t *testing.T) {
	recipeUUID := "550e8400-e29b-41d4-a716-446655440000"
	ingredientUUID := "660e8400-e29b-41d4-a716-446655440010"

	tests := []struct {
		name           string
		recipeUUID     string
		ingredientUUID string
		setupStore     func(*mockRecipeIngredientStore)
		setupRecipes   func(*mockRecipeChecker)
		expectedStatus int
	}{
		{
			name:           "success",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       1,
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ingredient not found",
			recipeUUID:     recipeUUID,
			ingredientUUID: "660e8400-e29b-41d4-a716-446655440999",
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ingredient belongs to different recipe",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       2, // Different recipe (mock returns ID=1)
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "recipe not found",
			recipeUUID:     "550e8400-e29b-41d4-a716-446655440999",
			ingredientUUID: ingredientUUID,
			setupRecipes: func(m *mockRecipeChecker) {
				m.GetRecipeFunc = func(ctx context.Context, recipeUUID string, opts *storage.RecipeQueryOpts) (storage.Recipe, error) {
					return storage.Recipe{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecipeIngredientStore{}
			recipes := &mockRecipeChecker{}

			if tt.setupStore != nil {
				tt.setupStore(store)
			}
			if tt.setupRecipes != nil {
				tt.setupRecipes(recipes)
			}

			h := handler.HandleRecipeIngredient(store, recipes)

			req := httptest.NewRequest(http.MethodGet, "/recipes/"+tt.recipeUUID+"/ingredients/"+tt.ingredientUUID, nil)
			req.SetPathValue("uuid", tt.recipeUUID)
			req.SetPathValue("ingredient_uuid", tt.ingredientUUID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleRecipeIngredient_Update(t *testing.T) {
	recipeUUID := "550e8400-e29b-41d4-a716-446655440000"
	ingredientUUID := "660e8400-e29b-41d4-a716-446655440010"

	tests := []struct {
		name           string
		recipeUUID     string
		ingredientUUID string
		body           dto.RecipeIngredientRequest
		setupStore     func(*mockRecipeIngredientStore)
		setupRecipes   func(*mockRecipeChecker)
		expectedStatus int
	}{
		{
			name:           "success",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType: "fermentable",
				Amount:         15.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       1,
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ingredient not found",
			recipeUUID:     recipeUUID,
			ingredientUUID: "660e8400-e29b-41d4-a716-446655440999",
			body: dto.RecipeIngredientRequest{
				IngredientType: "fermentable",
				Amount:         15.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "validation error",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			body: dto.RecipeIngredientRequest{
				IngredientType: "invalid",
				Amount:         15.0,
				AmountUnit:     "kg",
				UseStage:       "mash",
			},
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       1,
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecipeIngredientStore{}
			recipes := &mockRecipeChecker{}

			if tt.setupStore != nil {
				tt.setupStore(store)
			}
			if tt.setupRecipes != nil {
				tt.setupRecipes(recipes)
			}

			h := handler.HandleRecipeIngredient(store, recipes)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPatch, "/recipes/"+tt.recipeUUID+"/ingredients/"+tt.ingredientUUID, bytes.NewReader(body))
			req.SetPathValue("uuid", tt.recipeUUID)
			req.SetPathValue("ingredient_uuid", tt.ingredientUUID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestHandleRecipeIngredient_Delete(t *testing.T) {
	recipeUUID := "550e8400-e29b-41d4-a716-446655440000"
	ingredientUUID := "660e8400-e29b-41d4-a716-446655440010"

	tests := []struct {
		name           string
		recipeUUID     string
		ingredientUUID string
		setupStore     func(*mockRecipeIngredientStore)
		setupRecipes   func(*mockRecipeChecker)
		expectedStatus int
	}{
		{
			name:           "success",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       1,
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "ingredient not found",
			recipeUUID:     recipeUUID,
			ingredientUUID: "660e8400-e29b-41d4-a716-446655440999",
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{}, service.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "ingredient belongs to different recipe",
			recipeUUID:     recipeUUID,
			ingredientUUID: ingredientUUID,
			setupStore: func(m *mockRecipeIngredientStore) {
				m.GetRecipeIngredientFunc = func(ctx context.Context, ingredientUUID string) (storage.RecipeIngredient, error) {
					return storage.RecipeIngredient{
						RecipeID:       2, // Different recipe (mock returns ID=1)
						IngredientType: "fermentable",
						Amount:         10.0,
						AmountUnit:     "kg",
						UseStage:       "mash",
						ScalingFactor:  1.0,
					}, nil
				}
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &mockRecipeIngredientStore{}
			recipes := &mockRecipeChecker{}

			if tt.setupStore != nil {
				tt.setupStore(store)
			}
			if tt.setupRecipes != nil {
				tt.setupRecipes(recipes)
			}

			h := handler.HandleRecipeIngredient(store, recipes)

			req := httptest.NewRequest(http.MethodDelete, "/recipes/"+tt.recipeUUID+"/ingredients/"+tt.ingredientUUID, nil)
			req.SetPathValue("uuid", tt.recipeUUID)
			req.SetPathValue("ingredient_uuid", tt.ingredientUUID)
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
			}
		})
	}
}

// Helper functions for creating pointers
func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}
