package orchestrator

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DefaultOrchestrator struct {
	client *client.Client
}

func NewDefaultOrchestrator() *DefaultOrchestrator {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	return &DefaultOrchestrator{client: cli}
}

func (orc *DefaultOrchestrator) Close() {
	orc.client.Close()
}

func (orc *DefaultOrchestrator) StopContainers(containers []string) {
	for _, containerId := range containers {
		if err := orc.client.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
			panic(err)
		}
	}
}

func (orc *DefaultOrchestrator) GetContainerEndpoints(containers []string, networkName string, port int) []*Target {
	endpoints := make([]*Target, len(containers))
	for i, containerId := range containers {
		inspectRes, err := orc.client.ContainerInspect(context.Background(), containerId)
		if err != nil {
			panic(err)
		}

		endpoint, err := url.Parse(fmt.Sprintf("http://%s:%d", inspectRes.NetworkSettings.Networks[networkName].IPAddress, port))

		if err != nil {
			panic(err)
		}
		endpoints[i] = &Target{ContainerId: containerId, Url: endpoint, Health: 0}
	}

	log.Println("Endpoints: ", endpoints[0].ContainerId, endpoints[0].Url, endpoints[0].Health)

	return endpoints
}

func (orc *DefaultOrchestrator) StartContainers(image string, replicas int, networkName string) []string {
	pullResponse, err := orc.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	defer pullResponse.Close()

	io.Copy(os.Stdout, pullResponse)

	var containers []string
	for i := 0; i < replicas; i++ {
		createResponse, err := orc.client.ContainerCreate(context.Background(), &container.Config{
			Image: image,

			// ExposedPorts: nat.PortSet{
			// 	"8080/tcp": struct{}{},
			// },
			Env: os.Environ(),
		}, &container.HostConfig{

			// PortBindings: map[nat.Port][]nat.PortBinding{
			// 	"8080/tcp": {
			// 		{
			// 			HostIP: "",
			// 			HostPort: "8081",
			// 		},
			// 	},
			// },
		}, &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{networkName: {NetworkID: networkName}}}, nil, "")

		if err != nil {
			panic(err)
		}

		if err := orc.client.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		containers = append(containers, createResponse.ID)
	}

	return containers
}
