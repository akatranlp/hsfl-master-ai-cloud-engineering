FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib ../../lib

COPY ./src/book-service/go.mod ./src/book-service/go.sum  ./
RUN go mod download

COPY ./src/book-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-service

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /book-service /book-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/book-service"]
