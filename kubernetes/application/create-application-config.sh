#! /usr/bin/env bash

set -e
set -o allexport; source .env; set +o allexport

if [ "$GRPC_COMMUNICATION" = "false" ]; then
    AUTH_SERVICE_ENDPOINT=http://user-service:8080/validate-token
    BOOK_SERVICE_ENDPOINT=http://book-service:8080/valdiate-chapter-id
    USER_SERVICE_ENDPOINT=http://user-service:8080/move-user-amount
    TRANSACTION_SERVICE_ENDPOINT=http://transaction-service:8080/check-chapter-bought
else
    AUTH_SERVICE_ENDPOINT=http://user-service:8081
    BOOK_SERVICE_ENDPOINT=http://book-service:8081
    USER_SERVICE_ENDPOINT=http://user-service:8081
    TRANSACTION_SERVICE_ENDPOINT=http://transaction-service:8081
fi

kubectl create configmap application-config \
    --from-literal=AUTH_SERVICE_ENDPOINT=$AUTH_SERVICE_ENDPOINT \
    --from-literal=BOOK_SERVICE_ENDPOINT=$BOOK_SERVICE_ENDPOINT \
    --from-literal=USER_SERVICE_ENDPOINT=$USER_SERVICE_ENDPOINT \
    --from-literal=TRANSACTION_SERVICE_ENDPOINT=$TRANSACTION_SERVICE_ENDPOINT \
    --from-literal=GRPC_COMMUNICATION=$GRPC_COMMUNICATION \
    --from-literal=AUTH_IS_ACTIVE=$AUTH_IS_ACTIVE \
    --from-literal=RESET_ON_INIT=$RESET_ON_INIT \
    --from-literal=JWT_ACCESS_TOKEN_EXPIRATION=$JWT_ACCESS_TOKEN_EXPIRATION \
    --from-literal=JWT_REFRESH_TOKEN_EXPIRATION=$JWT_REFRESH_TOKEN_EXPIRATION \
    --namespace=hsfl-verse-vault
