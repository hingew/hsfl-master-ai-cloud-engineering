package main

import (
	"flag"
	"fmt"
	"loadtest/load"
	"loadtest/net"
	"log"
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
