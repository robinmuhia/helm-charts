package usecases_test

import (
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/infrastructure"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/infrastructure/helm/mock"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/usecases"
)

// Set up mocks
type Mock struct {
	Helm *mock.HelmMock
}

func initializeMocks() (*usecases.UsecaseHelmService, *Mock) {
	fakeHelm := mock.NewHelmServiceMock()

	infrastructure := infrastructure.NewInfrastructureInteractor(fakeHelm)

	usecases := usecases.NewUsecaseHelmImpl(*infrastructure)

	return usecases, &Mock{
		Helm: fakeHelm,
	}
}
