package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

func TestHandlersInterfacesImpl_ParseHelmLink(t *testing.T) {
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}

	validInput := domain.HelmLinkInput{
		Path: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
	}

	validPayload, err := json.Marshal(validInput)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	invalidDomain := domain.HelmLinkInput{
		Path: "https://github.com/helm/examples/",
	}

	invalidDomainPayload, err := json.Marshal(invalidDomain)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	type helmLinkInput struct {
		Path domain.HelmLinkInput `json:"url_link"`
	}

	invalidInput := helmLinkInput{
		Path: domain.HelmLinkInput{
			Path: "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz",
		},
	}

	invalidInputPayload, err := json.Marshal(invalidInput)
	if err != nil {
		t.Errorf("failed to marshal payload")
		return
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success: get image data",
			args: args{
				url:        fmt.Sprintf("%s/helm-link", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(validPayload),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: fail to bind json",
			args: args{
				url:        fmt.Sprintf("%s/helm-link", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidInputPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "fail: invalid download",
			args: args{
				url:        fmt.Sprintf("%s/helm-link", baseURL),
				httpMethod: http.MethodPost,
				body:       bytes.NewBuffer(invalidDomainPayload),
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Close = true

			client := http.DefaultClient

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("request error: %s", err)
				return
			}

			if resp == nil && !tt.wantErr {
				t.Errorf("nil response")
				return
			}

			defer resp.Body.Close()

			dataResponse, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read request body: %s", err)
				return
			}

			if dataResponse == nil {
				t.Errorf("nil response data")
				return
			}

			data := map[string]interface{}{}

			err = json.Unmarshal(dataResponse, &data)
			if tt.wantErr && err != nil {
				t.Errorf("bad data returned: %v", err)
				return
			}

			if tt.wantErr {
				errMsg, ok := data["error"]
				if !ok {
					t.Errorf("expected error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected")
					return
				}
			}

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %s", tt.wantStatus, resp.Status)
				return
			}
		})
	}
}
