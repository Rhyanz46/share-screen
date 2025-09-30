# Deployment Configuration

This document describes the required GitHub secrets for deploying the share-screen application.

## Required GitHub Secrets

⚠️ **IMPORTANT**: All required secrets must be configured or deployment will be cancelled automatically.

### Production Server Connection Secrets (REQUIRED)
- `SSH_USER` - SSH username for production server
- `SSH_PRIVATE_KEY` - SSH private key for production server
- `SERVER_HOST` - Production server hostname/IP
- `DEPLOY_PATH` - Deployment path on production server

### Staging Server Connection Secrets (REQUIRED for staging)
- `STAGING_SSH_USER` - SSH username for staging server
- `STAGING_SSH_PRIVATE_KEY` - SSH private key for staging server
- `STAGING_SERVER_HOST` - Staging server hostname/IP
- `STAGING_DEPLOY_PATH` - Deployment path on staging server

### Production Environment Secrets (Optional - defaults provided)
- `PROD_PORT` - Application port (default: 8080)
- `PROD_HTTP_PORT` - HTTP port mapping (default: 8080)
- `PROD_HTTPS_PORT` - HTTPS port mapping (default: 8443)
- `PROD_ENABLE_HTTPS` - Enable HTTPS (default: true)
- `PROD_STUN_SERVER` - STUN server URL (default: stun:stun.l.google.com:19302)
- `PROD_TOKEN_EXPIRY` - Token expiry duration (default: 30m)
- `PROD_TLS_CERT_FILE` - TLS certificate file path (default: /certs/server.crt)
- `PROD_TLS_KEY_FILE` - TLS private key file path (default: /certs/server.key)

### Staging Environment Secrets (Optional - defaults provided)
- `STAGING_PORT` - Application port (default: 8081)
- `STAGING_HTTP_PORT` - HTTP port mapping (default: 8081)
- `STAGING_HTTPS_PORT` - HTTPS port mapping (default: 8444)
- `STAGING_ENABLE_HTTPS` - Enable HTTPS (default: true)
- `STAGING_STUN_SERVER` - STUN server URL (default: stun:stun.l.google.com:19302)
- `STAGING_TOKEN_EXPIRY` - Token expiry duration (default: 15m)
- `STAGING_TLS_CERT_FILE` - TLS certificate file path (default: /certs/server.crt)
- `STAGING_TLS_KEY_FILE` - TLS private key file path (default: /certs/server.key)

## How to Configure

1. **Go to your GitHub repository**
2. **Settings → Secrets and variables → Actions**
3. **Add the required secrets above**

## Environment Setup

The deployment workflows will:
1. **Use GitHub secrets** for environment variables (no .env file needed in production)
2. **Apply sensible defaults** if optional secrets are not provided
3. **Pass environment variables** directly to Docker Compose

## Secret Validation

The deployment workflows include automatic validation that:

1. **Checks all required secrets** before starting deployment
2. **Fails fast** if any required secret is missing
3. **Shows status** of optional secrets (set vs using defaults)
4. **Prevents partial deployments** due to missing configuration

### Validation Output Example:
```
🔍 Validating required secrets for production deployment...
✅ All required secrets are configured
📊 Optional secrets status:
  - PROD_PORT: ⚠️ Using default (8080)
  - PROD_HTTP_PORT: ✅ Set
  - PROD_HTTPS_PORT: ✅ Set
  - PROD_ENABLE_HTTPS: ⚠️ Using default (true)
```

## Security Benefits

- ✅ **No sensitive data** in repository
- ✅ **Environment-specific** configuration
- ✅ **Encrypted storage** in GitHub
- ✅ **Access control** via GitHub permissions
- ✅ **No .env files** in production servers
- ✅ **Pre-deployment validation** prevents failed deployments

## Docker Compose Integration

The environment variables are passed to Docker Compose, which then passes them to the container as defined in `docker-compose.yml`:

```yaml
environment:
  - PORT=${PORT:-8080}
  - STUN_SERVER=${STUN_SERVER:-stun:stun.l.google.com:19302}
  - TOKEN_EXPIRY=${TOKEN_EXPIRY:-30m}
  - ENABLE_HTTPS=${ENABLE_HTTPS:-false}
  - TLS_CERT_FILE=${TLS_CERT_FILE:-/certs/server.crt}
  - TLS_KEY_FILE=${TLS_KEY_FILE:-/certs/server.key}
```