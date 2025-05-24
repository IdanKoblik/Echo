set -e

EXCLUDED="fixtures"
PACKAGES=$(go list ./... | grep -v "$EXCLUDED")

cd web 
npm i
npm run build 
cd ..

echo "Cleaning..."
go clean

echo "Running format and vet..."
go fmt ./...
go vet $PACKAGES

echo "Building..."
go build -ldflags "-s -w" -v

echo "Build completed successfully!"