#! /usr/bin/env bash

set -e

kubectl create secret generic user-jwt-secret \
    --from-file src/user-service/certs/access-key.pem \
    --from-file src/user-service/certs/access-public.pem \
    --from-file src/user-service/certs/refresh-key.pem \
    --from-file src/user-service/certs/refresh-public.pem \
    --namespace=hsfl-verse-vault
