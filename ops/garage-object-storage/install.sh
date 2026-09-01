#!/bin/bash

set -e


required_commands=("kubectl" "git" "helm")
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
echo "ok"

if [ ! -f values.override.yaml ]; then 
    echo "err: values.override.yaml is not found in $(pwd)"
    exit 1
fi
VALUES_FILE_ABS="$(realpath values.override.yaml)"

CLONE_DIR="/tmp/garage-install"
rm -rf "$CLONE_DIR"
echo "cloning garage repo"
git clone https://git.deuxfleurs.fr/Deuxfleurs/garage "$CLONE_DIR"
echo "cloned."
cd "$CLONE_DIR/script/helm"

echo "installing garage via helm"
helm install --create-namespace --namespace garage garage ./garage -f "$VALUES_FILE_ABS"

echo "waiting for pods to be created..."
for i in $(seq 1 30); do
    POD_COUNT=$(kubectl get pods -n garage -l app.kubernetes.io/name=garage --no-headers 2>/dev/null | wc -l)
    if [ "$POD_COUNT" -ge 1 ]; then
        break
    fi
    sleep 2
done

echo "waiting for pod to get ready"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=garage -n garage --timeout=180s


echo "pods are ready. garage nodes are up, but the cluster layout is NOT configured yet."
echo "garage will not accept S3 traffic until a layout is assigned and applied."
echo ""
kubectl exec --stdin --tty -n garage garage-0 -- /garage status || true
echo ""
echo "----------------------------------------------------------------------"
echo "MANUAL STEP REQUIRED:"
echo "the cluster layout must be assigned per node and applied. This cannot"
echo "be safely fully-automated here because it requires you to read the"
echo "node IDs from 'garage status' above and decide capacity/zone per node."
echo ""
echo "example for each node (replace <node_id> and adjust -z/-c):"
echo "  kubectl exec -n garage garage-0 -- /garage layout assign -z dc1 -c 100G <node_id>"
echo ""
echo "then review and apply:"
echo "  kubectl exec -n garage garage-0 -- /garage layout show"
echo "  kubectl exec -n garage garage-0 -- /garage layout apply --version <N>"
echo "----------------------------------------------------------------------"
echo "done. garage is installed but requires the layout step above before use."
