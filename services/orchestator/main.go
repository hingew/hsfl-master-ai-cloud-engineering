package main

import (
	"flag"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/orchestrator"
)

func main() {
	image := flag.String("image", "", "")
	replicas := flag.Int("replicas", 1, "")
	network := flag.String("network", "bridge", "")
	port := flag.Int("port", 3000, "")
	flag.Parse()

	orchestrator, err := orchestrator.NewOrchestrator(*image, *replicas, *network, *port)
	if err != nil {
		panic(err)
	}

	orchestrator.Start()
}
