package interceptor

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTTPCodeSimulatorInterceptor represents helper for simulating http code
func HTTPCodeSimulatorInterceptor(isEnabled bool, apis []string, timesleep int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if isEnabled {
			var flag bool
			for _, methodName := range apis {
				if strings.HasSuffix(info.FullMethod, methodName) {
					flag = true
					break
				}
			}

			if flag {
				handler(ctx, req)
				time.Sleep(time.Duration(timesleep) * time.Second)
				return nil, status.Newf(codes.Unavailable, "Unavailable").Err()
			}
		}
		return handler(ctx, req)
	}
}
