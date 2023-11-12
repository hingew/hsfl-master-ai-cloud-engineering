package main

import (
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	image := flag.String("image", "", "")
	replicas := flag.Int("replicas", 1, "")
	flag.Parse()

	containers := StartContainers(*image, *replicas)
	defer StopContainers(containers)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
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
