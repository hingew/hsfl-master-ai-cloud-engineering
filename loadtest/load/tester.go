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
	fmt.Println("Beginn der Ramp-Up-Phase")
	rampupTime := time.Duration(tester.config.Rampup) * time.Second
	totalTime := time.Duration(tester.config.Duration) * time.Second
	cooldownTime := time.Duration(tester.config.Cooldown) * time.Second
	tester.done = make(chan struct{})

	var wg sync.WaitGroup
	errors := make(chan error, 1)

	// Ramp-Up-Phase
	for i := 0; i < tester.config.NumberUsers; i++ {
		wg.Add(1)
		go tester.startUser(&wg, errors)

		// Ramp-Up-Verzögerung
		time.Sleep(rampupTime / time.Duration(tester.config.NumberUsers))
	}
	fmt.Println("Ende der Ramp-Up-Phase, Beginn der Phase konstanter Last")

	stableLoadTime := totalTime - rampupTime - cooldownTime
	time.Sleep(stableLoadTime)

	close(tester.done)

	fmt.Println("Beginn der Absenkungsphase auf Null")
	time.Sleep(cooldownTime)

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
			target := tester.config.Targets[rand.Intn(len(tester.config.Targets))]
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
