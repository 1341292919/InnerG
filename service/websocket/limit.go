package websocket

import "time"

// 一个简易的token令牌桶
type websocketMessageLimiter struct {
	rate       float64
	burst      float64
	tokens     float64
	lastRefill time.Time
}

func newWebsocketMessageLimiter(rate, burst int) *websocketMessageLimiter {
	now := time.Now()
	return &websocketMessageLimiter{
		rate:       float64(rate),
		burst:      float64(burst),
		tokens:     float64(burst),
		lastRefill: now,
	}
}

func (l *websocketMessageLimiter) Allow() bool {
	if l == nil || l.rate <= 0 || l.burst <= 0 {
		return true
	}

	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	l.lastRefill = now
	l.tokens += elapsed * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}
