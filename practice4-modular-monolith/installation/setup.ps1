$ErrorActionPreference = "Stop"

Set-Location "$PSScriptRoot\.."

Write-Host "Starting setup for Windows..."

$minikubeStatus = minikube status 2>&1
if ($minikubeStatus -notmatch "Running") {
    Write-Host "Starting minikube..."
    minikube start
} else {
    Write-Host "Minikube is already running."
}

Write-Host "Applying Kubernetes manifests from k8s/ folder..."
kubectl apply -f k8s/

Write-Host "Waiting for gateway deployment to be ready..."
kubectl wait --for=condition=available --timeout=120s deployment/gateway -n saga-system

Write-Host "Setup complete."
Write-Host "API will be available at http://localhost:8080"
Write-Host "Starting port-forwarding. Keep this terminal window open."

kubectl port-forward svc/gateway-service 8080:8080 -n saga-system