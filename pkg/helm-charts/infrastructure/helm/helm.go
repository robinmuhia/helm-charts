package helm

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/robinmuhia/helm-charts/pkg/helm-charts/domain"
)

// Service encapsulates the logic for processing Helm charts and fetching image details.
type Service struct {
	logger *log.Logger
}

// NewHelmService initializes and returns a new Service instance.
func NewHelmService(logger *log.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

// fetchImageDetails retrieves image metadata using the container registry API.
func (s *Service) fetchImageDetails(image string) (*domain.ImageDetails, error) {
	ref, err := name.ParseReference(image)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(ref)
	if err != nil {
		return nil, err
	}

	manifest, err := img.Manifest()
	if err != nil {
		return nil, err
	}

	size := int64(0)
	for _, layer := range manifest.Layers {
		size += layer.Size
	}

	return &domain.ImageDetails{
		Image:  image,
		Size:   size,
		Layers: len(manifest.Layers),
	}, nil
}

// downloadHelmChart downloads a Helm chart from a URL and saves it locally.
func (s *Service) downloadHelmChart(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req) // codeql:ignore
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download Helm chart: received status code %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp("", "helm-chart-*.tgz")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write Helm chart to file: %w", err)
	}

	return tmpFile.Name(), nil
}

// parseHelmChart extracts image references from a Helm chart.
func (s *Service) parseHelmChart(chartPath string) ([]string, error) {
	cmd := exec.Command("helm", "template", chartPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to render helm chart: %w", err)
	}

	var images []string

	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.Contains(line, "image:") {
			imageParts := strings.Split(line, ": ")
			if len(imageParts) > 1 {
				image := strings.TrimSpace(imageParts[1])
				image = strings.Trim(image, "\"")
				images = append(images, image)
			}
		}
	}

	return images, nil
}

// ProcessChartHandler handles HTTP requests to process Helm charts.
func (s *Service) ProcessHelmChart(ctx context.Context, path string) ([]*domain.ImageDetails, error) {
	chartPath, err := s.downloadHelmChart(ctx, path)
	if err != nil {
		return nil, err
	}

	images, err := s.parseHelmChart(chartPath)
	if err != nil {
		return nil, err
	}

	var results []*domain.ImageDetails

	var wg sync.WaitGroup

	var mu sync.Mutex

	wg.Add(len(images))

	for _, image := range images {
		go func(image string) {
			defer wg.Done()

			details, err := s.fetchImageDetails(image)
			if err != nil {
				s.logger.Printf("Failed to fetch details for image %s: %v", image, err)
			}

			mu.Lock()
			results = append(results, details)
			mu.Unlock()
		}(image)
	}

	wg.Wait()

	return results, nil
}
