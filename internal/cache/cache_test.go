package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetGet(t *testing.T) {
	c := NewMemoryCache()
	c.Set("India", 1)
	v, ok := c.Get("India")
	require.True(t, ok)
	require.Equal(t, 1, v)
}
func TestConcurrentAccess(t *testing.T) {
	c := NewMemoryCache()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			c.Set("k", i)
		}(i)
		go func() {
			defer wg.Done()
			_, _ = c.Get("k")
		}()
	}
	wg.Wait()
}
