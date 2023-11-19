FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib/. ../../lib

COPY ./src/user-service/go.mod ./src/user-service/go.sum  ./
RUN go mod download

COPY ./src/user-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /user-service /user-service

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/user-service"]
