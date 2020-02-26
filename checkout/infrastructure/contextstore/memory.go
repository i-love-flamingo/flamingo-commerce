package contextstore

import (
	"flamingo.me/flamingo-commerce/v3/checkout/domain/placeorder/process"
)

type (
	// Memory saves all contexts in a simple map
	Memory struct {
		storage map[string]process.Context
	}
)

// Inject dependencies
func (m *Memory) Inject() *Memory {
	m.storage = make(map[string]process.Context)

	return m
}

// Store a given context
func (m *Memory) Store(key string, value process.Context) error {
	m.storage[key] = value

	return nil
}

// Get a stored context
func (m *Memory) Get(key string) (process.Context, bool) {
	value, ok := m.storage[key]

	return value, ok
}

// Delete a stored context, nop if it doesn't exist
func (m *Memory) Delete(key string) error {
	delete(m.storage, key)

	return nil
}
