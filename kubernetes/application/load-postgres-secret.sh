#! /usr/bin/env bash

set -e
set -o allexport; source .env; set +o allexport

kubectl create secret generic db-secret \
    --from-literal=POSTGRES_USER=$POSTGRES_USER \
    --from-literal=POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
    --from-literal=POSTGRES_DB=$POSTGRES_DB \
    --namespace=hsfl-fape2866

