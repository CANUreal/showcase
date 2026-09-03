#!/usr/bin/env bash

# ATTENTION!! 
# THIS IS ONLY INSTALLING SEALED SECRETS
# WE WILL NEED KUBESEAL
# WHEN WE WANTED TO PLAY WITH SECRETS

set -euo pipefail

SEALED_SECRETS_VERSION="v0.39.1"
NAMESPACE="kube-system"

echo "Installing Sealed Secrets ${SEALED_SECRETS_VERSION}..."

kubectl apply -f \
  "https://github.com/bitnami/sealed-secrets/releases/download/${SEALED_SECRETS_VERSION}/controller.yaml"

echo "Waiting for Sealed Secrets controller..."

kubectl rollout status \
  deployment/sealed-secrets-controller \
  -n "$NAMESPACE" \
  --timeout=120s

echo "sealed secrets controller is ready."
