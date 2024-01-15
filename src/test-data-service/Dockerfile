FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib ../../lib

COPY ./src/test-data-service/go.mod ./src/test-data-service/go.sum  ./
RUN go mod download

COPY ./src/test-data-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /test-data-service

#FROM gcr.io/distroless/base-debian11
FROM debian:11-slim 

RUN apt-get update && apt-get install -y curl

WORKDIR /

COPY --chown=1000:1000 --from=builder-go /test-data-service /test-data-service

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/test-data-service"]
