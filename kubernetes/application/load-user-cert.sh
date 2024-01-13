#! /usr/bin/env bash

set -e

kubectl create secret generic user-jwt-secret \
    --from-file src/user-service/certs/id_rsa.pub \
    --from-file src/user-service/certs/id_rsa \
    --namespace=hsfl-verse-vault
