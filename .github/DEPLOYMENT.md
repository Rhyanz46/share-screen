# 🚀 Deployment Guide

This document explains how to set up the GitHub Actions CI/CD pipeline for automatic deployment to your server.

## 📋 Prerequisites

1. **Server Requirements:**
   - Ubuntu/Debian server with SSH access
   - Go 1.23.3+ installed
   - Systemd for service management
   - Sudo access for the deployment user

2. **GitHub Repository:**
   - Admin access to configure secrets
   - GitHub Container Registry enabled

## 🔐 Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

### Server Access Secrets
```bash
SSH_PRIVATE_KEY     # Your SSH private key for server access
SSH_USER           # Username for SSH connection (e.g., 'deploy', 'ubuntu')
SERVER_HOST        # Server IP address or domain (e.g., '192.168.1.100')
DEPLOY_PATH        # Deployment directory (e.g., '/opt/share-screen')
```

### Optional Secrets
```bash
SLACK_WEBHOOK      # Slack webhook URL for deployment notifications (optional)
```

## 🚀 Quick Setup

### Automated Setup (Recommended)
```bash
# Run the automated setup script
make setup-deployment

# This will:
# 1. Generate SSH keys
# 2. Guide you through server configuration
# 3. Create GitHub Secrets configuration
# 4. Generate server setup commands
```

## 🔑 Setup SSH Access

### 1. Generate SSH Key Pair (if you don't have one)
```bash
# On your local machine
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/github_actions_deploy

# This creates:
# ~/.ssh/github_actions_deploy (private key)
# ~/.ssh/github_actions_deploy.pub (public key)
```

### 2. Copy Public Key to Server
```bash
# Copy public key to your server
ssh-copy-id -i ~/.ssh/github_actions_deploy.pub user@your-server.com

# Or manually add to ~/.ssh/authorized_keys on server
cat ~/.ssh/github_actions_deploy.pub | ssh user@your-server.com "mkdir -p ~/.ssh && cat >> ~/.ssh/authorized_keys"
```

### 3. Add Private Key to GitHub Secrets
```bash
# Copy the private key content
cat ~/.ssh/github_actions_deploy

# Go to GitHub Repository > Settings > Secrets and variables > Actions
# Add new secret: SSH_PRIVATE_KEY with the private key content
```

## 🏗️ Server Setup

### 1. Install Go on Server
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y golang-1.23

# Or download directly
wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 2. Create Deployment Directory
```bash
sudo mkdir -p /opt/share-screen
sudo chown $USER:$USER /opt/share-screen
```

### 3. Setup Firewall (if needed)
```bash
# Allow HTTP and HTTPS ports
sudo ufw allow 8080/tcp
sudo ufw allow 8443/tcp
sudo ufw reload
```

## 📁 GitHub Secrets Configuration

Go to your GitHub repository and configure these secrets:

### Repository Settings > Secrets and variables > Actions

| Secret Name | Example Value | Description |
|-------------|---------------|-------------|
| `SSH_PRIVATE_KEY` | `-----BEGIN OPENSSH PRIVATE KEY-----\n...` | Your SSH private key content |
| `SSH_USER` | `ubuntu` | SSH username for server access |
| `SERVER_HOST` | `192.168.1.100` | Server IP address or domain |
| `DEPLOY_PATH` | `/opt/share-screen` | Directory where app will be deployed |
| `SLACK_WEBHOOK` | `https://hooks.slack.com/...` | Optional: Slack webhook for notifications |

## 🚀 Deployment Process

The CI/CD pipeline runs automatically on:
- **Push to main branch**: Full deployment
- **Push to develop branch**: Test and build only
- **Pull requests**: Test and build only

### Pipeline Stages:

1. **🧪 Test & Quality Check**
   - Run all unit tests (33/33 tests)
   - Check code formatting
   - Run go vet
   - Generate coverage reports

2. **🔒 Security Scan**
   - Run Gosec security scanner
   - Upload security reports

3. **🐳 Build Docker Image**
   - Build multi-platform Docker image
   - Push to GitHub Container Registry
   - Test Docker image

4. **🚀 Deploy to Server** (main branch only)
   - SSH to server
   - Copy application files
   - Build application
   - Setup systemd service
   - Generate SSL certificates
   - Start/restart service
   - Health check

5. **📢 Notify Deployment**
   - Send Slack notification (if configured)

## 🔧 Service Management

Once deployed, your application runs as a systemd service:

```bash
# Check service status
sudo systemctl status share-screen

# View logs
sudo journalctl -u share-screen -f

# Restart service
sudo systemctl restart share-screen

# Stop service
sudo systemctl stop share-screen

# Start service
sudo systemctl start share-screen
```

## 🌐 Access Your Application

After successful deployment:

- **HTTP**: `http://your-server-ip:8080`
- **HTTPS**: `https://your-server-ip:8443`

## 📊 Monitoring

### Application Logs
```bash
# View real-time logs
sudo journalctl -u share-screen -f

# View recent logs
sudo journalctl -u share-screen --since "1 hour ago"
```

### System Resources
```bash
# Check memory and CPU usage
htop

# Check disk usage
df -h

# Check network connections
ss -tlnp | grep share-screen
```

## 🐛 Troubleshooting

### Common Issues

1. **SSH Connection Failed**
   ```bash
   # Test SSH connection manually
   ssh -i ~/.ssh/github_actions_deploy user@server-ip

   # Check SSH key permissions
   chmod 600 ~/.ssh/github_actions_deploy
   ```

2. **Service Failed to Start**
   ```bash
   # Check service logs
   sudo journalctl -u share-screen --no-pager -l

   # Check if port is in use
   sudo ss -tlnp | grep :8080
   ```

3. **Certificate Issues**
   ```bash
   # Regenerate certificates
   cd /opt/share-screen
   ./scripts/generate-certs.sh
   sudo systemctl restart share-screen
   ```

4. **Permission Issues**
   ```bash
   # Fix ownership
   sudo chown -R $USER:$USER /opt/share-screen

   # Make scripts executable
   chmod +x /opt/share-screen/scripts/*.sh
   ```

## 🔄 Manual Deployment

If you need to deploy manually:

```bash
# 1. SSH to server
ssh user@your-server.com

# 2. Go to deployment directory
cd /opt/share-screen

# 3. Pull latest changes (if using git)
git pull origin main

# 4. Build application
go build -o bin/share-screen .

# 5. Restart service
sudo systemctl restart share-screen
```

## 📈 Scaling Considerations

For production scaling:

1. **Load Balancer**: Use nginx or HAProxy
2. **Multiple Servers**: Modify workflow for multiple deployment targets
3. **Database**: Add persistent storage if needed
4. **Monitoring**: Add Prometheus/Grafana monitoring
5. **Backup**: Implement backup strategy

## 🔐 Security Best Practices

1. **Firewall**: Only open necessary ports
2. **SSH**: Use key-based authentication only
3. **User**: Create dedicated deployment user
4. **Certificates**: Use Let's Encrypt for production
5. **Updates**: Keep server and dependencies updated

---

Happy Deploying! 🚀