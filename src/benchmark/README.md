# Benchmark

This is a tool to test if the incoming requests are being equally distributed between the servers and to test the limits of the servers.

## How to use our Benchmark-Tool

- create your own config.yaml from config-example.yaml
- execute the benchmarker

```bash
go run main.go --configPath config.yaml [--waitForResponse]
```
