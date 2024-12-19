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
