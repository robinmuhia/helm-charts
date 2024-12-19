package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

func TestUsecaseHelmService_ProcessHelmChart(t *testing.T) {
	type args struct {
		ctx     context.Context
		urlLink *domain.HelmLinkInput
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: process helm chart",
			args: args{
				ctx: context.Background(),
				urlLink: &domain.HelmLinkInput{
					Path: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
				},
			},
			wantErr: false,
		},
		{
			name: "fail: invalid url",
			args: args{
				ctx: context.Background(),
				urlLink: &domain.HelmLinkInput{
					Path: "foo",
				},
			},
			wantErr: true,
		},
		{
			name: "fail: fail to process chart",
			args: args{
				ctx: context.Background(),
				urlLink: &domain.HelmLinkInput{
					Path: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, mock := initializeMocks()

			if tt.name == "fail: fail to process chart" {
				mock.Helm.MockProcessHelmChartFn = func(_ context.Context, _ string) ([]*domain.ImageDetails, error) {
					return nil, fmt.Errorf("error")
				}
			}

			_, err := u.ProcessHelmChart(tt.args.ctx, tt.args.urlLink)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseHelmService.ProcessHelmChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
