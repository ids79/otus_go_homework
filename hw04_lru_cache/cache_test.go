package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache[int](10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache[int](5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Equal(t, 0, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache[int](3)

		_ = c.Set("aaa", 100)
		_ = c.Set("bbb", 200)
		_ = c.Set("ccc", 300)
		_ = c.Set("ddd", 400)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, 0, val)
		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		_, _ = c.Get("bbb")
		_, _ = c.Get("ccc")
		_ = c.Set("eee", 500)
		val, ok = c.Get("ddd")
		require.False(t, ok)
		require.Equal(t, 0, val)

		c.Clear()
		require.False(t, c.Set("bbb", 10))
		require.False(t, c.Set("ccc", 20))
		require.False(t, c.Set("ddd", 30))
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache[int](10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
