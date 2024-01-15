#! /usr/bin/env bash

set -e
set -o allexport; source .env; set +o allexport

kubectl create secret generic test-data-secret \
    --from-literal=TEST_DATA_USER_PASSWORD=$TEST_DATA_USER_PASSWORD \
    --namespace=hsfl-verse-vault

kubectl create configmap test-data-sql-config \
    --from-file src/test-data-service/init.sql \
    --namespace=hsfl-verse-vault
