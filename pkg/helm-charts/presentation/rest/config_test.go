package rest_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"log"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/common"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/presentation"
	"github.com/stretchr/testify/require"
)

var testServer *http.Server
var baseURL string

func startTestServer(ctx context.Context, _ *testing.T) {
	port := "8081"
	os.Setenv(common.Port.String(), port)

	go func() {
		err := presentation.StartServer(ctx, 8081)
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Panicf("failed to start test server: %v", err)
		}
	}()

	baseURL = fmt.Sprintf("http://localhost:%s/api/v1", port)
}

func stopTestServer(ctx context.Context, t *testing.T) {
	if testServer != nil {
		err := testServer.Shutdown(ctx)
		require.NoError(t, err)
	}
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Start the test server
	startTestServer(ctx, &testing.T{})

	// Run the tests
	exitCode := m.Run()

	// Stop the server
	stopTestServer(ctx, &testing.T{})

	defer os.Exit(exitCode)
}
