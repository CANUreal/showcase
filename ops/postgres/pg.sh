#!/bin/bash

set -e
echo "installing pg"

required_commands=("kubectl" "curl" "helm")

for cmd in "${required_commands[@]}"; do 
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "error: $cmd is not installed"
        echo "go install $cmd before running the script."
        exit 1
    fi
done

echo "dependencies are ok!"

echo "Let's check if kubernetes is installed."
if ! kubectl cluster-info >/dev/null 2>&1; then
    echo "kubernetes cluster is not reachable"
    exit 1
fi
echo "kubernetes check, deps check"

NAMESPACE="postgres"
RELEASE_NAME="showcase-pg"
DB_NAME="showcase"
DB_USER="showcase"

echo "adding bitnami helm repo"
helm repo add bitnami https://charts.bitnami.com/bitnami >/dev/null 2>&1 || true
helm repo update >/dev/null 2>&1

echo "creating namespace"
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

if helm status "$RELEASE_NAME" -n "$NAMESPACE" >/dev/null 2>&1; then
    echo "postgres release already exists, skipping install"
else
    echo "installing postgres via helm"
    # let's complete it here without writing a yaml
    helm install "$RELEASE_NAME" bitnami/postgresql \
        --namespace "$NAMESPACE" \
        --set auth.username="$DB_USER" \
        --set auth.database="$DB_NAME" \
        --set auth.postgresPassword="" \
        --set primary.persistence.size=2Gi
fi

echo "waiting for postgres pod to be ready"
kubectl wait --for=condition=ready pod \
    -l app.kubernetes.io/name=postgresql \
    -n "$NAMESPACE" \
    --timeout=180s

echo ""
echo "----------------------------------------------------------------------"
echo "postgres is up."
echo ""
echo "get the auto-generated password with:"
echo "  kubectl get secret --namespace $NAMESPACE $RELEASE_NAME-postgresql -o jsonpath=\"{.data.password}\" | base64 -d"
echo "ours is: kubectl get -n postgres secret showcase-pg-postgresql -o jsonpath='{.data.postgres-password}' | base64 -d; echo"
echo ""
echo "the service is reachable inside the cluster at:"
echo "  $RELEASE_NAME-postgresql.$NAMESPACE.svc.cluster.local:5432"
echo ""
echo "for local access (e.g. running the go backend outside the cluster), port-forward:"
echo "  kubectl port-forward -n $NAMESPACE svc/$RELEASE_NAME-postgresql 5432:5432"
echo ""
echo "then DATABASE_URL would look like:"
echo "  postgres://$DB_USER:<password>@localhost:5432/$DB_NAME?sslmode=disable"
echo "----------------------------------------------------------------------"

