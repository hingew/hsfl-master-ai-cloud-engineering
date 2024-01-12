package main

import (
	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/curve/curve_tester"
	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

func main() {
	//configPath := flag.String("config", "config.json", "Path to the configuration file")
	//flag.Parse()

	//config, err := curve_tester.ReadConfig(*configPath)
	//if err != nil {
	//	log.Fatalf("Failed to read config: %v", err)
	//}

	config := &curve_tester.LoadtestConfig{
		Graph: []curve_tester.NextGraphPoint{
			{Seconds2TargetRPS: 10, TargetRPS: 10},
			{Seconds2TargetRPS: 5, TargetRPS: 5},
			{Seconds2TargetRPS: 10, TargetRPS: 25},
			{Seconds2TargetRPS: 10, TargetRPS: 0},
		},
		Target: "192.168.2.149:32350",
		Paths: []string{
			"/",
		},
	}

	client := net.NewTcpClient()

	curve_tester.RunLoadTest(*config, client)
}
