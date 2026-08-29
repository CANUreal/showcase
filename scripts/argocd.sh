#!/bin/bash

echo "Hello, let's install argocd!"

required_commands=("kubectl" "curl")

for cmd in "${required_commands[@]}"; do 
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "error: $cmd is not installed"
        echo "go install $cmd before running the script."
        exit 1
    fi
done

echo "Dependencies are installed let's check if kubernetes is installed."
if ! kubectl cluster-info >/dev/null 2>&1; then
    echo "kubernetes cluster is not reachable"
    exit 1
fi
echo "kubernetes check, deps check"
echo "Let's start!!"
echo "[NAMESPACE] Creating namespace"
kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
echo "[NAMESPACE] Namespace created."
echo "Applying default..."
kubectl apply \
    -n argocd \
    -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
echo "Default config applied."
kubectl wait --for=condition=available --timeout=120s deployment/argocd-server -n argocd
echo "Port forwarding to :8080"
kubectl port-forward -n argocd svc/argocd-server 8080:443 >/tmp/argocd-port-forward.log 2>&1 &

PORT_FORWARD_PID=$!

echo "Port forward proc id: $PORT_FORWARD_PID"
echo "Ok, you can use argocd on http://localhost:8080"
