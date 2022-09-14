export VERSION=0.0.1

echo "Building collector $VERSION"
go build -o build cmd/api.go

docker build -t sarafdarpundlik/collector:$VERSION -f build/Dockerfile .

echo "Pushing image to minikube"
minikube image load sarafdarpundlik/collector:$VERSION
