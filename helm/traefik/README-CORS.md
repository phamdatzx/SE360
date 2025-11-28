# Traefik CORS Configuration Guide

## Overview

This directory contains Traefik middleware configurations for your microservices, including:
- **middlewares-auth.yaml**: Forward authentication middleware
- **middlewares-cors.yaml**: CORS (Cross-Origin Resource Sharing) middleware

## CORS Middleware Configuration

The CORS middleware (`middlewares-cors.yaml`) is configured with:

### Allowed Methods
- GET, POST, PUT, DELETE, PATCH, OPTIONS

### Allowed Headers
- Content-Type, Authorization, X-Requested-With, Accept, Origin
- Custom headers: X-User-Id, X-User-Role, X-Username (matching auth middleware)

### Allowed Origins
Currently configured for local development:
- `http://localhost:3000` (React/Next.js default)
- `http://localhost:8080` (Common backend port)
- `http://localhost:5173` (Vite default)

**Important**: Update the `accessControlAllowOriginList` in `middlewares-cors.yaml` with your production domains.

### Other Settings
- **Credentials**: Enabled (allows cookies and auth headers)
- **Max Age**: 3600 seconds (1 hour cache for preflight requests)
- **Exposed Headers**: X-User-Id, X-User-Role, X-Username

## Deployment

### 1. Apply the CORS Middleware

```bash
kubectl apply -f helm/traefik/middlewares-cors.yaml
```

### 2. Apply the Auth Middleware (if not already done)

```bash
kubectl apply -f helm/traefik/middlewares-auth.yaml
```

### 3. Verify Middlewares

```bash
kubectl get middleware -n default
```

You should see both `cors-headers` and `forward-auth` middlewares.

## Usage in IngressRoute

### Option 1: CORS Only (Public Endpoints)

For public endpoints that don't require authentication:

```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: product-service
  namespace: default
spec:
  entryPoints:
    - web
  routes:
    - match: Host(`api.yourdomain.com`) && PathPrefix(`/api/product`)
      kind: Rule
      services:
        - name: product-service
          port: 8080
      middlewares:
        - name: cors-headers  # CORS only
```

### Option 2: CORS + Authentication (Protected Endpoints)

For protected endpoints that require both CORS and authentication:

```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: product-service-protected
  namespace: default
spec:
  entryPoints:
    - web
  routes:
    - match: Host(`api.yourdomain.com`) && PathPrefix(`/api/product`)
      kind: Rule
      services:
        - name: product-service
          port: 8080
      middlewares:
        - name: cors-headers    # Apply CORS first
        - name: forward-auth    # Then authentication
```

### Option 3: Mixed Routes (Some Public, Some Protected)

```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: product-service-mixed
  namespace: default
spec:
  entryPoints:
    - web
  routes:
    # Public routes - CORS only
    - match: Host(`api.yourdomain.com`) && PathPrefix(`/api/product/public`)
      kind: Rule
      services:
        - name: product-service
          port: 8080
      middlewares:
        - name: cors-headers
    
    # Protected routes - CORS + Auth
    - match: Host(`api.yourdomain.com`) && PathPrefix(`/api/product/admin`)
      kind: Rule
      services:
        - name: product-service
          port: 8080
      middlewares:
        - name: cors-headers
        - name: forward-auth
```

## Customization

### Adding Production Domains

Edit `middlewares-cors.yaml` and add your domains:

```yaml
accessControlAllowOriginList:
  - "http://localhost:3000"
  - "https://yourdomain.com"
  - "https://www.yourdomain.com"
  - "https://app.yourdomain.com"
```

### Adding Custom Headers

If your application uses custom headers, add them to both sections:

```yaml
accessControlAllowHeaders:
  - "Content-Type"
  - "Authorization"
  - "Your-Custom-Header"  # Add here

accessControlExposeHeaders:
  - "X-User-Id"
  - "Your-Custom-Response-Header"  # Add here
```

### Restricting Methods

If you want to restrict HTTP methods:

```yaml
accessControlAllowMethods:
  - "GET"
  - "POST"
  # Remove methods you don't want to allow
```

## Troubleshooting

### CORS Errors Still Occurring

1. **Check middleware is applied**:
   ```bash
   kubectl get middleware cors-headers -n default -o yaml
   ```

2. **Verify IngressRoute configuration**:
   ```bash
   kubectl get ingressroute -n default -o yaml
   ```

3. **Check Traefik logs**:
   ```bash
   kubectl logs -n default -l app.kubernetes.io/name=traefik
   ```

### Preflight Requests Failing

- Ensure OPTIONS method is in `accessControlAllowMethods`
- Check that all required headers are in `accessControlAllowHeaders`
- Verify the origin is in `accessControlAllowOriginList`

### Credentials Not Working

- Ensure `accessControlAllowCredentials: true` is set
- Frontend must use `credentials: 'include'` in fetch/axios
- Origin must be explicitly listed (cannot use wildcard `*` with credentials)

## Testing CORS

### Using curl

```bash
# Test preflight request
curl -X OPTIONS http://api.yourdomain.com/api/product \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

# Test actual request
curl -X GET http://api.yourdomain.com/api/product \
  -H "Origin: http://localhost:3000" \
  -v
```

### Expected Response Headers

You should see these headers in the response:

```
Access-Control-Allow-Origin: http://localhost:3000
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, PATCH, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization, ...
Access-Control-Allow-Credentials: true
Access-Control-Max-Age: 3600
```

## Best Practices

1. **Order matters**: Apply CORS middleware before authentication middleware
2. **Be specific**: List exact origins instead of using wildcards in production
3. **Minimize headers**: Only expose headers that your frontend actually needs
4. **Cache preflight**: Use appropriate `accessControlMaxAge` to reduce preflight requests
5. **Namespace consistency**: Ensure middlewares and IngressRoutes are in the same namespace

## Related Files

- [values.yaml](file:///home/phamdatzx/projects/e-commerce/helm/traefik/values.yaml) - Main Traefik Helm values
- [middlewares-auth.yaml](file:///home/phamdatzx/projects/e-commerce/helm/traefik/middlewares-auth.yaml) - Authentication middleware
- [middlewares-cors.yaml](file:///home/phamdatzx/projects/e-commerce/helm/traefik/middlewares-cors.yaml) - CORS middleware
