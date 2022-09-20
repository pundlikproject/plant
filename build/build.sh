#!/bin/bash
#echo "Building collector"
#sh ../collector/build/build.sh

#echo "Building Schedular"
#sh ../robot-communicator/build/build.sh
#Create namespace
kubectl create namespace plant

#Create and deploy 
kubectl apply -f postgres-configmap.yaml
kubectl apply -f postgres-storage.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml

kubectl delete -f client-ui-deployment.yaml
kubectl delete -f collector-deployment.yaml
kubectl delete -f schedular-deployment.yaml

kubectl apply -f client-ui-deployment.yaml
kubectl apply -f collector-deployment.yaml
kubectl apply -f schedular-deployment.yaml