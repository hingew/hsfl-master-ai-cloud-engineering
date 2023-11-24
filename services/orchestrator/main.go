package main

import (
	"os"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/orchestrator"
)


func main() {
	dockerImage := os.Getenv("DOCKER_IMAGE")
	orchestrator, err := orchestrator.NewOrchestrator(dockerImage, 6, "bridge", 3000)
	if err != nil {
		panic(err)
	}

	orchestrator.Start()
}
