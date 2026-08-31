#!/bin/bash

set -e
echo "installing pg"

required_commands=("kubectl" "curl")

for cmd in "${required_commands[@]}"; do 
    if ! command -v "$cmd" >/dev/null 2>&1; then
        echo "error: $cmd is not installed"
        echo "go install $cmd before running the script."
        exit 1
    fi
done

echo "dependencies are ok!"
# left here 