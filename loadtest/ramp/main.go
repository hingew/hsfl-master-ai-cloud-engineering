package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/ramp/ramp_tester"
)

func main() {
	fmt.Println("Start des Lasttests")
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	config, err := ramp_tester.ReadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	client := net.NewTcpClient()
	tester := ramp_tester.NewTester(*config, client)
	if err := tester.Run(); err != nil {
		log.Fatal(err)
	}
}
