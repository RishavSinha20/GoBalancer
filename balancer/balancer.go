package balancer

import (
	"hash/fnv"
	"net/http"
	"sync/atomic"
)

type LoadBalancer struct {
	backends []*Backend
	counter  uint64
	strategy string
}

func NewLoadBalancer(urls []string, strategy string) *LoadBalancer {
	var backends []*Backend

	for _, u := range urls {
		backends = append(backends, NewBackend(u))
	}

	return &LoadBalancer{
		backends: backends,
		strategy: strategy,
	}
}
func (lb *LoadBalancer) GetBackends() []*Backend {
	return lb.backends
}

// ---STRATEGIES-----

func (lb *LoadBalancer) getRoundRobin() *Backend {
	n := len(lb.backends)

	for i := 0; i < n; i++ {
		idx := atomic.AddUint64(&lb.counter, 1)
		b := lb.backends[int(idx)%n]

		if b.IsAlive() {
			return b
		}
	}

	return nil
}

func (lb *LoadBalancer) getLeastConnections() *Backend {
	var selected *Backend

	for _, b := range lb.backends {
		if !b.IsAlive() {
			continue
		}

		if selected == nil || b.GetConnections() < selected.GetConnections() {
			selected = b
		}
	}
	return selected
}

func (lb *LoadBalancer) getIPHash(ip string) *Backend {
	h := fnv.New32a()
	h.Write([]byte(ip))
	idx := h.Sum32() % uint32(len(lb.backends))

	return lb.backends[idx]
}

func (lb *LoadBalancer) Next(r *http.Request) *Backend {
	switch lb.strategy {

	case "least_conn":
		return lb.getLeastConnections()

	case "ip_hash":
		return lb.getIPHash(r.RemoteAddr)

	default:
		return lb.getRoundRobin()
	}

}
