#!/bin/bash

# Install git hooks for the project
# Run this after cloning the repository

HOOK_DIR=".git/hooks"
HOOK_FILE="$HOOK_DIR/pre-commit"

echo "Installing git hooks..."

# Create hooks directory if it doesn't exist
if [ ! -d "$HOOK_DIR" ]; then
    echo "Error: Not a git repository (no .git/hooks directory found)"
    exit 1
fi

# Create pre-commit hook
cat > "$HOOK_FILE" << 'EOF'
#!/bin/bash

# Git pre-commit hook - runs checks before allowing commit
# This ensures code quality standards are met before code is committed

echo "Running pre-commit checks..."
echo ""

# Run make ci (fmt, vet, staticcheck, tests)
make ci

# Capture the exit code
EXIT_CODE=$?

# If make ci failed, prevent the commit
if [ $EXIT_CODE -ne 0 ]; then
    echo ""
    echo "❌ Pre-commit checks failed!"
    echo "Please fix the errors above before committing."
    echo ""
    echo "To bypass this check (not recommended):"
    echo "  git commit --no-verify"
    exit 1
fi

echo ""
echo "✅ Pre-commit checks passed - proceeding with commit"
exit 0
EOF

# Make hook executable
chmod +x "$HOOK_FILE"

echo "✅ Git hooks installed successfully!"
echo ""
echo "The pre-commit hook will now run 'make ci' before every commit."
echo "This ensures code is formatted, linted, and tested before committing."
echo ""
echo "To bypass the hook (not recommended): git commit --no-verify"
