package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	loadbalancer "github.com/hingew/hsfl-master-ai-cloud-engineering/lib/load_balancer"
)

func main() {
	image := flag.String("image", "", "")
	replicas := flag.Int("replicas", 1, "")
	network := flag.String("network", "bridge", "")
	port := flag.Int("port", 3000, "")
	flag.Parse()

	containers := StartContainers(*image, *replicas)
	endpoints := EvaluateContainerEndpoint(containers, *network, *port)
	defer StopContainers(containers)

	balancer := loadbalancer.NewLoadBalancer(endpoints, loadbalancer.RoundRobin)

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: balancer,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		server.Shutdown(context.Background())

	}()

	server.ListenAndServe()
}

func EvaluateContainerEndpoint(containerIds []string, networkName string, port int) []*url.URL {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	endpoints := make([]*url.URL, len(containerIds))

	for i, id := range containerIds {
		container, err := cli.ContainerInspect(context.Background(), id)
		if err != nil {
			panic(err)
		}

		ip := container.NetworkSettings.Networks[networkName].IPAddress

		rawUrl := fmt.Sprintf("http://%s:%d", ip, port)
		endpoint, err := url.Parse(rawUrl)
		if err != nil {
			panic(err)
		}

		endpoints[i] = endpoint
	}

	return endpoints
}

func StopContainers(containers []string) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if err := cli.ContainerRemove(context.Background(), container, types.ContainerRemoveOptions{Force: true}); err != nil {
			panic(err)
		}
	}

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

	containers := make([]string, replicas)
	for i := 0; i < replicas; i++ {
		createResponse, err := cli.ContainerCreate(context.Background(), &container.Config{Image: image}, &container.HostConfig{}, nil, nil, "")
		if err != nil {
			panic(err)
		}

		if err := cli.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		containers[i] = createResponse.ID
	}
	return containers
}
