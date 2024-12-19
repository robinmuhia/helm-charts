package helpers

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// GetEnvVar retrieves the environment variable with the supplied name and fails
// if it is not able to do so
func GetEnvVar(envVarName string) (string, error) {
	envVar := os.Getenv(envVarName)
	if envVar == "" {
		return "", fmt.Errorf("the environment variable '%s' is not set", envVarName)
	}

	return envVar, nil
}

// Validate and sanitize the URL before making a request.
func ValidateURL(userInputURL string) (string, error) {
	parsedURL, err := url.Parse(userInputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid URL scheme: only HTTP/HTTPS are allowed")
	}

	trustedDomains := []string{"bitnami.com", "helm.sh", "artifacthub.io", "hashicorp.com", "github.io", "jetstack.io", "github.com"}

	isValidHost := false

	for _, domain := range trustedDomains {
		if strings.Contains(parsedURL.Host, domain) {
			isValidHost = true
			break
		}
	}

	if !isValidHost {
		return "", fmt.Errorf("untrusted domain: %s", parsedURL.Host)
	}

	return parsedURL.String(), nil
}
