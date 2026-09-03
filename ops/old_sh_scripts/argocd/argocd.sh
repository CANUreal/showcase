#!/bin/bash

set -e

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
kubectl apply --server-side -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
echo "Default config applied."
kubectl wait --for=condition=available --timeout=120s deployment/argocd-server -n argocd

echo "applying argocd ingress"
kubectl apply -f "$(dirname "$0")/argocd-ingress.yaml"

source "$(dirname "$0")/../common/setup-firewall.sh"

echo "done!"
echo "make sure 192.168.0.101 argocd.local is in /etc/hosts on machines that need access"
echo ""
echo "get admin password with: "
echo "  kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d"



