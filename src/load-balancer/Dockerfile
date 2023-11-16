FROM golang:1.21.3 AS builder-go

WORKDIR /app/build

COPY ./src/load-balancer/go.mod ./src/load-balancer/go.sum  ./
RUN go mod download

COPY ./src/load-balancer .
RUN CGO_ENABLED=0 GOOS=linux go build -o /balance

FROM debian:12

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /balance /balance

EXPOSE 8080

## USER nonroot:nonroot

ENTRYPOINT ["/balance"]
