package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/cespare/xxhash"
)

type BalancerAlgorithm int

const (
	RoundRobin BalancerAlgorithm = iota
	WeightedRoundRobin
	IPHash
)

type Target struct {
	Proxy   http.Handler
	Weight  int
	Current int
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
	case IPHash:
		balancer.balanceByIPHash(w, r)
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

	target := balancer.findTargetByWeight()

	if target != nil {
		target.Current++
		defer func() { target.Current-- }()
		target.Proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}
}

func (balancer *LoadBalancer) findTargetByWeight() *Target {
	var min float64 = -1
	var target *Target
	for _, t := range balancer.targets {
		ratio := float64(t.Current) / float64(t.Weight)
		if ratio < min || min == -1 {
			min = ratio
			target = t
		}
	}
	return target
}

func (balancer *LoadBalancer) balanceByIPHash(w http.ResponseWriter, r *http.Request) {
	ip := getClientIp(r)
	next := xxhash.Sum64([]byte(ip)) % uint64(len(balancer.targets))
	target := balancer.targets[next]
	target.Proxy.ServeHTTP(w, r)
}

func getClientIp(r *http.Request) string {
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip != "" {
		return ip
	}
	return r.RemoteAddr
}
