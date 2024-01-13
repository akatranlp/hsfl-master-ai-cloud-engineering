# Web-Service

The Web-Service provides the front-end.
It uses react, combined with tailwind css, shadcn-ui, tanstack-query, react-md-editor, react-hot-toast and vite to create an appealing frontend.

## How to use the Web-Service

The web-service needs to be deployed behind a reverse-proxy with all the correct paths and all services to be running.
These paths can be found in the following config-file: `./src/reverse-proxy/config-sample.yaml`

When this is done the frontend must be build with the following commands:

```bash
pnpm install
pnpm build
```

and can then be started with: `go run main.go` when the environment variable `PORT=<port>` is set.

### Create Docker-Image

If you want to use an docker-image instead, the following commands must be executed from the root of this project:

```bash
docker build -t web-service -f ./src/web-service/Dockerfile .
docker run -dit -p port:port --envFile envFile web-service
```

### Develop the web-service

For development we still need the backend and reverse-proxy to be running.

Then we can run the web-service with `pnpm dev` and use vite's hot-reloading while developing.
