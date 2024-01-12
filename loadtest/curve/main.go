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
		CurvePoints: []curve_tester.NextCurvePoint{
			{Seconds2TargetRPS: 5, TargetRPS: 3},
			{Seconds2TargetRPS: 5, TargetRPS: 1},
			{Seconds2TargetRPS: 5, TargetRPS: 3},
			{Seconds2TargetRPS: 5, TargetRPS: 3},
			{Seconds2TargetRPS: 5, TargetRPS: 0},
		},
		Target: "192.168.178.98:31153",
		Paths: []string{
			"/",
		},
	}

	client := net.NewTcpClient()

	tester := curve_tester.NewCurveTester(*config, client)
	tester.Run()
}
