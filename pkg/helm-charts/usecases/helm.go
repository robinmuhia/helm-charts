package usecases

import (
	"context"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/helpers"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

func (u *UsecaseHelmService) ProcessHelmChart(ctx context.Context, urlLink domain.HelmLinkInput) ([]*domain.ImageDetails, error) {
	validPath, err := helpers.ValidateURL(urlLink.Path)
	if err != nil {
		return nil, err
	}

	images, err := u.Infrastructure.Helm.ProcessHelmChart(ctx, validPath)
	if err != nil {
		return nil, err
	}

	return images, nil
}
