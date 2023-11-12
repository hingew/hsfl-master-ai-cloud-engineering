package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type BalancerAlgorithm int

const (
	RoundRobin BalancerAlgorithm = iota
	WeightedRoundRobin
	HaukesStuff // TODO ich erinnere mich leider nicht mehr welchen Algorithmus Hauke implementieren wollte
)

type Target struct {
	Proxy   http.Handler
	Weight  int // Kapazität des Servers
	Current int // Aktuelle Anzahl der Verbindungen
}

type LoadBalancer struct {
	targets   []*Target
	index     int
	algorithm BalancerAlgorithm
	mutex     sync.Mutex
}

func NewLoadBalancer(targetUrls []*url.URL, weights []int, algorithm BalancerAlgorithm) *LoadBalancer {
	targets := make([]*Target, len(targetUrls))
	for i, targetUrl := range targetUrls {
		targets[i] = &Target{
			Proxy:  httputil.NewSingleHostReverseProxy(targetUrl),
			Weight: weights[i],
		}
	}

	return &LoadBalancer{
		index:     0,
		targets:   targets,
		algorithm: algorithm,
	}
}

func (balancer *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Use Container Target: %d", balancer.index)
	switch balancer.algorithm {
	case RoundRobin:
		balancer.balanceByRoundRobin(w, r)
	case WeightedRoundRobin:
		balancer.balanceByWeightedRoundRobin(w, r)
	case HaukesStuff:
		balancer.balanceBy(w, r)
	}
}

func (balancer *LoadBalancer) balanceByRoundRobin(w http.ResponseWriter, r *http.Request) {
	target := balancer.targets[balancer.index]
	target.Proxy.ServeHTTP(w, r)
	balancer.index = (balancer.index + 1) % len(balancer.targets)
}

func (balancer *LoadBalancer) balanceByWeightedRoundRobin(w http.ResponseWriter, r *http.Request) {
	balancer.mutex.Lock()
	defer balancer.mutex.Unlock()

	// Wählen Sie den Server mit dem niedrigsten Verhältnis von aktiven Verbindungen zu Gewicht
	var min float64 = -1
	var target *Target
	for _, t := range balancer.targets {
		ratio := float64(t.Current) / float64(t.Weight)
		if ratio < min || min == -1 {
			min = ratio
			target = t
		}
	}

	if target != nil {
		// Simulieren Sie den Beginn einer neuen Verbindung
		target.Current++
		defer func() { target.Current-- }()
		target.Proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}

func (balancer *LoadBalancer) balanceBy(w http.ResponseWriter, r *http.Request) {
	// TODO Hauke
}
