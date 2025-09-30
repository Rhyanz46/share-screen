# Deployment Configuration

This document describes the required GitHub secrets for deploying the share-screen application.

## Required GitHub Secrets

⚠️ **IMPORTANT**: All required secrets must be configured or deployment will be cancelled automatically.

### Production Server Connection Secrets (REQUIRED)
- `SSH_USER` - SSH username for production server
- `SSH_PRIVATE_KEY` - SSH private key for production server
- `SERVER_HOST` - Production server hostname/IP
- `DEPLOY_PATH` - Deployment path on production server

### Production Environment Secrets (REQUIRED)
- `PROD_PORT` - Application port (e.g., 8080)
- `PROD_HTTP_PORT` - HTTP port mapping (e.g., 8080)
- `PROD_HTTPS_PORT` - HTTPS port mapping (e.g., 8443)
- `PROD_ENABLE_HTTPS` - Enable HTTPS (true/false)
- `PROD_STUN_SERVER` - STUN server URL (e.g., stun:stun.l.google.com:19302)
- `PROD_TOKEN_EXPIRY` - Token expiry duration (e.g., 30m)
- `PROD_TLS_CERT_FILE` - TLS certificate file path (e.g., /certs/server.crt)
- `PROD_TLS_KEY_FILE` - TLS private key file path (e.g., /certs/server.key)

## How to Configure

1. **Go to your GitHub repository**
2. **Settings → Secrets and variables → Actions**
3. **Add the required secrets above**

## Environment Setup

The deployment workflows will:
1. **Use GitHub secrets** for environment variables (no .env file needed in production)
2. **Apply sensible defaults** if optional secrets are not provided
3. **Pass environment variables** directly to Docker Compose

## Environment Dependencies Check

⚠️ **CRITICAL**: The CI/CD pipeline will FAIL if production environment is missing required secrets.

The CI/CD pipeline includes production environment validation that:

1. **Checks production environment** on every pipeline run
2. **FAILS IMMEDIATELY** if production environment missing required secrets
3. **Reports complete status** of all secrets and configurations
4. **Prevents all pipeline execution** until production environment is properly configured
5. **Enforces complete setup** before any CI/CD operations

### Environment Check Success Example:
```
🔍 Checking production environment dependencies...
📋 Branch: main
📋 Event: push
🎯 Production deployment conditions met, checking secrets...
✅ Production deployment ENABLED - All required secrets configured

🌍 Production Environment Check:
  ✅ All required secrets configured
  - PROD_PORT: ⚠️ Using default (8080)
  - PROD_HTTP_PORT: ✅ Set
  - PROD_HTTPS_PORT: ✅ Set
  - PROD_ENABLE_HTTPS: ⚠️ Using default (true)

🎯 Deployment Summary:
  Production: true

✅ ENVIRONMENT CHECK PASSED
🎯 Production environment has all required secrets configured
```

### Environment Check Failure Example:
```
🔍 Checking production environment dependencies...
📋 Branch: main
📋 Event: push
🎯 Production deployment conditions met, checking secrets...
❌ Production deployment DISABLED - Missing required secrets: SSH_PRIVATE_KEY SERVER_HOST

🌍 Production Environment Check:
  ❌ Missing required secrets

❌ CRITICAL: Production environment missing required secrets!
📋 Required secrets: SSH_USER, SSH_PRIVATE_KEY, SERVER_HOST, DEPLOY_PATH

💥 ENVIRONMENT CHECK FAILED
📖 Please configure all missing secrets in GitHub Settings → Secrets and variables → Actions
🚫 CI/CD pipeline stopped - production environment incomplete
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