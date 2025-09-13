package cache

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
}

type MemoryCache struct {
	data map[string]any
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{data: make(map[string]any)}
}
func (m *MemoryCache) Get(key string) (any, bool) {
	v, ok := m.data[key]
	return v, ok
}
func (m *MemoryCache) Set(key string, value any) {
	m.data[key] = value
}
