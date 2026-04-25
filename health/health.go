package health

import (
	"log"
	"net/http"
	"time"

	"loadbalancer/balancer"
)

func StartHealthCheck(backends []*balancer.Backend) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, b := range backends {
				go checkBackend(b)
			}
		}
	}
}

func checkBackend(b *balancer.Backend) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(b.URL.String())
	if err != nil || resp.StatusCode != 200 {
		b.SetAlive(false)
		log.Println("Backend DOWN:", b.URL)
		return
	}

	b.SetAlive(true)
	log.Println("Backend UP:", b.URL)
}
