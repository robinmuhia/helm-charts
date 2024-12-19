package helm

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestService_fetchImageDetails(t *testing.T) {
	type args struct {
		image string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: get image details",
			args: args{
				image: "docker.io/bitnami/redis",
			},
			wantErr: false,
		},
		{
			name: "fail: invalid image path",
			args: args{
				image: "docker.io/bitnami/redis@latest@v1",
			},
			wantErr: true,
		},
		{
			name: "fail: fail to pull from remote",
			args: args{
				image: "docker.io/bitnami/rendisy",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.Writer(), "HelmService: ", log.LstdFlags)

			s := NewHelmService(logger)

			_, err := s.fetchImageDetails(tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.fetchImageDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_downloadHelmChart(t *testing.T) {
	type args struct {
		ctx context.Context
		url string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: download chart",
			args: args{
				ctx: context.Background(),
				url: "http://example.com/helm-chart.tgz",
			},
			wantErr: false,
		},
		{
			name: "fail: nil context",
			args: args{
				ctx: nil,
				url: "http://example.com/helm-chart.tgz",
			},
			wantErr: true,
		},
		{
			name: "fail: invalid status code",
			args: args{
				ctx: context.Background(),
				url: "http://example.com/helm-chart.tgz",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.Writer(), "HelmService: ", log.LstdFlags)

			s := NewHelmService(logger)

			if tt.name == "success: download chart" {
				httpmock.RegisterResponder(http.MethodGet, tt.args.url, func(_ *http.Request) (*http.Response, error) {
					resp := []byte("fake tarball content")

					return httpmock.NewJsonResponse(http.StatusOK, resp)
				})
			}

			if tt.name == "fail: invalid status code" {
				httpmock.RegisterResponder(http.MethodGet, tt.args.url, func(_ *http.Request) (*http.Response, error) {
					return httpmock.NewJsonResponse(http.StatusBadRequest, nil)
				})
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			_, err := s.downloadHelmChart(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.downloadHelmChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_parseHelmChart(t *testing.T) {
	type args struct {
		chartPath string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: parse helm chart and extract images",
			args: args{
				chartPath: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
			},
			wantErr: false,
		},
		{
			name: "error: helm command fails",
			args: args{
				chartPath: "https://github.com/",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.Writer(), "HelmService: ", log.LstdFlags)

			s := NewHelmService(logger)

			chartPath, err := s.downloadHelmChart(context.Background(), tt.args.chartPath)
			if err != nil {
				t.Errorf("failed to download chart")
			}

			_, err = s.parseHelmChart(chartPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.parseHelmChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_ProcessHelmChart(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: parse helm chart and extract images",
			args: args{
				ctx:  context.Background(),
				path: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
			},
			wantErr: false,
		},
		{
			name: "error: nil context",
			args: args{
				ctx:  nil,
				path: "https://github.com/",
			},
			wantErr: true,
		},
		{
			name: "error: bad url (no helm chart to download)",
			args: args{
				ctx:  context.Background(),
				path: "https://github.com/",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.New(log.Writer(), "HelmService: ", log.LstdFlags)

			s := NewHelmService(logger)

			_, err := s.ProcessHelmChart(tt.args.ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ProcessHelmChart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
