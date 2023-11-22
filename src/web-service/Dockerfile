###################################
#           builder web           #
###################################
FROM node AS builder-web
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

WORKDIR /app/build

COPY ./src/web-service/package.json ./src/web-service/pnpm-lock.yaml ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

COPY ./src/web-service .
RUN pnpm build


###################################
#           builder go            #
###################################
FROM golang:1.21.4 AS builder-go

WORKDIR /app/build
COPY ./lib/. ../../lib

COPY ./src/web-service/go.mod ./src/web-service/go.sum ./
RUN go mod download

COPY ./src/web-service/main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o /web

###################################
#           runner go             #
###################################
# FROM gcr.io/distroless/base-debian11
FROM debian:11-slim 

WORKDIR /

COPY --chown=1000:1000 --from=builder-go /web /web
COPY --chown=1000:1000 --from=builder-web /app/build/dist /dist

EXPOSE 8080

# USER nonroot:nonroot

ENTRYPOINT ["/web"]
