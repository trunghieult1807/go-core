package interceptor

import (
	"context"
	"go-core/util"
	"reflect"

	"google.golang.org/grpc"
)

// SetIDUnaryServerInterceptor represents set id
func SetIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		value := reflect.ValueOf(req)
		typeOf := reflect.TypeOf(req)
		if typeOf.Kind() != reflect.Ptr {
			return handler(ctx, req)
		}

		trans := value.Elem().FieldByName("RequestId")
		if trans.Kind() != reflect.String {
			return handler(ctx, req)
		}
		if trans.String() != "" {
			return handler(ctx, req)
		}
		id := util.GetID()
		if trans.CanSet() {
			trans.SetString(id)
		}
		return handler(ctx, req)
	}
}
