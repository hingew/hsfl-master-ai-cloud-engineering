package curve_tester

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

type RequestHandler struct {
	target string
	routes []string
	client net.Client
}

func NewRequestHandler(target string, routes []string, client net.Client) *RequestHandler {
	return &RequestHandler{target: target, routes: routes, client: client}
}

func (handler *RequestHandler) Run(stopChan chan struct{}, rpsChan chan int, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	doRequestWg := sync.WaitGroup{}

	currentRPS := 0

	for {
		select {
		case <-stopChan:
			fmt.Println("Wait for done Request to return")
			doRequestWg.Wait()
			return
		case currentRPS = <-rpsChan:
		default:
			startTime := time.Now()

			for i := 0; i < currentRPS; i++ {
				go handler.doRequest(&doRequestWg)
			}

			timeSinceStart := time.Since(startTime)
			if timeSinceStart > time.Second {
				fmt.Printf("Unable to do %d requests within 1 second. Took %.3f seconds\n", currentRPS, timeSinceStart.Seconds())
			} else {
				durationUntilNextSecond := time.Second - timeSinceStart
				time.Sleep(durationUntilNextSecond)
			}
		}
	}
}

func (handler *RequestHandler) doRequest(requestWg *sync.WaitGroup) {
	requestWg.Add(1)
	defer requestWg.Done()

	route := handler.routes[rand.Intn(len(handler.routes))]
	err := handler.client.Send(handler.target, route)
	if err != nil {
		log.Fatal(err)
		return
	}
}
