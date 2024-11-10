package codec

import (
	"sync"
)

// ContextKey is trpc context key type, the specific value is judged
// by interface, the interface will both judge value and type. Defining
// a new type can avoid string value conflict.
type ContextKey string

// MetaData is request penetrate message.
type MetaData map[string][]byte

var msgPool = sync.Pool{
	New: func() interface{} {
		return &msg{}
	},
}

// Clone returns a copied meta data.
func (m MetaData) Clone() MetaData {
	if m == nil {
		return nil
	}
	md := MetaData{}
	for k, v := range m {
		md[k] = v
	}
	return md
}

// CommonMeta is common meta message.
type CommonMeta map[interface{}]interface{}

// Clone returns a copied common meta message.
func (c CommonMeta) Clone() CommonMeta {
	if c == nil {
		return nil
	}
	cm := CommonMeta{}
	for k, v := range c {
		cm[k] = v
	}
	return cm
}

// trpc context key data
const (
	ContextKeyMessage = ContextKey("TRPC_MESSAGE")
	// ServiceSectionLength is the length of service section,
	// service name example: trpc.app.server.service
	ServiceSectionLength = 4
)

// Msg defines core message data for Multi-protocol, business protocol
// should set this message when packing and unpacking data.
type Msg interface {
	// WithLogger sets logger into context.
	WithLogger(interface{})
	// Logger returns logger from context.
	Logger() interface{}
}
