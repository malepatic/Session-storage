# Session Management Backend

## Overview
A cloud-native session management backend built in Go using Redis for token caching with TTL, PostgreSQL for user data persistence, Docker for containerization, and deployment automation with Jenkins and Kubernetes.

## Features
- User registration and authentication with JWT tokens
- Redis-based session caching with TTL
- PostgreSQL data persistence
- RESTful API for session management
- Docker containerization
- Kubernetes deployment configurations
- Automated CI/CD with Jenkins/GitHub Actions
- Health check endpoint

## Tech Stack
- **Backend**: Go, Gin framework
- **Caching**: Redis
- **Database**: PostgreSQL
- **Containerization**: Docker
- **Orchestration**: Kubernetes
- **CI/CD**: Jenkins/GitHub Actions
- **Cloud Provider**: AWS (EKS)

## API Endpoints

| Method | Endpoint   | Description          | Authentication |
|--------|------------|----------------------|----------------|
| POST   | /register  | Register a new user  | No             |
| POST   | /login     | User login           | No             |
| POST   | /logout    | User logout          | Yes            |
| GET    | /profile   | Get user profile     | Yes            |
| GET    | /health    | Health check         | No             |

## Local Development

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Make (optional for using the Makefile)

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/session-app.git
   cd session-app
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

4. Access the API at http://localhost:8080

### Run Tests
```bash
go test ./...
```

## Environment Variables

| Variable          | Description                            | Default                                             |
|-------------------|----------------------------------------|-----------------------------------------------------|
| PORT              | HTTP server port                       | 8080                                                |
| REDIS_URL         | Redis connection URL                   | redis://localhost:6379/0                           |
| DB_URL            | PostgreSQL connection URL              | postgres://postgres:postgres@localhost:5432/session_db?sslmode=disable |
| JWT_SECRET        | Secret key for JWT signing             | your-secret-key                                     |
| TOKEN_EXPIRATION  | Token expiration time in hours         | 24                                                  |

## Deployment

### AWS EKS Deployment
1. Create an EKS cluster:
   ```bash
   eksctl create cluster --name session-cluster --region us-west-2 --node-type t3.medium --nodes 2
   ```

2. Update kubeconfig:
   ```bash
   aws eks update-kubeconfig --name session-cluster --region us-west-2
   ```

3. Apply Kubernetes manifests:
   ```bash
   kubectl apply -f k8s/
   ```

### CI/CD Setup
1. Set up Jenkins or GitHub Actions using the provided configuration files.
2. Configure the required credentials:
   - Docker registry credentials
   - AWS credentials
   - Kubernetes configuration

## Architecture Diagram

```
┌─────────────────┐     ┌───────────────┐     ┌──────────────────┐
│                 │     │               │     │                  │
│  Kubernetes EKS │     │  Redis Cache  │     │  PostgreSQL/RDS  │
│  (Go Service)   │◄────┤  (Sessions)   │     │  (User Data)     │
│                 │     │               │     │                  │
└─────────┬───────┘     └───────────────┘     └──────────────────┘
          │
          │
┌─────────▼───────┐     ┌───────────────┐
│                 │     │               │
│  Load Balancer  │◄────┤   Internet    │
│                 │     │               │
└─────────────────┘     └───────────────┘
```

## Security Considerations
- JWT tokens are stored in Redis with TTL and can be invalidated
- Passwords are hashed using bcrypt
- Kubernetes secrets are used for sensitive data
- Connection information is never hard-coded
- HTTPS should be enabled in production

## Future Enhancements
- Add refresh token functionality
- Implement rate limiting
- Add more authentication methods (OAuth, SSO)
- Implement user roles and permissions
- Add monitoring and logging (Prometheus, ELK)

## License
This project is licensed under the MIT License - see the LICENSE file for details.