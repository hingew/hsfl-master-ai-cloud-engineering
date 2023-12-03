package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/load"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

func main() {
	fmt.Println("Start des Lasttests")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	config, err := load.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	client := net.NewTcpClient()
	tester := load.NewTester(*config, client)
	if err := tester.Run(); err != nil {
		log.Fatal(err)
	}
}
