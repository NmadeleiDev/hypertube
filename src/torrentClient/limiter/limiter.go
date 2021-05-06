package limiter

import "sync"

type RateLimiter struct {
	mu    sync.Mutex
	limit int
	value int
	add   chan struct{}
}

func (r *RateLimiter) Init(limit int) {
	r.limit = limit
	r.add = make(chan struct{}, limit)
}

func (r *RateLimiter) Add() {
	r.add <- struct{}{}

	r.mu.Lock()
	r.value ++
	r.mu.Unlock()
}

func (r *RateLimiter) Pop() {
	<- r.add

	r.mu.Lock()
	r.value --
	r.mu.Unlock()
}

func (r *RateLimiter) GetVal() int {
	r.mu.Lock()
	val := r.value
	r.mu.Unlock()

	return val
}

func (r *RateLimiter) Destroy() {
	close(r.add)
}
