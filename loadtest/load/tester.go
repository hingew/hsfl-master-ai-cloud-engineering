package load

import (
	"fmt"
	"loadtest/net"
	"math/rand"
	"sync"
	"time"
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
	fmt.Printf("Starte neue Benutzer-Goroutine\n")
	defer wg.Done()

	for {
		select {
		case <-tester.done:
			return
		default:
			target := tester.config.Targets[rand.Intn(len(tester.config.Targets))]
			fmt.Printf("Sende Anfrage an %s\n", target)
			if err := tester.client.Send(target, tester.config.Path); err != nil {
				fmt.Printf("Fehler beim Senden der Anfrage: %v\n", err)
				errors <- err
				return
			}

			// Pausieren zwischen den Anfragen
			time.Sleep(1 * time.Second)
		}
	}
}
