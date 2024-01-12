package curve_tester

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

type Tester struct {
	config           LoadtestConfig
	client           net.Client
	runningUsers     []*RequestUser
	currentRPS       int
	logHeaderPrinted bool
	startTime        time.Time
	waitGroup        sync.WaitGroup
	errors           chan error
}

func newTester(config LoadtestConfig, client net.Client) *Tester {
	return &Tester{
		config:       config,
		client:       client,
		currentRPS:   0.0,
		waitGroup:    sync.WaitGroup{},
		errors:       make(chan error, 1),
		runningUsers: make([]*RequestUser, 0),
	}
}

func RunLoadTest(config LoadtestConfig, client net.Client) {
	fmt.Println("Starting load test")

	tester := newTester(config, client)

	tester.startTime = time.Now()
	for _, nextGraphPoint := range tester.config.Graph {
		tester.interpolate2RPS(nextGraphPoint.TargetRPS, nextGraphPoint.Seconds2TargetRPS)
	}

	fmt.Println("Signal to kill all running request user")

	for len(tester.runningUsers) > 0 {
		tester.killOldestUser()
	}

	fmt.Println("Wait for death of all running request user")
	tester.waitGroup.Wait()

	select {
	case err := <-tester.errors:
		log.Fatal(err)
	default:
		fmt.Println("Load test done")
	}
}

func (tester *Tester) interpolate2RPS(targetRPS int, seconds2Interpolate int) {
	rps := float64(tester.currentRPS)

	rpsChange := float64(targetRPS) - rps
	rpsIncrement := rpsChange / float64(seconds2Interpolate)

	for elapsedSeconds := 0; elapsedSeconds < seconds2Interpolate; elapsedSeconds++ {
		rps += rpsIncrement
		tester.currentRPS = int(rps)

		tester.logRPS(targetRPS)

		for len(tester.runningUsers) < tester.currentRPS {
			tester.spawnNewUser()
		}
		for len(tester.runningUsers) > tester.currentRPS {
			tester.killOldestUser()
		}

		time.Sleep(1 * time.Second)
	}

	tester.currentRPS = targetRPS
}

func (tester *Tester) spawnNewUser() {
	user := NewRequestUser(tester.config.Target, tester.config.Paths, tester.client)

	tester.waitGroup.Add(1)
	go user.StartSendingRequestsPeriodically(&tester.waitGroup, tester.errors, 1*time.Second)

	tester.runningUsers = append(tester.runningUsers, user)
}

func (tester *Tester) killOldestUser() {
	tester.runningUsers[0].stopChan <- struct{}{}
	tester.runningUsers = tester.runningUsers[1:]
}

func (tester *Tester) logRPS(targetRPS int) {
	if !tester.logHeaderPrinted {
		fmt.Println("ElapsedTime\tTargetRPS\tCurrentRPS")
		tester.logHeaderPrinted = true
	}
	elapsedTime := time.Since(tester.startTime)
	fmt.Printf("%.1f\t\t%d\t\t%d\n", elapsedTime.Seconds(), targetRPS, tester.currentRPS)
}
