package usecases

import (
	"context"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/helpers"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("github.com/robinmuhia/helm-charts/pkg/helm-charts/usecases")

func (u *UsecaseHelmService) ProcessHelmChart(ctx context.Context, urlLink *domain.HelmLinkInput) ([]*domain.ImageDetails, error) {
	ctx, span := tracer.Start(ctx, "ProcessHelmChart")
	defer span.End()

	validPath, err := helpers.ValidateURL(urlLink.Path)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, err
	}

	images, err := u.Infrastructure.Helm.ProcessHelmChart(ctx, validPath)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		return nil, err
	}

	return images, nil
}
