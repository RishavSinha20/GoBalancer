# GoLoadBalancer

A production-style Layer 7 Load Balancer written in Go featuring multiple routing strategies, active health checks, automatic failover, retry mechanisms, rate limiting, response caching, and observability.

The project was built to understand how modern traffic management systems such as NGINX, HAProxy, and Envoy handle request routing, fault tolerance, and infrastructure monitoring.

---

# Architecture

```text
                 ┌─────────────┐
                 │   Client    │
                 └──────┬──────┘
                        │
                        ▼
             ┌────────────────────┐
             │   Load Balancer    │
             │                    │
             │ Round Robin        │
             │ Least Connections  │
             │ IP Hashing         │
             └─────────┬──────────┘
                       │
       ┌───────────────┼───────────────┐
       │               │               │
       ▼               ▼               ▼
┌────────────┐ ┌────────────┐ ┌────────────┐
│ Backend 1  │ │ Backend 2  │ │ Backend N  │
└────────────┘ └────────────┘ └────────────┘

       ▲               ▲               ▲
       └────── Health Checks ──────────┘
```

---

# Features

## Core Routing

* Reverse Proxy implementation
* Round Robin load balancing
* Least Connections load balancing
* IP Hash load balancing

## Fault Tolerance

* Active health checks
* Automatic backend failover
* Backend recovery detection
* Graceful degradation

## Reliability

* Retry mechanism for failed requests
* Request timeout handling
* Health monitoring subsystem

## Performance

* In-memory response caching
* Configurable cache TTL

## Security

* Request rate limiting
* Optional TLS support

## Observability

* Request metrics collection
* Failure metrics collection
* Health monitoring logs

## Infrastructure

* YAML-based configuration
* Graceful shutdown support

---

# Project Structure

```text
loadbalancer/
│
├── balancer/
│   ├── backend.go
│   └── balancer.go
│
├── proxy/
│   └── proxy.go
│
├── health/
│   └── health.go
│
├── middleware/
│   └── ratelimiter.go
│
├── metrics/
│   └── metrics.go
│
├── cache/
│   └── cache.go
│
├── config/
│   ├── config.go
│   └── watcher.go
│
├── main.go
├── go.mod
├── go.sum
└── config.yaml
```

---

# Routing Strategies

## Round Robin

Requests are distributed sequentially across available backend servers.

```text
Request 1 → Backend 1
Request 2 → Backend 2
Request 3 → Backend 1
Request 4 → Backend 2
```

---

## Least Connections

Traffic is routed to the backend currently handling the fewest active connections.

Useful when requests have varying execution times.

---

## IP Hash

A client's IP address is hashed to determine backend selection.

This ensures the same client is consistently routed to the same backend.

---

# Health Monitoring

Every backend is periodically checked using an HTTP health probe.

If a backend:

* Times out
* Becomes unreachable
* Returns a non-200 response

it is marked as unavailable and removed from routing.

Recovered backends are automatically reintroduced into the load balancing pool.

---

# Configuration

Example `config.yaml`:

```yaml
port: ":8080"

strategy: "round_robin"

backends:
  - url: "http://localhost:8081"
  - url: "http://localhost:8082"
```

Supported strategies:

```text
round_robin
least_conn
ip_hash
```

---

# Installation

Clone the repository:

```bash
git clone https://github.com/RishavSinha20/goloadbalancer.git

cd goloadbalancer
```

Install dependencies:

```bash
go mod tidy
```

---

# Running the Project

## Start Backend Server 1

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Response from Backend 1")
    })

    http.ListenAndServe(":8081", nil)
}
```

Run:

```bash
go run backend1.go
```

---

## Start Backend Server 2

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Response from Backend 2")
    })

    http.ListenAndServe(":8082", nil)
}
```

Run:

```bash
go run backend2.go
```

---

## Start Load Balancer

```bash
go run main.go
```

Server will start on:

```text
http://localhost:8080
```

---

# Testing

## Basic Load Balancing

```bash
curl http://localhost:8080
curl http://localhost:8080
curl http://localhost:8080
```

Expected:

```text
Response from Backend 1
Response from Backend 2
Response from Backend 1
```

---

## Health Check Testing

1. Stop Backend 1

```bash
Ctrl + C
```

2. Send requests

```bash
curl http://localhost:8080
```

Expected:

```text
Response from Backend 2
```

Backend 1 should automatically be removed from routing.

---

## Backend Recovery

Restart Backend 1.

Within the next health-check interval, it should automatically rejoin the load balancing pool.

---

## Rate Limiting

Generate burst traffic:

```bash
for i in {1..20}; do
  curl http://localhost:8080 &
done
```

Expected:

```text
429 Too Many Requests
```

for some requests.

---

## Failure Handling

Stop all backend servers.

Then:

```bash
curl http://localhost:8080
```

Expected:

```text
503 Service Unavailable
```

---

# Metrics

The load balancer tracks:

* Total Requests
* Failed Requests

Metrics endpoint:

```text
/metrics
```

Example response:

```json
{
  "total_requests": 120,
  "failed_requests": 4
}
```

---

# Tradeoffs

## Caching

Pros:

* Reduced backend load
* Faster response times

Cons:

* Potential stale data
* Memory overhead

---

## Retry Logic

Pros:

* Better resiliency
* Improved user experience

Cons:

* Potential duplicate requests
* Increased backend traffic

---

## Least Connections

Pros:

* Better load distribution

Cons:

* Requires active connection tracking

---


