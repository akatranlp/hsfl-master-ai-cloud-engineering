FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./src/benchmark/go.mod ./src/benchmark/go.sum  ./
RUN go mod download

COPY ./src/benchmark .
RUN CGO_ENABLED=0 GOOS=linux go build -o /benchmark

# gcr.io/distroless/base-debian11
FROM debian:11-slim 

WORKDIR /

COPY --chown=1000:1000 --from=builder-go /benchmark /benchmark

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/benchmark"]
