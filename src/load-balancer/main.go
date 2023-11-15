package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/akatranlp/hsfl-master-ai-cloud-engineering/load-balancer/load"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	image := flag.String("image", "", "")
	replicas := flag.Int("replicas", 1, "")
	network := flag.String("network", "bridge", "")
	flag.Parse()

	containers := StartContainers(*image, *replicas)
	defer StopContainers(containers)

	endpoints := GetContainerEndpoints(containers, *network)

	lb := load.NewBalancer(endpoints)

	server := &http.Server{
		Addr:    ":3000",
		Handler: lb,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		server.Shutdown(context.Background())
	}()

	server.ListenAndServe()
}
func StopContainers(containers []string) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	for _, containerId := range containers {
		if err := cli.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			panic(err)
		}
	}
}

func GetContainerEndpoints(containers []string, networkName string) []*url.URL {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	endpoints := make([]*url.URL, len(containers))
	for i, containerId := range containers {
		inspectRes, err := cli.ContainerInspect(context.Background(), containerId)
		if err != nil {
			panic(err)
		}

		endpoint, err := url.Parse(fmt.Sprintf("http://%s:3000", inspectRes.NetworkSettings.Networks[networkName].IPAddress))

		if err != nil {
			panic(err)
		}
		endpoints[i] = endpoint
	}

	return endpoints
}

func StartContainers(image string, replicas int) []string {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	pullResponse, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer pullResponse.Close()

	io.Copy(os.Stdout, pullResponse)

	var containers []string
	for i := 0; i < replicas; i++ {
		createResponse, err := cli.ContainerCreate(context.Background(), &container.Config{
			Image: image,
		}, &container.HostConfig{}, nil, nil, "")

		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		containers = append(containers, createResponse.ID)
	}

	return containers
}
