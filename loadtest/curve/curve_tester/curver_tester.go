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
	rpsChan          chan int
	currentRPS       int
	targetRPS        int
}

func NewCurveTester(config LoadtestConfig, client net.Client) *CurveTester {
	return &CurveTester{
		config:     config,
		client:     client,
		isRunning:  false,
		currentRPS: 0,
		targetRPS:  0,
	}
}

func (tester *CurveTester) Run() {
	if tester.isRunning {
		fmt.Println("Curve load test already running")
		return
	}

	requestHandler := NewRequestHandler(tester.config.Target, tester.config.Paths, tester.client)

	requestWg := sync.WaitGroup{}
	stopRequestHandlerChan := make(chan struct{})
	tester.rpsChan = make(chan int, tester.maxTargetRPSChange())

	fmt.Println("Starting curve load test")

	tester.startTime = time.Now()

	go requestHandler.Run(stopRequestHandlerChan, tester.rpsChan, &requestWg)
	go tester.doLogs(stopRequestHandlerChan)

	tester.currentRPS = 0
	for _, nextCurvePoint := range tester.config.CurvePoints {
		tester.targetRPS = nextCurvePoint.TargetRPS
		tester.interpolate2RPS(nextCurvePoint.Seconds2TargetRPS)
		tester.currentRPS = tester.targetRPS
	}

	close(stopRequestHandlerChan)
	fmt.Println("Stopped sending requests")

	fmt.Println("Wait for outstanding requests")
	requestWg.Wait()

	fmt.Println("Load test done")
}

func (tester *CurveTester) interpolate2RPS(seconds2Interpolate int) {

	steps := tester.calcInterpolationSteps()

	if steps == 0 {
		time.Sleep(time.Duration(seconds2Interpolate) * time.Second)
		return
	}

	interpolationInterval := tester.calcInterpolationIntervalDuration(seconds2Interpolate, steps)
	rpsIncrement := tester.calcRpsIncrement()

	for i := 0; i < steps; i++ {
		tester.currentRPS += rpsIncrement

		select {
		case tester.rpsChan <- tester.currentRPS:
		default:
		}

		time.Sleep(interpolationInterval)
	}
}

func (tester *CurveTester) doLogs(stopLogs chan struct{}) {
	fmt.Println("ElapsedTime\tTargetRPS\tCurrentRPS")

	for {
		select {
		case <-stopLogs:
			return
		default:
			elapsedTime := time.Since(tester.startTime)
			fmt.Printf("%.0f\t\t%d\t\t%d\n", elapsedTime.Seconds(), tester.targetRPS, tester.currentRPS)
			time.Sleep(time.Second)
		}
	}
}

func (tester *CurveTester) calcRpsIncrement() int {
	change := tester.targetRPS - tester.currentRPS

	if change < 0 {
		return -1
	} else if change > 0 {
		return 1
	}

	return 0
}

func (tester *CurveTester) calcInterpolationSteps() int {
	change := tester.targetRPS - tester.currentRPS

	if change < 0 {
		change = change * -1
	}

	return change
}

func (tester *CurveTester) calcInterpolationIntervalDuration(seconds2Interpolate int, steps int) time.Duration {
	durationPerStepSeconds := float64(seconds2Interpolate) / float64(steps)
	durationPerStepMilliseconds := int(durationPerStepSeconds * 1000)
	return time.Duration(durationPerStepMilliseconds) * time.Millisecond
}

func (tester *CurveTester) maxTargetRPSChange() int {
	maxTargetRPSChange := 0
	for i := 0; i < len(tester.config.CurvePoints)-1; i++ {
		curvePoint := tester.config.CurvePoints[i]
		nextCurvePoint := tester.config.CurvePoints[i+1]

		rpsChange := nextCurvePoint.TargetRPS - curvePoint.TargetRPS
		if rpsChange < 0 {
			rpsChange = rpsChange * -1
		}

		if rpsChange > maxTargetRPSChange {
			maxTargetRPSChange = rpsChange
		}
	}

	return maxTargetRPSChange
}
