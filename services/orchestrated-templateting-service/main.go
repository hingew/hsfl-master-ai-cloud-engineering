package main

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/orchestrator"
)

func main() {
	orchestrator, err := orchestrator.NewOrchestrator("hauing/templateing-service", 6, "bridge", 3000)
	if err != nil {
		panic(err)
	}

	orchestrator.Start()
}
