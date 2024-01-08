package loadbalancer

import (
	"net/http"
	"testing"

	"gotest.tools/assert"
)

type TestProxy struct {
}

func (p *TestProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func TestLoadBalancer(t *testing.T) {
	t.Run("balanceByRoundRobin", func(t *testing.T) {
		lb := LoadBalancer{
			targets: []*Target{
				{
					Proxy: &TestProxy{},
				},
				{
					Proxy: &TestProxy{},
				},
				{
					Proxy: &TestProxy{},
				},
			},
			index: 0,
		}

		for i := 0; i < 6; i++ {
			assert.Equal(t, lb.index, i%3)

			lb.balanceByRoundRobin(nil, nil)
		}
		assert.Equal(t, lb.index, 0)
	})

	t.Run("balanceByWeightedRoundRobin", func(t *testing.T) {
		lb := LoadBalancer{
			targets: []*Target{
				{
					Proxy:   &TestProxy{},
					Current: 5,
					Weight:  1,
				},
				{
					Proxy:   &TestProxy{},
					Current: 4,
					Weight:  2,
				},
				{
					Proxy:   &TestProxy{},
					Current: 7,
					Weight:  3,
				},
			},
			index: 0,
		}

		t.Run("findTargetByWeight", func(t *testing.T) {
			target := lb.findTargetByWeight()

			assert.Equal(t, target, lb.targets[1])
		})
	})
}
