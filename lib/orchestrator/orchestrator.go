package orchestrator

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
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

type Orchestrator struct {
	image    string
	replicas int
	network  string
	port     int
	client   *client.Client
}

func NewOrchestrator(image string, replicas int, network string, port int) (*Orchestrator, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	orc := Orchestrator{
		image:    image,
		replicas: replicas,
		network:  network,
		port:     port,
		client:   cli,
	}

	return &orc, nil
}

func (orc *Orchestrator) Start() {
	containerIds, err := orc.startContainers()
	if err != nil {
		panic(err)
	}

	containerEndpoints, err := orc.evaluateContainerEndpoints(*containerIds)
	if err != nil {
		panic(err)
	}

	defer orc.stopContainers(*containerIds)

	weights := make([]int, orc.replicas)
	for index, _ := range weights {
		weights[index] = rand.Intn(orc.replicas)
	}

	balancer_algorithm_tag := os.Getenv("BALANCER_ALGORITHM")
	balancer_algorithm := loadbalancer.RoundRobin
	switch balancer_algorithm_tag {
	case "RoundRobin":
		balancer_algorithm = loadbalancer.RoundRobin
		break
	case "WeightedRoundRobin":
		balancer_algorithm = loadbalancer.WeightedRoundRobin
		break
	case "IPHash":
		balancer_algorithm = loadbalancer.IPHash
		break
	default:
		log.Print("No balancing algorithm provided. Use RoundRobin as default")
	}

	balancer := loadbalancer.NewLoadBalancer(containerEndpoints, weights, balancer_algorithm)

	addr := fmt.Sprintf(":%d", orc.port)
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

func (orc *Orchestrator) evaluateContainerEndpoints(containerIds []string) ([]*url.URL, error) {
	endpoints := make([]*url.URL, len(containerIds))

	for i, id := range containerIds {
		container, err := orc.client.ContainerInspect(context.Background(), id)
		if err != nil {
			return nil, err
		}

		ip := container.NetworkSettings.Networks[orc.network].IPAddress

		rawUrl := fmt.Sprintf("http://%s:%d", ip, orc.port)
		endpoint, err := url.Parse(rawUrl)
		if err != nil {
			return nil, err
		}

		endpoints[i] = endpoint
	}

	return endpoints, nil
}

func (orc *Orchestrator) startContainers() (*[]string, error) {
	pullResponse, err := orc.client.ImagePull(context.Background(), orc.image, types.ImagePullOptions{})
	if err != nil {
		log.Print("Could not pull images")
		return nil, err
	}
	defer pullResponse.Close()

	io.Copy(os.Stdout, pullResponse)

	containerIds := make([]string, orc.replicas)
	for i := 0; i < orc.replicas; i++ {
		createResponse, err := orc.client.ContainerCreate(context.Background(), &container.Config{Image: orc.image}, &container.HostConfig{}, nil, nil, "")
		if err != nil {
			return nil, err
		}

		if err := orc.client.ContainerStart(context.Background(), createResponse.ID, types.ContainerStartOptions{}); err != nil {
			return nil, err
		}

		containerIds[i] = createResponse.ID
	}
	return &containerIds, err
}

func (orc *Orchestrator) stopContainers(containerIds []string) {
	for _, id := range containerIds {
		if err := orc.client.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true}); err != nil {
			panic(err)
		}
	}
}
