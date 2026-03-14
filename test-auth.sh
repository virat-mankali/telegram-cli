#!/bin/bash
# Quick test script for tgcli auth

set -e

echo "=== tgcli Auth Test ==="
echo ""
echo "Before running this, get your credentials from https://my.telegram.org/apps"
echo ""

# Check if credentials are set
if [ -z "$TGCLI_APP_ID" ] || [ -z "$TGCLI_APP_HASH" ]; then
    echo "ERROR: TGCLI_APP_ID and TGCLI_APP_HASH environment variables not set"
    echo ""
    echo "Set them like this:"
    echo "  export TGCLI_APP_ID=12345678"
    echo "  export TGCLI_APP_HASH=abcdef0123456789..."
    echo ""
    echo "Or run inline:"
    echo "  TGCLI_APP_ID=12345678 TGCLI_APP_HASH=abcdef... ./test-auth.sh"
    exit 1
fi

echo "Using APP_ID: $TGCLI_APP_ID"
echo "Using APP_HASH: ${TGCLI_APP_HASH:0:10}..."
echo ""

# Build and run
echo "Building tgcli..."
go build -o ./dist/tgcli ./cmd/tgcli/

echo "Running auth..."
./dist/tgcli auth

echo ""
echo "✓ Auth test complete!"
echo "Session saved to ~/.tgcli/session.json"
