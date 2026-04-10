package balancer

import "sync/atomic"

type LoadBalancer struct{
	backends []*Backend
	counter uint64
}

func NewLoadBalancer(urls []string) *LoadBalancer{
	var backends []*Backend

	for _, u := range urls{
		backends = append(backends, NewBackend(u))
	}

	return &LoadBalancer{
		backends : backends,
	}
}
func (lb *LoadBalancer) Next() *Backend {
	idx := atomic.AddUint64(&lb.counter, 1)
	return lb.backends[int(idx)%len(lb.backends)]
}	