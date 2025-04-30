set -e

EXCLUDED="fixtures"
PACKAGES=$(go list ./... | grep -v "$EXCLUDED")

echo "Cleaning..."
go clean

echo "Running format and vet..."
go fmt ./...
go vet $PACKAGES

echo "Building..."
go build -v

echo "Build completed successfully!"