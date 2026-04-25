package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"loadbalancer/balancer"
	"loadbalancer/config"
	"loadbalancer/health"
	"loadbalancer/middleware"
	"loadbalancer/proxy"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	for _, b := range cfg.Backends {
		urls = append(urls, b.URL)
	}

	lb := balancer.NewLoadBalancer(urls, cfg.Strategy)

	go health.StartHealthCheck(lb.GetBackends())

	proxy := proxy.NewProxy(lb)

	// Rate limiter
	rl := middleware.NewRateLimiter(10)

	handler := rl.Middleware(proxy)

	server := &http.Server{
		Addr:    cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Println("Server running on", cfg.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Println("Server stopped:", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)

	log.Println("Server exited")
}
