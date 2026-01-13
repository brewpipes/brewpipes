package handler_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type VolumeGetter struct {
	GetVolumesFunc func(context.Context) ([]storage.Volume, error)
}

func (v VolumeGetter) GetVolumes(ctx context.Context) ([]storage.Volume, error) {
	return v.GetVolumesFunc(ctx)
}

func TestVolumesHandler(t *testing.T) {
	handler := handler.HandleGetVolumes(VolumeGetter{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/volumes", nil)
	handler.ServeHTTP(recorder, request)
}
