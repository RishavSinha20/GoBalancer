package proxy


import(
	"log"
	"net/http"
	"net/http/httputil"
	"loadbalancer/balancer"
)

type Proxy struct{
	lb *balancer.LoadBalancer
}

func NewProxy(lb *balancer.LoadBalancer) *Proxy {
	return &Proxy{lb: lb}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := p.lb.Next()

	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("Proxy error:", err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}

	log.Println("Forwarding request to:", backend.URL)

	proxy.ServeHTTP(w, r)
}