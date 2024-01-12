package curve_tester

import (
	"fmt"
	"sync"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

type CurveTester struct {
	config           LoadtestConfig
	client           net.Client
	logHeaderPrinted bool
	startTime        time.Time
	isRunning        bool
}

func NewCurveTester(config LoadtestConfig, client net.Client) *CurveTester {
	return &CurveTester{
		config:    config,
		client:    client,
		isRunning: false,
	}
}

func (tester *CurveTester) Run() {
	if tester.isRunning {
		fmt.Println("Curve load test already running")
		return
	}

	fmt.Println("Starting curve load test")

	tester.startTime = time.Now()

	requestHandler := NewRequestHandler(tester.config.Target, tester.config.Paths, tester.client)

	requestWg := sync.WaitGroup{}
	stopRequestHandlerChan := make(chan struct{})
	rpsChan := make(chan int)

	go requestHandler.Run(stopRequestHandlerChan, rpsChan, &requestWg)

	currentRPS := 0
	for _, nextCurvePoint := range tester.config.CurvePoints {
		tester.interpolate2RPS(currentRPS, nextCurvePoint.TargetRPS, nextCurvePoint.Seconds2TargetRPS, rpsChan)
		currentRPS = nextCurvePoint.TargetRPS
		fmt.Println("Step to next curve point")
	}

	fmt.Println("Stop sending requests")
	close(stopRequestHandlerChan)

	fmt.Println("Wait for outstanding requests")
	requestWg.Wait()

	fmt.Println("Load test done")
}

func (tester *CurveTester) interpolate2RPS(startRps, targetRPS int, seconds2Interpolate int, rpsChan chan int) {
	rpsChange := targetRPS - startRps
	steps := rpsChange
	if steps < 0 {
		steps = steps * -1
	}

	if steps == 0 {
		time.Sleep(time.Duration(seconds2Interpolate) * time.Second)
		return
	}

	rpsIncrement := 0
	if rpsChange < 0 {
		rpsIncrement = -1
	} else if rpsChange > 0 {
		rpsIncrement = 1
	}

	interpolationInterval := time.Duration(float64(seconds2Interpolate)/float64(steps)) * time.Second

	currentRPS := startRps

	for i := 0; i < steps; i++ {
		currentRPS += rpsIncrement

		rpsChan <- int(currentRPS)

		tester.logRPS(targetRPS, int(currentRPS))

		time.Sleep(interpolationInterval)
	}
}

func (tester *CurveTester) logRPS(targetRPS int, currentRPS int) {
	if !tester.logHeaderPrinted {
		fmt.Println("ElapsedTime\tTargetRPS\tCurrentRPS")
		tester.logHeaderPrinted = true
	}
	elapsedTime := time.Since(tester.startTime)
	fmt.Printf("%.1f\t\t%d\t\t%d\n", elapsedTime.Seconds(), targetRPS, currentRPS)
}
