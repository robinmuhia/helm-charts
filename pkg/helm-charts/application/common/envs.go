package common

type EnvironmentVariable string

const (
	Environment             EnvironmentVariable = "ENVIRONMENT"
	Port                    EnvironmentVariable = "PORT"
	JaegerCollectorEndpoint EnvironmentVariable = "JAEGER_ENDPOINT"
)

// String converts environment variable to its string type
func (e EnvironmentVariable) String() string {
	return string(e)
}
