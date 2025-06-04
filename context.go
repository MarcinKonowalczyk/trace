package trace

import (
	"context"
)

type tracerKey struct{}

func WithTracer(ctx context.Context, tracer Tracer) context.Context {
	return context.WithValue(ctx, tracerKey{}, tracer)
}

func GetTracer(ctx context.Context) Tracer {
	if ctx == nil {
		// no context, no tracer
		return nil
	}
	tracer, ok := ctx.Value(tracerKey{}).(Tracer)
	if !ok {
		return nil
	}
	return tracer
}

// WITHOUT CONTEXT:
//
//	if tracer != nil {
//	    defer tracer.Un(tracer.Trace(trace.Here()))
//	}
//
// WITH CONTEXT:
// NOTE the extra `()`
//
//	defer trace.TraceCtx(ctx)()
//
//go:noinline
func TraceCtx(ctx context.Context, where ...string) func() {
	tracer := GetTracer(ctx)

	if tracer == nil {
		return func() {}
	} else {
		// we have a tracer
		where_str := whereToString(where...)
		t := tracer.Trace(where_str)
		return func() {
			tracer.Un(t)
		}
	}
}

//go:noinline
func HereCtx(ctx context.Context) string {
	tracer := GetTracer(ctx)
	if tracer != nil {
		return here(callerName(1))
	}
	return ""
}

func MessageCtx(ctx context.Context, args ...any) {
	tracer := GetTracer(ctx)
	if tracer != nil {
		tracer.Message(args...)
	}
}
