package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type BalancerAlgorithm int

const (
	RoundRobin BalancerAlgorithm = iota
	WeightedRoundRobin
	HaukesStuff // TODO ich erinnere mich leider nicht mehr welchen Algorithmus Hauke implementieren wollte
)

type LoadBalancer struct {
	targets   []http.Handler
	index     int
	algorithm BalancerAlgorithm
}

func NewLoadBalancer(targetUrls []*url.URL, algorithm BalancerAlgorithm) *LoadBalancer {
	targets := make([]http.Handler, len(targetUrls))
	for i, targetUrl := range targetUrls {
		targets[i] = httputil.NewSingleHostReverseProxy(targetUrl)
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
	target.ServeHTTP(w, r)
	balancer.index = (balancer.index + 1) % len(balancer.targets)
}

func (balancer *LoadBalancer) balanceByWeightedRoundRobin(w http.ResponseWriter, r *http.Request) {
	// TODO Robert
}

func (balancer *LoadBalancer) balanceBy(w http.ResponseWriter, r *http.Request) {
	// TODO Hauke
}
