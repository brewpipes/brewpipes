package handler_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/brewpipes/brewpipes/service/production/handler"
	"github.com/brewpipes/brewpipes/service/production/storage"
)

type VolumeStore struct {
	ListVolumesFunc     func(context.Context) ([]storage.Volume, error)
	GetVolumeByUUIDFunc func(context.Context, string) (storage.Volume, error)
	CreateVolumeFunc    func(context.Context, storage.Volume) (storage.Volume, error)
}

func (v VolumeStore) ListVolumes(ctx context.Context) ([]storage.Volume, error) {
	if v.ListVolumesFunc == nil {
		return nil, nil
	}
	return v.ListVolumesFunc(ctx)
}

func (v VolumeStore) GetVolumeByUUID(ctx context.Context, volumeUUID string) (storage.Volume, error) {
	if v.GetVolumeByUUIDFunc == nil {
		return storage.Volume{}, nil
	}
	return v.GetVolumeByUUIDFunc(ctx, volumeUUID)
}

func (v VolumeStore) CreateVolume(ctx context.Context, volume storage.Volume) (storage.Volume, error) {
	if v.CreateVolumeFunc == nil {
		return volume, nil
	}
	return v.CreateVolumeFunc(ctx, volume)
}

func TestVolumesHandler(t *testing.T) {
	handler := handler.HandleVolumes(VolumeStore{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/volumes", nil)
	handler.ServeHTTP(recorder, request)
}
