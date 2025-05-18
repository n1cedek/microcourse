package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"microservices_course/week8/rate_limiter/grpc/internal/metric"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	res, err := handler(ctx, req)
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		//metric.HistogramResponseTimeObserve("error", diffTime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		//metric.HistogramResponseTimeObserve("success", diffTime.Seconds())
	}
	return res, nil
}
