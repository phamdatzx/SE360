# Docker Compose CORS Configuration Guide

## Overview

Your `docker-compose.yaml` has been updated to include CORS middleware configuration for Traefik, matching your Kubernetes setup.

## Changes Applied

### CORS Middleware Definition

Added CORS headers middleware with the following configuration:

```yaml
# CORS Middleware - matches Kubernetes configuration
- "traefik.http.middlewares.cors-headers.headers.accesscontrolallowmethods=GET,POST,PUT,DELETE,PATCH,OPTIONS"
- "traefik.http.middlewares.cors-headers.headers.accesscontrolallowheaders=Content-Type,Authorization,X-Requested-With,Accept,Origin,X-User-Id,X-User-Role,X-Username"
- "traefik.http.middlewares.cors-headers.headers.accesscontrolalloworiginlist=http://localhost:3000,http://localhost:8080,http://localhost:5173"
- "traefik.http.middlewares.cors-headers.headers.accesscontrolallowcredentials=true"
- "traefik.http.middlewares.cors-headers.headers.accesscontrolexposeheaders=X-User-Id,X-User-Role,X-Username"
- "traefik.http.middlewares.cors-headers.headers.accesscontrolmaxage=3600"
- "traefik.http.middlewares.cors-headers.headers.addvaryheader=true"
```

### Middleware Application

#### Public Routes (`/api/user/public`)
- **Middleware**: CORS only
- **Purpose**: Allow cross-origin requests for public endpoints like login/register

```yaml
- "traefik.http.routers.user-service-public.middlewares=cors-headers"
```

#### Private Routes (`/api/user`)
- **Middleware**: CORS + Authentication (in order)
- **Purpose**: Handle CORS first, then validate authentication

```yaml
- "traefik.http.routers.user-service-private.middlewares=cors-headers,forward-auth"
```

## Configuration Details

### Allowed Origins
Currently configured for local development:
- `http://localhost:3000` - React/Next.js default
- `http://localhost:8080` - Backend default
- `http://localhost:5173` - Vite default (your frontend)

### Allowed Methods
- GET, POST, PUT, DELETE, PATCH, OPTIONS

### Allowed Headers
- Standard: Content-Type, Authorization, X-Requested-With, Accept, Origin
- Custom: X-User-Id, X-User-Role, X-Username (from auth middleware)

### Exposed Headers
- X-User-Id, X-User-Role, X-Username

### Other Settings
- **Credentials**: Enabled (allows cookies and auth headers)
- **Max Age**: 3600 seconds (1 hour preflight cache)
- **Vary Header**: Added for proper caching

## Usage

### Starting the Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f traefik user-service

# Restart after changes
docker-compose restart
```

### Testing CORS

#### Test Preflight Request

```bash
curl -X OPTIONS http://k8stestdomain.com/api/user/public/register \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v
```

Expected response headers:
```
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Methods: GET,POST,PUT,DELETE,PATCH,OPTIONS
Access-Control-Allow-Headers: Content-Type,Authorization,...
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 3600
```

#### Test Actual Request

```bash
curl -X POST http://k8stestdomain.com/api/user/public/register \
  -H "Origin: http://localhost:5173" \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123","name":"Test User"}' \
  -v
```

## Comparison: Docker Compose vs Kubernetes

Both environments now have identical CORS configuration:

| Feature | Docker Compose | Kubernetes |
|---------|---------------|------------|
| CORS Middleware | Traefik labels | `middlewares-cors.yaml` |
| Allowed Origins | localhost:3000, 8080, 5173 | localhost:3000, 8080, 5173 |
| Middleware Order | cors-headers, forward-auth | cors-headers, forward-auth |
| Public Routes | CORS only | CORS only |
| Private Routes | CORS + Auth | CORS + Auth |

## Adding Production Origins

To add production domains, update the `accesscontrolalloworiginlist` label:

```yaml
- "traefik.http.middlewares.cors-headers.headers.accesscontrolalloworiginlist=http://localhost:3000,http://localhost:5173,https://yourdomain.com,https://www.yourdomain.com"
```

Then restart:
```bash
docker-compose restart user-service
```

## Troubleshooting

### CORS Errors Still Occurring

1. **Check Traefik logs**:
   ```bash
   docker-compose logs traefik
   ```

2. **Verify middleware is loaded**:
   ```bash
   # Access Traefik dashboard
   open http://localhost:8081/dashboard/
   # Check HTTP > Middlewares section
   ```

3. **Ensure origin matches exactly**:
   - Frontend URL: `http://localhost:5173`
   - Must match origin in allowlist (no trailing slash)

### Service Not Accessible

1. **Check /etc/hosts**:
   ```bash
   cat /etc/hosts | grep k8stestdomain
   ```
   Should contain: `127.0.0.1 k8stestdomain.com`

2. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

3. **Test Traefik routing**:
   ```bash
   curl http://k8stestdomain.com/api/user/public/health -v
   ```

## Next Steps

1. ✅ CORS middleware configured in docker-compose
2. ✅ Removed CORS from Go application code
3. ✅ Frontend updated to use k8stestdomain.com
4. 🔄 Test the setup with your frontend
5. 📝 Add production domains before deploying

## Related Files

- [docker-compose.yaml](file:///home/phamdatzx/projects/e-commerce/docker-compose.yaml) - Main compose file
- [helm/traefik/middlewares-cors.yaml](file:///home/phamdatzx/projects/e-commerce/helm/traefik/middlewares-cors.yaml) - Kubernetes CORS config
- [helm/traefik/README-CORS.md](file:///home/phamdatzx/projects/e-commerce/helm/traefik/README-CORS.md) - Kubernetes CORS guide
