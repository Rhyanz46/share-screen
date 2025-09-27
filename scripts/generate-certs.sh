#!/bin/bash

# Certificate generation script for share-screen
# This script generates self-signed certificates for development/testing
# For production, use certificates from a trusted CA

set -euo pipefail

# Configuration
CERTS_DIR="./certs"
DAYS_VALID=365
KEY_SIZE=2048
COUNTRY="ID"
STATE="Jakarta"
CITY="Jakarta"
ORG="ShareScreen"
OU="IT Department"

# Get local IP for SAN
LOCAL_IP=$(hostname -I | awk '{print $1}' 2>/dev/null || echo "127.0.0.1")

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if openssl is installed
if ! command -v openssl &> /dev/null; then
    print_error "OpenSSL is not installed. Please install it first."
    exit 1
fi

# Create certs directory
print_info "Creating certificates directory: $CERTS_DIR"
mkdir -p "$CERTS_DIR"

# Generate private key
print_info "Generating private key..."
openssl genrsa -out "$CERTS_DIR/server.key" $KEY_SIZE

# Set secure permissions for private key
chmod 600 "$CERTS_DIR/server.key"

# Create certificate configuration
cat > "$CERTS_DIR/server.conf" << EOF
[req]
default_bits = $KEY_SIZE
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = v3_req

[dn]
C = $COUNTRY
ST = $STATE
L = $CITY
O = $ORG
OU = $OU
CN = localhost

[v3_req]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
IP.3 = $LOCAL_IP
EOF

# Generate certificate signing request
print_info "Generating certificate signing request..."
openssl req -new -key "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.csr" -config "$CERTS_DIR/server.conf"

# Generate self-signed certificate
print_info "Generating self-signed certificate (valid for $DAYS_VALID days)..."
openssl x509 -req -in "$CERTS_DIR/server.csr" -signkey "$CERTS_DIR/server.key" -out "$CERTS_DIR/server.crt" -days $DAYS_VALID -extensions v3_req -extfile "$CERTS_DIR/server.conf"

# Set appropriate permissions
chmod 644 "$CERTS_DIR/server.crt"

# Clean up
rm "$CERTS_DIR/server.csr" "$CERTS_DIR/server.conf"

# Display certificate information
print_info "Certificate generated successfully!"
echo ""
print_info "Certificate details:"
openssl x509 -in "$CERTS_DIR/server.crt" -text -noout | grep -A 5 "Subject:"
echo ""
openssl x509 -in "$CERTS_DIR/server.crt" -text -noout | grep -A 10 "Subject Alternative Name:"
echo ""

print_info "Files created:"
echo "  - $CERTS_DIR/server.key (private key)"
echo "  - $CERTS_DIR/server.crt (certificate)"
echo ""

print_warning "This is a self-signed certificate for development/testing only!"
print_warning "Browsers will show security warnings. For production, use certificates from a trusted CA."
echo ""

print_info "To trust this certificate in your browser:"
echo "  1. Open the HTTPS URL in your browser"
echo "  2. Click 'Advanced' when you see the security warning"
echo "  3. Click 'Proceed to localhost (unsafe)'"
echo ""
echo "  Or add the certificate to your system's trusted certificates store."
echo ""

print_info "Certificate generation completed successfully!"