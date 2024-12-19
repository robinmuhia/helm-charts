package usecases

import "github.com/robinmuhia/helm-charts/pkg/helm-charts/infrastructure"

type UsecaseHelmService struct {
	Infrastructure infrastructure.Infrastructure
}

func NewUsecaseHelmImpl(infra infrastructure.Infrastructure) *UsecaseHelmService {
	return &UsecaseHelmService{
		Infrastructure: infra,
	}
}
