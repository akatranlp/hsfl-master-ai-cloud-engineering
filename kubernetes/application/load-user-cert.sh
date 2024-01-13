#! /usr/bin/env bash

set -e

kubectl create secret generic user-jwt-secret \
    --from-file src/user-service/certs/key.pem \
    --from-file src/user-service/certs/public.pem \
    --namespace=hsfl-verse-vault
