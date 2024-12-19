package presentation

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/common"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/application/helpers"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/infrastructure"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/infrastructure/helm"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/presentation/rest"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var allowedOriginPatterns = []string{
	`^https://.+\.web\.app$`,
}

// Compile the regex patterns into a slice of *regexp.Regexp
func compilePatterns(patterns []string) []*regexp.Regexp {
	var compiledPatterns []*regexp.Regexp

	for _, pattern := range patterns {
		compiledPattern := regexp.MustCompile(pattern)
		compiledPatterns = append(compiledPatterns, compiledPattern)
	}

	return compiledPatterns
}

// Check if the origin is allowed by matching it against the compiled regex patterns
func isAllowedOrigin(origin string, compiledPatterns []*regexp.Regexp) bool {
	for _, pattern := range compiledPatterns {
		if pattern.MatchString(origin) {
			return true
		}
	}

	return false
}

// StartServer sets up gin
func StartServer(_ context.Context, port int) error {
	logger := log.New(log.Writer(), "HelmService: ", log.LstdFlags)

	helm := helm.NewHelmService(logger)

	infra := infrastructure.NewInfrastructureInteractor(helm)

	usecases := usecases.NewUsecaseHelmImpl(*infra)

	r := gin.Default()

	SetupRoutes(r, usecases)

	addr := fmt.Sprintf(":%d", port)

	if err := r.Run(addr); err != nil {
		return err
	}

	return nil
}

func SetupRoutes(r *gin.Engine, usecases *usecases.UsecaseHelmService) {
	compiledPatterns := compilePatterns(allowedOriginPatterns)

	r.Use(cors.New(cors.Config{
		AllowWildcard: true,
		AllowMethods:  []string{http.MethodPut, http.MethodPatch, http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{
			"Accept",
			"Accept-Charset",
			"Accept-Language",
			"Accept-Encoding",
			"Origin",
			"Host",
			"User-Agent",
			"Content-Length",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if origin == "http://localhost:8080" {
				return true
			}

			allowed := isAllowedOrigin(origin, compiledPatterns)
			return allowed
		},
		MaxAge:          12 * time.Hour,
		AllowWebSockets: true,
	}))

	environment, err := helpers.GetEnvVar(common.Environment.String())
	if err != nil {
		log.Panic(err)
	}

	handlers := rest.NewHandlersInterfaces(usecases)

	r.Use(otelgin.Middleware(fmt.Sprintf("helm-chart-%v", environment)))

	r.Use(gin.Recovery())

	apiV1routes := r.Group("api/v1")

	// endpoints
	apiV1routes.POST("/helm-link", handlers.ParseHelmLink)
}
