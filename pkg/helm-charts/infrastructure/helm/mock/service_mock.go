package mock

import (
	"context"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

// HelmMock mocks the interface for methods exposed our helm infrastructure
type HelmMock struct {
	MockProcessHelmChartFn func(ctx context.Context, path string) ([]*domain.ImageDetails, error)
}

// NewHelmServiceMock ...
func NewHelmServiceMock() *HelmMock {
	return &HelmMock{
		MockProcessHelmChartFn: func(_ context.Context, _ string) ([]*domain.ImageDetails, error) {
			return []*domain.ImageDetails{
				{
					Image:  "nginx:1.16.0",
					Size:   123456,
					Layers: 2,
				},
			}, nil
		},
	}
}

// ProcessHelmChart mocks the implementation of processing a helm chart
func (h HelmMock) ProcessHelmChart(ctx context.Context, path string) ([]*domain.ImageDetails, error) {
	return h.MockProcessHelmChartFn(ctx, path)
}
