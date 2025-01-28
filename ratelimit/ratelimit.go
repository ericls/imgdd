package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ClientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	clients map[string]*ClientLimiter
	mu      sync.Mutex
	limit   rate.Limit
	burst   int
}

func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*ClientLimiter),
		mu:      sync.Mutex{},
		limit:   limit,
		burst:   burst,
	}
	return rl
}

func (rl *RateLimiter) Cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		for ip, client := range rl.clients {
			if time.Since(client.lastSeen) > 10*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) IsRateLimited(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	client, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &ClientLimiter{
			limiter:  rate.NewLimiter(rl.limit, rl.burst),
			lastSeen: time.Now(),
		}
		return false
	}
	client.lastSeen = time.Now()
	return !client.limiter.Allow()
}
