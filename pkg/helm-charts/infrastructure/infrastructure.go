package infrastructure

import (
	"context"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

// Helm is the interface for methods exposed our helm infrastructure
type Helm interface {
	ProcessHelmChart(ctx context.Context, path string) ([]*domain.ImageDetails, error)
}

// Infrastructure implements the infrastructure interface(s)
type Infrastructure struct {
	Helm Helm
}

// NewInfrastructureInteractor initializes a new Infrastructure
func NewInfrastructureInteractor(helm Helm) *Infrastructure {
	return &Infrastructure{
		Helm: helm,
	}
}
