package balancer

import "sync/atomic"

type LoadBalancer struct {
	backends []*Backend
	counter  uint64
}

func NewLoadBalancer(urls []string) *LoadBalancer {
	var backends []*Backend

	for _, u := range urls {
		backends = append(backends, NewBackend(u))
	}

	return &LoadBalancer{
		backends: backends,
	}
}
func (lb *LoadBalancer) Next() *Backend {
	n := len(lb.backends)

	for i := 0; i < n; i++ {
		idx := atomic.AddUint64(&lb.counter, 1)
		b := lb.backends[int(idx)%n]

		if b.IsAlive() {
			return b
		}
	}

	return nil // all backends dead
}
func (lb *LoadBalancer) GetBackends() []*Backend {
	return lb.backends
}
