kubectl delete -f client-ui-deployment.yaml
kubectl delete -f collector-deployment.yaml

kubectl apply -f client-ui-deployment.yaml
kubectl apply -f collector-deployment.yaml