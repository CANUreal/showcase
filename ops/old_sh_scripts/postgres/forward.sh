#!/bin/bash
set -e
kubectl port-forward --namespace postgres svc/showcase-pg-postgresql 5432:5432
