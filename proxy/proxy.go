package proxy

import (
	"loadbalancer/balancer"
	"loadbalancer/metrics"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

type Proxy struct {
	lb *balancer.LoadBalancer
}

func NewProxy(lb *balancer.LoadBalancer) *Proxy {
	return &Proxy{lb: lb}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests()

	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		backend := p.lb.Next(r)

		if backend == nil {
			break
		}

		backend.IncrementConnections()
		proxy := httputil.NewSingleHostReverseProxy(backend.URL)
		errorOccurred := false

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Println("Error :", err)
			errorOccurred = true
		}

		proxy.ServeHTTP(w, r)

		backend.DecrementConnections()

		if !errorOccurred {
			return
		}

		log.Println("Retrying request....")
		time.Sleep(100 * time.Millisecond)
	}

	metrics.IncFailures()
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}
