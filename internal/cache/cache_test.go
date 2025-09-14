package cache

import (
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
