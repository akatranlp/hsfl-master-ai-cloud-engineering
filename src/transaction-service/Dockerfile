FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib ../../lib

COPY ./src/transaction-service/go.mod ./src/transaction-service/go.sum  ./
RUN go mod download

COPY ./src/transaction-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /transaction-service

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /transaction-service /transaction-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/transaction-service"]
