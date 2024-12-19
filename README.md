# Helm Chart Image Processor

![Lint and Tests](https://github.com/robinmuhia/helm-charts/actions/workflows/ci.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/robinmuhia/helm-charts/badge.svg?branch=main)](https://coveralls.io/github/robinmuhia/helm-charts?branch=main)

This project provides an API that processes a Helm chart to extract container image details. The API accepts a Helm chart path, identifies all container images, downloads metadata from their respective registries, and returns the image size and number of layers.

---

## Features

- **Helm Chart Parsing:** Processes a Helm chart using `helm template`.
- **Image Extraction:** Extracts container image references from Kubernetes manifests.
- **Image Metadata Retrieval:** Fetches size and layer details for each image using Docker registries.
- **REST API:** Exposes functionality through a simple HTTP POST API.

---

## Prerequisites

- **Helm CLI**: Install the Helm CLI from [Helm Installation Guide](https://helm.sh/docs/intro/install/).
- **Go**: Install Go (1.20 or later).
- **Docker Registry Access**: Ensure the machine running the code has internet access to query Docker registries.

---

## Installation

1. Clone the repository:

   ```bash
    git clone https://github.com/robinmuhia/helm-charts.git
    cd helm-charts
   ```

2. Install dependencies

   ```bash
   go mod tidy
   ```

3. Run the project (Populate environemnt variables from env.example)
   ```bash
   go run server.go
   ```

## Usage

### API Endpoint

**POST** `/api/v1/helm-link`  
**Content-Type:** `application/json`

#### Request Example (cURL)

````bash
curl -X POST http://localhost:8080/api/v1/helm-link \
-H "Content-Type: application/json" \
-d '{
  "url_link": "https://github.com/helm/examples/releases/download/hello-world-0.1.0/hello-world-0.1.0.tgz"
}'

2. Expected Response

   ```bash
   [
    {
        "image": "nginx:1.16.0",
        "size": 44815103,
        "layers": 3
    }
   ]
````

3. In case of an error

   ```bash
   {
    "error": "failed to download Helm chart: received status code 404"
   }
   ```

## Linting and Testing

1. To lint

```bash
   golangci-lint run -c .golangci.yaml
```

2. To run tests

```bash
   go test ./...
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
```

## Deployment

1. You can build the binary and run it as follows;

```bash
go build -o helm-chart-processor
./helm-chart-processor
```

2. Optionally, used docker

```bash
docker build -t helm-chart-processor:latest .
docker-compose up --build
```

## Observability

View traces on

```bash
http://localhost:16686
```
