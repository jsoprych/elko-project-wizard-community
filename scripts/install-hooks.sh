#!/bin/bash
# install-hooks.sh — installs git pre-commit hook for elko-project-wizard
set -e
HOOK=".git/hooks/pre-commit"
cp scripts/hooks/pre-commit "$HOOK"
chmod +x "$HOOK"
echo "✓ Pre-commit hook installed: $HOOK"
echo "  Tests will run automatically before every commit."
