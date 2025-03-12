package lru

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLRUCache(t *testing.T) {
	t.Parallel()

	const maxSize = 10

	t.Run("just works", func(t *testing.T) {
		cache := NewLRUCache[int, int](maxSize/2, time.Second)

		for i := range maxSize {
			cache.Put(i, i)
		}

		for i := range maxSize {
			v, ok := cache.Get(i)

			if i < maxSize/2 {
				require.False(t, ok, fmt.Sprint(i))
			} else {
				require.True(t, ok)
				require.Equal(t, v, i)
			}
		}
	})

	t.Run("displacement works", func(t *testing.T) {
		t.Parallel()

		cache := NewLRUCache[int, int](maxSize, time.Second)
		cache.Put(1, 1)

		time.Sleep(time.Second)

		_, ok := cache.Get(1)
		require.False(t, ok)
	})

	t.Run("concurrent", func(t *testing.T) {
		t.Parallel()

		var (
			cache = NewLRUCache[int, int](maxSize, time.Second)
			wg    = &sync.WaitGroup{}
		)

		for i := range maxSize * 3 {
			wg.Add(1)

			go func() {
				defer wg.Done()

				cache.Put(i, i)
			}()
		}

		wg.Wait()

		require.Equal(t, maxSize, cache.len)
		time.Sleep(time.Second)

		for i := maxSize * 2; i < maxSize*3; i++ {
			_, ok := cache.Get(i)

			require.False(t, ok)
		}
	})

	t.Run("add multiple times", func(t *testing.T) {
		t.Parallel()

		cache := NewLRUCache[int, int](maxSize, time.Second)

		cache.Put(1, 1)

		v, ok := cache.Get(1)
		require.True(t, ok)
		require.Equal(t, 1, v)

		time.Sleep(time.Second / 2)

		cache.Put(1, 1)
		time.Sleep(time.Second / 2)

		v, ok = cache.Get(1)
		require.True(t, ok)
		require.Equal(t, 1, v)
	})

	t.Run("random access", func(t *testing.T) {
		t.Parallel()

		var (
			cache = NewLRUCache[int, int](maxSize, time.Second)
			wg    = sync.WaitGroup{}
		)

		for i := range maxSize * 10 {
			wg.Add(1)

			go func() {
				defer wg.Done()

				cache.Put(i%maxSize, i)
			}()
		}

		wg.Wait()

		time.Sleep(time.Second / 2)

		for range maxSize {
			key := rand.Intn(maxSize)

			_, ok := cache.Get(key)
			require.True(t, ok)
		}

		time.Sleep(time.Second / 2)

		for range maxSize {
			key := rand.Intn(maxSize)

			_, ok := cache.Get(key)
			require.False(t, ok)
		}
	})
}
