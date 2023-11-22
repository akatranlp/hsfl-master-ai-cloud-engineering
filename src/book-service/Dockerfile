FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib ../../lib

COPY ./src/book-service/go.mod ./src/book-service/go.sum  ./
RUN go mod download

COPY ./src/book-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-service

# gcr.io/distroless/base-debian11
FROM debian:11-slim 

WORKDIR /

COPY --chown=1000:1000 --from=builder-go /book-service /book-service

EXPOSE 8080

# USER 1000:1000

ENTRYPOINT ["/book-service"]
