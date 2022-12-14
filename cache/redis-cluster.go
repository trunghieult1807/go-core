package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"go-core/opentracing/jaeger"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go/ext"
)

func initRedisCluster(addrs []string) (*redis.ClusterClient, error) {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	_, err := clusterClient.Ping().Result()
	return clusterClient, err
}

type clusterRedisHelper struct {
	clusterClient *redis.ClusterClient
}

func (h *clusterRedisHelper) GetTransaction(ctx context.Context, transactionID string) CacheTransactionExecution {
	txPipeline := h.clusterClient.TxPipeline()
	return &redisCacheTransaction{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     txPipeline,
			transactionID: transactionID,
		},
	}
}
func (h *clusterRedisHelper) GetPipeline(ctx context.Context, transactionID string) CachePipelineExecution {
	pipeline := h.clusterClient.Pipeline()
	return &redisCachePipeline{
		baseRedisCachePipeline: baseRedisCachePipeline{
			Pipeliner:     pipeline,
			transactionID: transactionID,
		},
	}
}

func (h *clusterRedisHelper) Exists(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Exists", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	indicator, err := h.clusterClient.Exists(key).Result()
	if err != nil {
		return err
	}
	if indicator == 0 {
		return redis.Nil
	}
	return nil
}

func (h *clusterRedisHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Get", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := h.clusterClient.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Set", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = h.clusterClient.Set(key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (isSuccess bool, err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/SetNX", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	_, err = h.clusterClient.SetNX(key, string(data), expiration).Result()
	if err != nil {
		return false, err
	}
	return isSuccess, nil
}

func (h *clusterRedisHelper) Del(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Del", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.clusterClient.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Expire", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.clusterClient.Expire(key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) GetInterface(ctx context.Context, key string, value interface{}) (interface{}, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/GetInterface", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := h.clusterClient.Get(key).Result()
	if err != nil {
		return nil, err
	}

	typeValue := reflect.TypeOf(value)
	kind := typeValue.Kind()

	var outData interface{}
	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outData = reflect.New(typeValue).Interface()
	default:
		outData = reflect.Zero(typeValue).Interface()
	}
	err = json.Unmarshal([]byte(data), &outData)
	if err != nil {
		return nil, err
	}

	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		return reflect.ValueOf(outData).Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *clusterRedisHelper) DelMulti(ctx context.Context, keys ...string) error {
	var err error
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/DelMulti", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return err
}

func (h *clusterRedisHelper) GetKeysByPattern(ctx context.Context, pattern string, cursor uint64, limit int64) ([]string, uint64, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/GetKeysByPattern", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	var keys []string
	return keys, 0, err
}

func (h *clusterRedisHelper) SubscribeMessage(ctx context.Context, keySpace string, subscribeFunc SubscribeFunc) {
	subscribes := h.clusterClient.Subscribe(keySpace)
	messageChan := subscribes.Channel()
	for {
		select {
		case message, ok := <-messageChan:
			if ok {
				go subscribeFunc(CacheMessage{Message: *message})
			}
		}
	}
}

func (h *clusterRedisHelper) PublishMessage(ctx context.Context, keySpace string, message interface{}) error {
	result := h.clusterClient.Publish(keySpace, message)
	var out int64
	var err error
	if out, err = result.Result(); err != nil {
		return err
	}
	if out == 0 {
		return fmt.Errorf("published message with response:  %v", out)
	}
	return nil
}
func (h *clusterRedisHelper) GetMulti(ctx context.Context, data interface{}, keys ...string) ([]interface{}, error) {
	return nil, nil
}

func (h *clusterRedisHelper) RenameKey(ctx context.Context, oldkey, newkey string) error {
	return nil
}

func (h *clusterRedisHelper) GetStrLenght(ctx context.Context, key string) (int64, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/GetStrLenght", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return h.clusterClient.StrLen(key).Result()
}

func (h *clusterRedisHelper) GetType(ctx context.Context, key string) (string, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/GetType", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return h.clusterClient.Type(key).Result()
}

func (h *clusterRedisHelper) DebugObjectByKey(ctx context.Context, key string) (string, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/DebugObjectByKey", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return h.clusterClient.DebugObject(key).Result()
}

func (h *clusterRedisHelper) TimeExpire(ctx context.Context, key string) (time.Duration, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/TimeExpire", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return h.clusterClient.TTL(key).Result()
}

func (h *clusterRedisHelper) HSet(ctx context.Context, key, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}
func (h *clusterRedisHelper) HSetNX(ctx context.Context, key string, mapKey string, mapValue interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}
func (h *clusterRedisHelper) HGet(ctx context.Context, key, mapKey string) (value string, err error) {
	return value, err
}
func (h *clusterRedisHelper) HGetAll(ctx context.Context, key string, mapKeys []string) (values map[string]string, err error) {
	return values, err
}
func (h *clusterRedisHelper) HIncreaseBy(ctx context.Context, key, mapKey string, increase int64) (isIncreased bool, value string, err error) {
	return isIncreased, value, err
}

func (h *clusterRedisHelper) HMSet(ctx context.Context, key string, mapData map[string]interface{}, expiration time.Duration) (isSet bool, err error) {
	return isSet, err
}

func (h *clusterRedisHelper) HMGet(ctx context.Context, key string, fields []string) (result map[string]interface{}, err error) {
	return result, err
}
