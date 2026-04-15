#!/bin/bash

set -eo pipefail

cd "$(dirname "$0")/.."

COMMIT_HASH=$(git rev-parse --short HEAD || echo "unknown")
echo "Current Commit Hash: $COMMIT_HASH"

echo "Deleting existing Minikube cluster..."
minikube delete


echo "Starting fresh Minikube instance..."
minikube start

echo "Injecting variables and applying Kubernetes manifests..."

for file in k8s/*.yaml; do
    sed "s/{{COMMIT_HASH}}/$COMMIT_HASH/g" "$file" | kubectl apply -f -
done

echo "Waiting for services to initialize..."
sleep 5

echo "Waiting for gateway deployment to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/gateway -n saga-system

echo "API Gateway: http://localhost:8080"
echo "Starting port-forwarding... (Press Ctrl+C to stop)"
kubectl port-forward svc/gateway-service 8080:8080 -n saga-system