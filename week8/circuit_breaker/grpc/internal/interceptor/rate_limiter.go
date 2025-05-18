package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"microservices_course/week8/circuit_breaker/grpc/internal/rate_limiter"
)

type RateLimiterInterceptor struct {
	rateLimiter *rate_limiter.TokenBucketLimiter
}

func NewRateLimiterInterceptor(limiter *rate_limiter.TokenBucketLimiter) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{rateLimiter: limiter}
}

func (r *RateLimiterInterceptor) LimiterInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !r.rateLimiter.Allow() {
		return nil, status.Error(codes.ResourceExhausted, "to many request")
	}
	return handler(ctx, req)
}
