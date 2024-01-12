package curve_tester

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/loadtest/net"
)

type RequestUser struct {
	target   string
	routes   []string
	stopChan chan struct{}
	client   net.Client
}

func NewRequestUser(target string, routes []string, client net.Client) *RequestUser {
	return &RequestUser{target: target, routes: routes, client: client, stopChan: make(chan struct{})}
}

func (user *RequestUser) StartSendingRequestsPeriodically(wg *sync.WaitGroup, errors chan<- error, sendInterval time.Duration) {
	samplingRate := 2
	samplingInterval := sendInterval / 2
	sampleCounter := 0

	defer wg.Done()

	for {
		select {
		case <-user.stopChan:
			fmt.Println("Request user stopped")
			return
		default:
			sampleCounter++
			if sampleCounter%samplingRate == 0 {
				route := user.routes[rand.Intn(len(user.routes))]
				if err := user.client.Send(user.target, route); err != nil {
					errors <- err
					return
				}
				sampleCounter = 0
				time.Sleep(samplingInterval)
			} else {
				time.Sleep(samplingInterval)
			}
		}
	}
}
