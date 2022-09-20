#!/bin/bash
export VERSION=0.0.1

echo "Building collector $VERSION"
go build -o build ./cmd/api.go

echo "Bulding docker"
docker build -t sarafdarpundlik/collector:$VERSION -f build/Dockerfile .

echo "Removing image to minikube"
minikube image rm sarafdarpundlik/collector:$VERSION
echo "Pushing image to minikube"
minikube image load sarafdarpundlik/collector:$VERSION
