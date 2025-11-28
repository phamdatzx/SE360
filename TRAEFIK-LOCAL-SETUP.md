# Traefik Local Setup Guide

This docker-compose configuration sets up Traefik as a reverse proxy with forward authentication, matching the Kubernetes Helm configuration.

## Features

✅ **Traefik Reverse Proxy** - Routes traffic to services
✅ **Forward Authentication** - Protects private routes with JWT verification
✅ **Dashboard** - Web UI for monitoring routes and services
✅ **Public & Private Routes** - Separate routing for authenticated and unauthenticated endpoints
✅ **Access Logs** - Track all incoming requests

## Architecture

```
Client Request → Traefik (Port 80) → Route Matching → Forward Auth (for private) → User Service
```

### Routes Configuration

| Route       | Host              | Path                | Authentication | Priority |
| ----------- | ----------------- | ------------------- | -------------- | -------- |
| Public API  | k8stestdomain.com | /api/user/public/\* | ❌ No          | 100      |
| Private API | k8stestdomain.com | /api/user/\*        | ✅ Yes         | 90       |
| Dashboard   | traefik.local     | /                   | ❌ No          | -        |

**Priority Note**: Public routes have higher priority (100) to match before the broader private route (90).

## Setup Instructions

### 1. Configure Local DNS

Add these entries to your `/etc/hosts` file:

```bash
sudo nano /etc/hosts
```

Add:

```
127.0.0.1 k8stestdomain.com
127.0.0.1 traefik.local
```

### 2. Start Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View Traefik logs only
docker-compose logs -f traefik

# View user-service logs only
docker-compose logs -f user-service
```

### 3. Access Services

- **Traefik Dashboard**: http://traefik.local or http://localhost:8081
- **User Service Public API**: http://k8stestdomain.com/api/user/public/*
- **User Service Private API**: http://k8stestdomain.com/api/user/* (requires authentication)

## Testing

### Test Public Endpoints (No Auth Required)

```bash
# Health check or public endpoint
curl http://k8stestdomain.com/api/user/public/verify

# Register
curl -X POST http://k8stestdomain.com/api/user/public/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://k8stestdomain.com/api/user/public/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Test Private Endpoints (Auth Required)

```bash
# Get your token from login response
TOKEN="your-jwt-token-here"

# Access private endpoint
curl http://k8stestdomain.com/api/user/profile \
  -H "Authorization: Bearer $TOKEN"

# Without token (should fail with 401/403)
curl http://k8stestdomain.com/api/user/profile
```

## Forward Authentication Flow

For private routes, Traefik performs forward authentication:

1. **Client** sends request to `http://k8stestdomain.com/api/user/profile`
2. **Traefik** intercepts and forwards auth headers to `http://user-service:8080/api/user/public/verify`
3. **Auth Service** validates JWT and returns:
   - `200 OK` + headers (`X-User-Id`, `X-User-Role`, `X-Username`) → Request proceeds
   - `401/403` → Request rejected
4. **Traefik** forwards the request to user-service with auth headers
5. **User Service** processes the request with user context

## Traefik Middleware Configuration

The forward auth middleware is configured via Docker labels:

```yaml
traefik.http.middlewares.forward-auth.forwardauth.address=http://user-service:8080/api/user/public/verify
traefik.http.middlewares.forward-auth.forwardauth.trustForwardHeader=true
traefik.http.middlewares.forward-auth.forwardauth.authResponseHeaders=X-User-Id,X-User-Role,X-Username
```

## Troubleshooting

### Check Traefik Dashboard

Visit http://traefik.local to see:

- Registered routers and their rules
- Active services
- Middleware configurations

### View Active Routes

```bash
# List all running containers
docker-compose ps

# Check Traefik configuration
docker exec traefik traefik version

# Inspect network
docker network inspect e-commerce_app-network
```

### Common Issues

**Issue**: Cannot access k8stestdomain.com

- **Solution**: Verify `/etc/hosts` entry and clear browser cache

**Issue**: 404 on all routes

- **Solution**: Check that user-service has `traefik.enable=true` label

**Issue**: Forward auth not working

- **Solution**: Verify the `/api/user/public/verify` endpoint returns proper status codes

**Issue**: Dashboard not accessible

- **Solution**: Try http://localhost:8081 as alternative

## Stop and Clean Up

```bash
# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Rebuild and restart
docker-compose up -d --build
```

## Comparison with Kubernetes Setup

| Feature           | Kubernetes (Helm) | Docker Compose (Local) |
| ----------------- | ----------------- | ---------------------- |
| Reverse Proxy     | Traefik Ingress   | Traefik Container      |
| Forward Auth      | Middleware CRD    | Docker Labels          |
| Service Discovery | K8s DNS           | Docker Network         |
| Dashboard         | IngressRoute      | Label + Port           |
| Load Balancing    | Built-in          | Single instance        |
| SSL/TLS           | Cert-manager      | Not configured (local) |

## Next Steps

- Add SSL certificates for local HTTPS testing
- Configure additional services with Traefik routing
- Set up metrics and monitoring (Prometheus)
- Add rate limiting middleware
- Configure CORS middleware

## Additional Resources

- [Traefik Documentation](https://doc.traefik.io/traefik/)
- [Docker Labels Reference](https://doc.traefik.io/traefik/providers/docker/)
- [ForwardAuth Middleware](https://doc.traefik.io/traefik/middlewares/http/forwardauth/)
