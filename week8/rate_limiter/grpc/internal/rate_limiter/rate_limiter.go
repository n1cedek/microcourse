package rate_limiter

import (
	"context"
	"time"
)

type TokenBucketLimiter struct {
	tokenBucketCh chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int, period time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{tokenBucketCh: make(chan struct{}, limit)}

	for i := 0; i < limit; i++ {
		limiter.tokenBucketCh <- struct{}{}
	}

	replInterval := period.Nanoseconds() / int64(limit)
	go limiter.replenishTokens(ctx, time.Duration(replInterval))

	return limiter
}

// replenishTokens — добавляет токены в канал с заданным интервалом
func (t *TokenBucketLimiter) replenishTokens(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.tokenBucketCh <- struct{}{}
		}
	}
}

// Allow — проверяет, можно ли выполнить действие (т.е. есть ли токен)
func (t *TokenBucketLimiter) Allow() bool {
	select {
	case <-t.tokenBucketCh:
		return true
	default:
		return false
	}
}
