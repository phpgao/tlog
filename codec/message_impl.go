package codec

import "context"

// msg is the context of rpc.
type msg struct {
	context context.Context
	logger  interface{}
}

// resetDefault reset all fields of msg to default value.
func (m *msg) resetDefault() {
	m.logger = nil
}

// WithLogger sets logger into context message. Generally, the logger is
// created from WithFields() method.
func (m *msg) WithLogger(l interface{}) {
	m.logger = l
}

// Logger returns logger from context message.
func (m *msg) Logger() interface{} {
	return m.logger
}

// Message returns the message of context.
func Message(ctx context.Context) Msg {
	if m, ok := ctx.Value(ContextKeyMessage).(*msg); ok {
		return m
	}
	return &msg{context: ctx}
}

// EnsureMessage returns context and message, if there is a message in context,
// returns the original one, if not, returns a new one.
func EnsureMessage(ctx context.Context) (context.Context, Msg) {
	if m, ok := ctx.Value(ContextKeyMessage).(*msg); ok {
		return ctx, m
	}
	return WithNewMessage(ctx)
}

// WithNewMessage creates a new empty message, retrieves it from the message pool,
// and associates it with the provided context.
//
// Important: The returned message is obtained from a pool to optimize memory usage.
// Users are responsible for manually invoking codec.PutBackMessage(msg) after use.
// Failure to return the message to the pool doesn't result in a traditional memory leak,
// where memory is never reclaimed. Instead, it may lead to a gradual increase in memory
// footprint over time, as messages are not being recycled as efficiently. This can
// eventually lead to higher than normal memory consumption, although the memory
// may still be eventually released.
func WithNewMessage(ctx context.Context) (context.Context, Msg) {
	m := msgPool.Get().(*msg)
	ctx = context.WithValue(ctx, ContextKeyMessage, m)
	m.context = ctx
	return ctx, m
}
