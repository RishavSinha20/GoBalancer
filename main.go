package main

import (
	"log"
	"net/http"

	"loadbalancer/balancer"
	"loadbalancer/config"
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

	lb := balancer.NewLoadBalancer(urls)
	proxy := proxy.NewProxy(lb)

	log.Println("Load balancer running on", cfg.Port)
	http.ListenAndServe(cfg.Port, proxy)
}