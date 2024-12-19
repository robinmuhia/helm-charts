package helpers

import (
	"testing"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/common"
)

func TestGetEnvVar(t *testing.T) {
	type args struct {
		envVarName string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success: get env",
			args: args{
				envVarName: common.Environment.String(),
			},
			want:    "test",
			wantErr: false,
		},
		{
			name: "fail: fail to get env",
			args: args{
				envVarName: "FOO",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEnvVar(tt.args.envVarName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnvVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("GetEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	type args struct {
		userInputURL string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success: parse url",
			args: args{
				userInputURL: "https://github.com",
			},
			want:    "https://github.com",
			wantErr: false,
		},
		{
			name: "fail: invalid url",
			args: args{
				userInputURL: "https://  htpps://  github.com.com",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "fail: url not in allowed hosts",
			args: args{
				userInputURL: "https://bar.com",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "fail: invalid scheme",
			args: args{
				userInputURL: "ssh://github.com",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateURL(tt.args.userInputURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ValidateURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
