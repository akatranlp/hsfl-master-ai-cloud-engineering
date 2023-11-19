FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib ../../lib

COPY ./src/reverse-proxy/go.mod ./src/reverse-proxy/go.sum  ./
RUN go mod download

COPY ./src/reverse-proxy .
RUN CGO_ENABLED=0 GOOS=linux go build -o /proxy

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /proxy /proxy

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/proxy"]
