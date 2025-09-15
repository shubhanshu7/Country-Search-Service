package cache

import "sync"

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
}

type MemoryCache struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{data: make(map[string]any)}
}
func (m *MemoryCache) Get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}
func (m *MemoryCache) Set(key string, value any) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}
