package helpers

import (
	"fmt"
	"os"
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
