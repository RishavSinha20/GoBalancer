package proxy

import (
	"bytes"
	"io"
	"loadbalancer/balancer"
	"loadbalancer/cache"
	"loadbalancer/metrics"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type Proxy struct {
	lb    *balancer.LoadBalancer
	cache *cache.Cache
}

func NewProxy(lb *balancer.LoadBalancer) *Proxy {
	return &Proxy{
		lb:    lb,
		cache: cache.NewCache(10 * time.Second),
	}
}

func getClientIP(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics.IncRequests()

	cacheKey := r.URL.String()

	// Only cache GET requests
	if r.Method == "GET" {
		if data, found := p.cache.Get(cacheKey); found {
			w.Write(data)
			return
		}
	}

	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		backend := p.lb.Next(r)

		if backend == nil {
			break
		}

		backend.IncrementConnections()

		proxy := httputil.NewSingleHostReverseProxy(backend.URL)

		var responseBody []byte

		proxy.ModifyResponse = func(resp *http.Response) error {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			responseBody = body
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			return nil
		}

		errorOccurred := false

		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Println("Error:", err)
			errorOccurred = true
		}

		proxy.ServeHTTP(w, r)

		backend.DecrementConnections()

		if !errorOccurred {
			// Cache response
			if r.Method == "GET" {
				p.cache.Set(cacheKey, responseBody)
			}
			return
		}

		time.Sleep(100 * time.Millisecond)
	}

	metrics.IncFailures()
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}
