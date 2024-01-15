FROM golang:1.21.4 AS builder-go

WORKDIR /app/build

COPY ./src/load-balancer/go.mod ./src/load-balancer/go.sum  ./
RUN go mod download

COPY ./src/load-balancer .
RUN CGO_ENABLED=0 GOOS=linux go build -o /balance

FROM debian:11-slim

RUN apt-get update && apt-get install -y curl

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /balance /balance

EXPOSE 8080

## USER nonroot:nonroot

ENTRYPOINT ["/balance"]
