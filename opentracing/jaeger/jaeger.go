package jaeger

import (
	"context"
	"io"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type clientSpanTagKey struct{}

// Start returns a new span or gets a span from parent
func Start(ctx context.Context, methodName string, spanKind opentracing.Tag, tags ...opentracing.Tag) opentracing.Span {
	var parentSpanCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentSpanCtx = parent.Context()
	}
	opts := []opentracing.StartSpanOption{opentracing.ChildOf(parentSpanCtx)}
	opts = append(opts, spanKind)
	for _, tag := range tags {
		opts = append(opts, tag)
	}
	if tagx := ctx.Value(clientSpanTagKey{}); tagx != nil {
		if opt, ok := tagx.(opentracing.StartSpanOption); ok {
			opts = append(opts, opt)
		}
	}
	span := opentracing.GlobalTracer().StartSpan(methodName, opts...)
	return span
}

// Continue returns a reference span
func Continue(spanCtx opentracing.SpanContext, methodName string, spanKind opentracing.Tag, tags ...opentracing.Tag) opentracing.Span {
	opts := []opentracing.StartSpanOption{opentracing.FollowsFrom(spanCtx)}
	opts = append(opts, spanKind)
	for _, tag := range tags {
		opts = append(opts, tag)
	}
	span := opentracing.GlobalTracer().StartSpan(methodName, opts...)
	return span
}

// Finish finalizes span
func Finish(span opentracing.Span, err error) {
	if err != nil && err != io.EOF && err != redis.Nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("event", "error"), log.String("message", err.Error()))
	}
	span.Finish()
}
