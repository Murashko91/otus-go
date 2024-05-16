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
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

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
		require.Nil(t, val)
	})

	t.Run("Test Capacity", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("key1", nil)
		require.False(t, wasInCache)
		val, ok := c.Get("key1")
		require.Nil(t, val)
		require.True(t, ok)

		wasInCache = c.Set("key2", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("key3", 300)
		require.False(t, wasInCache)
		wasInCache = c.Set("key4", 400)
		require.False(t, wasInCache)

		// key1 key has been removed from cache
		val, ok = c.Get("key1")
		require.Nil(t, val)
		require.False(t, ok)

		// key2 and key3 exist
		val, ok = c.Get("key2")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("key3")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("key4")
		require.True(t, ok)
		require.Equal(t, 400, val)

		// Test removing moved item:
		wasInCache = c.Set("key3", 300)
		require.True(t, wasInCache)
		c.Set("key5", 500)
		c.Set("key6", 600)
		c.Set("key7", 600)

		val, ok = c.Get("key3")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		testDataMap := map[Key]string{
			"key1": "val1",
			"key2": "val2",
			"key3": "val3",
			"key4": "val4",
			"key5": "val5",
		}

		for key, val := range testDataMap {
			c.Set(key, val)
		}

		for key, val := range testDataMap {
			result, ok := c.Get(key)

			require.True(t, ok)
			require.Equal(t, result, val)
		}

		c.Clear()

		for key := range testDataMap {
			result, ok := c.Get(key)

			require.False(t, ok)
			require.Equal(t, result, nil)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	success := true
	c := NewCache(10)
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

	require.True(t, success)
}
