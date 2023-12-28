package load

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

type Tester struct {
	config TesterConfig
	client net.Client
	done   chan struct{}
}

func NewTester(config TesterConfig, client net.Client) *Tester {
	return &Tester{config: config, client: client}
}

func (tester *Tester) Run() error {
	var wg sync.WaitGroup
	errors := make(chan error, 1)
	tester.done = make(chan struct{})

	activeRoutines := 0
	targetRoutines := 0
	currentRPS := 0
	for _, spec := range tester.config.RampSpecifications {
		rpsChange := spec.TargetRPS - currentRPS
		steps := spec.Duration
		rpsIncrement := rpsChange / steps

		for i := 0; i < steps; i++ {
			currentRPS += rpsIncrement
			targetRoutines = currentRPS * spec.Duration

			for activeRoutines < targetRoutines {
				wg.Add(1)
				go tester.startUser(&wg, errors)
				activeRoutines++
			}
			for activeRoutines > targetRoutines {
				activeRoutines--
			}

			time.Sleep(1 * time.Second)
		}
		currentRPS = spec.TargetRPS
	}

	close(tester.done)
	wg.Wait()

	select {
	case err := <-errors:
		return err
	default:
		fmt.Println("Lasttest abgeschlossen")
		return nil
	}
}

func (tester *Tester) startUser(wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

	for {
		select {
		case <-tester.done:
			return
		default:
			target := tester.config.Target
			path := tester.config.Paths[rand.Intn(len(tester.config.Paths))]
			fmt.Printf("Sende Anfrage an %s%s\n", target, path)
			if err := tester.client.Send(target, path); err != nil {
				errors <- err
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
