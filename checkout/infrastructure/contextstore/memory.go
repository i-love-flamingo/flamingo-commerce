package contextstore

import (
	"context"
	"sync"

	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Memory saves all contexts in a simple map
	Memory struct {
		mx      sync.RWMutex
		storage map[string]process.Context
	}
)

var _ process.ContextStore = new(Memory)

// Inject dependencies
func (m *Memory) Inject() *Memory {
	m.storage = make(map[string]process.Context)

	return m
}

// Store a given context
func (m *Memory) Store(_ context.Context, key string, value process.Context) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.storage[key] = value

	return nil
}

// Get a stored context
func (m *Memory) Get(_ context.Context, key string) (process.Context, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	value, ok := m.storage[key]

	return value, ok
}

// Delete a stored context, nop if it doesn't exist
func (m *Memory) Delete(_ context.Context, key string) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.storage, key)

	return nil
}
