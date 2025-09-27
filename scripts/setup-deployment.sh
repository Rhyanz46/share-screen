#!/bin/bash

# Setup Deployment Script
# This script helps you configure GitHub Actions deployment

set -e

echo "🚀 Share Screen - Deployment Setup"
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running on macOS or Linux
if [[ "$OSTYPE" == "darwin"* ]]; then
    CLIPBOARD_CMD="pbcopy"
elif command -v xclip &> /dev/null; then
    CLIPBOARD_CMD="xclip -selection clipboard"
elif command -v xsel &> /dev/null; then
    CLIPBOARD_CMD="xsel --clipboard --input"
else
    CLIPBOARD_CMD=""
fi

copy_to_clipboard() {
    if [[ -n "$CLIPBOARD_CMD" ]]; then
        echo "$1" | $CLIPBOARD_CMD
        success "Copied to clipboard!"
    else
        warning "Clipboard not available. Please copy manually."
    fi
}

# Main setup
main() {
    echo
    info "This script will help you set up GitHub Actions deployment."
    echo

    # Check if we're in the right directory
    if [[ ! -f "main.go" ]] || [[ ! -d ".github/workflows" ]]; then
        error "Please run this script from the share-screen project root directory."
        exit 1
    fi

    # Generate SSH key
    echo "1. 🔑 SSH Key Generation"
    echo "======================"
    echo

    read -p "Do you want to generate a new SSH key for deployment? (y/N): " generate_key

    if [[ $generate_key =~ ^[Yy]$ ]]; then
        KEY_NAME="share-screen-deploy"
        SSH_DIR="$HOME/.ssh"
        PRIVATE_KEY="$SSH_DIR/$KEY_NAME"
        PUBLIC_KEY="$SSH_DIR/${KEY_NAME}.pub"

        info "Generating SSH key pair..."
        ssh-keygen -t ed25519 -C "github-actions-share-screen-deploy" -f "$PRIVATE_KEY" -N ""

        success "SSH key generated!"
        echo "📁 Private key: $PRIVATE_KEY"
        echo "📁 Public key: $PUBLIC_KEY"
        echo

        info "Here's your public key to add to your server:"
        echo "================================================"
        cat "$PUBLIC_KEY"
        echo "================================================"

        copy_to_clipboard "$(cat "$PUBLIC_KEY")"
        echo

        warning "Add this public key to your server's ~/.ssh/authorized_keys file"
        read -p "Press Enter when you've added the public key to your server..."

        echo
        info "Here's your private key for GitHub Secrets (SSH_PRIVATE_KEY):"
        echo "=============================================================="
        cat "$PRIVATE_KEY"
        echo "=============================================================="

        copy_to_clipboard "$(cat "$PRIVATE_KEY")"
        echo
    fi

    # Collect server information
    echo "2. 🖥️ Server Information"
    echo "======================"
    echo

    read -p "Enter your server IP or hostname: " SERVER_HOST
    read -p "Enter SSH username (e.g., ubuntu, deploy): " SSH_USER
    read -p "Enter deployment path (e.g., /opt/share-screen): " DEPLOY_PATH

    # Optional: Slack webhook
    echo
    echo "3. 📢 Notifications (Optional)"
    echo "=============================="
    echo

    read -p "Enter Slack webhook URL (optional, press Enter to skip): " SLACK_WEBHOOK

    # Display GitHub Secrets configuration
    echo
    echo "4. 🔐 GitHub Secrets Configuration"
    echo "=================================="
    echo

    info "Add these secrets to your GitHub repository:"
    echo "Repository Settings > Secrets and variables > Actions"
    echo

    echo "Required Secrets:"
    echo "=================="
    echo "SSH_PRIVATE_KEY     = [Your SSH private key content]"
    echo "SSH_USER           = $SSH_USER"
    echo "SERVER_HOST        = $SERVER_HOST"
    echo "DEPLOY_PATH        = $DEPLOY_PATH"

    if [[ -n "$SLACK_WEBHOOK" ]]; then
        echo
        echo "Optional Secrets:"
        echo "================="
        echo "SLACK_WEBHOOK      = $SLACK_WEBHOOK"
    fi

    # Create secrets summary file
    SECRETS_FILE="deployment-secrets.txt"
    cat > "$SECRETS_FILE" << EOF
# GitHub Secrets Configuration for Share Screen Deployment
# Add these secrets to: Repository Settings > Secrets and variables > Actions

# Required Secrets
SSH_PRIVATE_KEY=[Your SSH private key content from above]
SSH_USER=$SSH_USER
SERVER_HOST=$SERVER_HOST
DEPLOY_PATH=$DEPLOY_PATH

# Optional Secrets
SLACK_WEBHOOK=$SLACK_WEBHOOK

# For Staging Environment (if using)
STAGING_SSH_PRIVATE_KEY=[Same as SSH_PRIVATE_KEY or different key]
STAGING_SSH_USER=$SSH_USER
STAGING_SERVER_HOST=[Your staging server IP]
STAGING_DEPLOY_PATH=/opt/share-screen-staging
EOF

    success "Secrets configuration saved to $SECRETS_FILE"

    # Server setup commands
    echo
    echo "5. 🛠️ Server Setup Commands"
    echo "=========================="
    echo

    info "Run these commands on your server ($SERVER_HOST):"
    echo

    SERVER_SETUP_COMMANDS="# Install Go (if not installed)
sudo apt update
sudo apt install -y golang-1.23 || {
    wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
    echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc
    source ~/.bashrc
}

# Create deployment directory
sudo mkdir -p $DEPLOY_PATH
sudo chown \$USER:\$USER $DEPLOY_PATH

# Setup firewall (if ufw is enabled)
sudo ufw allow 8080/tcp
sudo ufw allow 8443/tcp

# Verify setup
go version
ls -la $DEPLOY_PATH"

    echo "$SERVER_SETUP_COMMANDS"

    # Create server setup script
    SERVER_SCRIPT="setup-server.sh"
    cat > "$SERVER_SCRIPT" << EOF
#!/bin/bash
# Server Setup Script for Share Screen Deployment
# Run this on your server: $SERVER_HOST

set -e

echo "🛠️ Setting up server for Share Screen deployment..."

$SERVER_SETUP_COMMANDS

echo "✅ Server setup completed!"
echo "🚀 Ready for GitHub Actions deployment!"
EOF

    chmod +x "$SERVER_SCRIPT"
    success "Server setup script created: $SERVER_SCRIPT"

    copy_to_clipboard "$SERVER_SETUP_COMMANDS"

    # Final instructions
    echo
    echo "6. 🎯 Next Steps"
    echo "==============="
    echo

    success "Setup completed! Here's what to do next:"
    echo
    echo "📋 Checklist:"
    echo "□ Add SSH public key to server's ~/.ssh/authorized_keys"
    echo "□ Run server setup commands on your server"
    echo "□ Add GitHub Secrets to your repository"
    echo "□ Push to main branch to trigger deployment"
    echo

    info "Test your SSH connection:"
    echo "ssh -i ~/.ssh/$KEY_NAME $SSH_USER@$SERVER_HOST"
    echo

    info "Manual deployment test:"
    echo "git push origin main"
    echo

    success "Happy deploying! 🚀"
}

# Run main function
main "$@"