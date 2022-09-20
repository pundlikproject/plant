export VERSION=0.0.1

echo "Building schedular $VERSION"
go build -o build cmd/api.go

docker build -t sarafdarpundlik/schedular:$VERSION -f build/Dockerfile .

echo "Removing image to minikube"
minikube image rm sarafdarpundlik/schedular:$VERSION

echo "Pushing image to minikube"
minikube image load sarafdarpundlik/schedular:$VERSION
