# Load-Balancer

The loadbalancer creates a number of replicas of your specified image and balances incoming requests among them.

## How to use the Load-Balancer

You can specify which image and how manby replicas through commandline-arguments
like: `--image akatranlp/user-service:latest --replicas 2 --network backend`

The loadbalancer will copy it's environment to the underlying containers, so all env-variables for the service must be specified on the loadbalancer.

### Create Docker-Image

If you want to use an docker-image instead, the following commands must be executed from the root of this project:

```bash
docker build -t load-balancer -f ./src/load-balancer/Dockerfile .
docker run -dit -v /var/run/docker.sock:/var/run/docker.sock -p port:port --envFile envFile load-balancer
```
