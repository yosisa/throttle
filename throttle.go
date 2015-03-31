package throttle

import (
	"sync"
	"time"
)

type Bucket struct {
	tokens   float64
	rate     float64
	cap      float64
	lastFill time.Time
	m        sync.Mutex
}

func NewBucket(rate, cap int64) *Bucket {
	return &Bucket{
		tokens:   float64(cap),
		rate:     float64(rate),
		cap:      float64(cap),
		lastFill: time.Now(),
	}
}

func (b *Bucket) Take(n int64) (int64, time.Duration) {
	b.m.Lock()
	defer b.m.Unlock()
	now := time.Now()
	diff := now.Sub(b.lastFill).Seconds()
	if (b.cap-b.tokens)/b.rate <= diff {
		b.tokens = b.cap
	} else {
		b.tokens += b.rate * diff
	}
	b.lastFill = now

	intTokens := int64(b.tokens)
	if intTokens >= n {
		b.tokens -= float64(n)
		return n, 0
	}

	absence := n - intTokens
	b.tokens -= float64(intTokens)
	wait := time.Duration((float64(absence) - b.tokens) / b.rate * 1e9)
	return intTokens, wait
}

func (b *Bucket) TakeExactly(n int64) {
	for rem := n; rem > 0; {
		tokens, wait := b.Take(rem)
		rem -= tokens
		time.Sleep(wait)
	}
}
