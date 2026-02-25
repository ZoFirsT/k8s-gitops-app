# Go REST API - GitOps Microservice

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-0088CC?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)

This repository contains the source code for a cloud-native, highly performant REST API built with Go. It serves as the core application for the **[GitOps Infrastructure Project](https://github.com/ZoFirsT/k8s-gitops-infra)**, demonstrating end-to-end CI/CD and Kubernetes deployment best practices.

## Key Features

| Feature | Description |
|---------|-------------|
| **High Performance** | Built with Go and the [Gin](https://gin-gonic.com/) web framework |
| **Micro-Container** | Multi-stage Docker build targeting `scratch` base image; final production image is approximately 10MB |
| **Observability Built-in** | Custom middleware exposes Prometheus metrics for monitoring request rates, durations, and system health |
| **Automated CI/CD** | GitHub Actions pipeline automatically tests, builds multi-architecture images (amd64/arm64), pushes to GHCR, and triggers GitOps updates |

## Project Structure

Following the standard Go project layout:

```
.
├── .github/workflows/   # CI Pipeline (Test -> Build -> Push -> Update Infra)
├── services/api/        # Main API microservice
│   ├── cmd/             # Application entrypoint (main.go)
│   ├── internal/        # Private application and library code
│   │   ├── handler/     # HTTP route handlers (Business logic)
│   │   ├── middleware/  # Gin middlewares (Prometheus interceptors)
│   │   └── model/       # Data structures and DTOs
│   ├── Dockerfile       # Production-ready multi-stage Dockerfile
│   └── go.mod           # Go dependencies
└── docker-compose.yml   # Local development environment (API + Redis)
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Liveness and Readiness probe endpoint for Kubernetes |
| `GET` | `/metrics` | Exposes Prometheus metrics (e.g., `http_requests_total`) |
| `POST` | `/api/v1/tasks` | Creates a new task and caches it in Redis. Requires JSON body |

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "GitOps Task", "description": "Deploying to K8s"}'
```

## Local Development

To run the application locally without a Kubernetes cluster, you can use Docker Compose. This will spin up the Go API and a local Redis instance.

```bash
# Clone the repository
git clone https://github.com/ZoFirsT/k8s-gitops-app.git
cd k8s-gitops-app

# Start the services
docker compose up --build
```

The API will be available at `http://localhost:8080`.

## CI/CD Workflow

This repository uses GitHub Actions as the Continuous Integration engine. Upon pushing to the main branch, the following automated steps occur:

| Stage | Description |
|-------|-------------|
| **Test** | Provisions a Go environment, runs `go test`, and generates coverage reports |
| **Build & Push** | Uses Docker Buildx to build multi-platform images (`linux/amd64`, `linux/arm64`) and pushes them to GitHub Container Registry (`ghcr.io`) |
| **Update Infra** | Commits the new Docker image tag (using the short Git SHA) directly to the `k8s-gitops-infra` repository, which subsequently triggers ArgoCD to sync the new deployment |

---

*Developed by **ZoFirsT***
