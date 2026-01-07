package auth

import "log/slog"

type Service struct {
}

func NewService() (*Service, error) {
	slog.Info("auth service initialized")
	return &Service{}, nil
}

func (s *Service) Start() error {
	slog.Info("auth service starting")
	slog.Info("auth service running")
	return nil
}

func (s *Service) Stop() error {
	slog.Info("auth service stopping")
	slog.Info("auth service stopped")
	return nil
}

func (s *Service) HealthCheck() bool {
	slog.Info("auth service health check")
	return true
}
