package interceptor

import (
	"context"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerInterceptor struct {
	cd *gobreaker.CircuitBreaker
}

func NewCircuitBreakerInter(cd *gobreaker.CircuitBreaker) *CircuitBreakerInterceptor {
	return &CircuitBreakerInterceptor{cd: cd}
}

func (c *CircuitBreakerInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	res, err := c.cd.Execute(func() (interface{}, error) {
		return handler(ctx, req)
	})
	if err != nil {
		if err == gobreaker.ErrOpenState {
			return nil, status.Error(codes.Unavailable, "service unavailable")
		}
		return nil, err
	}
	return res, nil
}
