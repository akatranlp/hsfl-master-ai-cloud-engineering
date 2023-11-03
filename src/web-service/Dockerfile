FROM node AS builder-web
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app/build

COPY ./src/web-service/package.json ./src/web-service/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./src/web-service .
RUN pnpm build

FROM golang:1.21.3 AS builder-go

WORKDIR /app/build

COPY ./src/web-service/go.mod .
RUN go mod download

COPY ./src/web-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o /web

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --chown=nonroot:nonroot --from=builder-go /web /web
COPY --chown=nonroot:nonroot --from=builder-web /app/build/dist /dist

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/web"]
