#!/bin/bash

set -e

echo "installing jenkins on k3s!"
required_commands=("kubectl" "helm")

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

echo "things seem ok lets go"

echo "adding the helm repo"
helm repo add jenkinsci https://charts.jenkins.io
echo "done"
echo "updating repo helm repo list"
helm repo update
echo "done!"

echo "creating namespace"
kubectl create namespace jenkins --dry-run=client -o yaml | kubectl apply -f -

echo "installing jenkins via helm"
helm upgrade --install jenkins jenkinsci/jenkins -n jenkins -f "$(dirname "$0")/jenkins-values.yaml"

echo "waiting jenkins pod to be ready"
kubectl -n jenkins wait --for=condition=ready pod -l app.kubernetes.io/component=jenkins-controller --timeout=180s

# run this if you need remote kubectl access sudo firewall-cmd --zone=public --add-port=6443/tcp --permanent && sudo firewall-cmd --reload
source "$(dirname "$0")/../common/setup-firewall.sh"

echo "done, jenkins is reachable at https://jenkins.local"
echo "but just make sure '192.168.0.101 jenkins.local' is in your /etc/hosts"
echo ""
echo "get admin password with   "
# never hardcode path or pod name!?
echo "  kubectl -n jenkins exec -it \$(kubectl -n jenkins get pods -l app.kubernetes.io/component=jenkins-controller -o jsonpath='{.items[0].metadata.name}') -- cat /run/secrets/additional/chart-admin-password"



