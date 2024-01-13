# Reverse-Proxy

The Reverse-Proxy is only used when the project is run locally.
It forwards requests to the servers defined in the config.

## How to use the Reverse-Proxy

- create your own config.yaml from config-example.yaml
- specify the port and configFilePath via Environment-variable:

```bash
PORT=<port>
CONFIG_FILE_PATH=<config-Path>
```

- execute the reverse-proxy with `go run main.go`

### Create Docker-Image

If you want to use an docker-image instead, the following commands must be executed from the root of this project:

```bash
docker build -t reverse-proxy -f ./src/reverse-proxy/Dockerfile .
docker run -dit -v configPath:configPath -p port:port --envFile envFile load-balancer
```
