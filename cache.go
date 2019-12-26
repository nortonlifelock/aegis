package ttl

import (
	"context"
	"sync"
	"time"
)

// Cache is a ttl cache which will retain an entry for the duration of
// of the ttl time.Duration that is passed in for each entry. After the
// ttl has expired the entry is removed from the cache.
type Cache struct {
	data   sync.Map
	ctxmap sync.Map
	mutex  sync.Mutex
}

// Delete removes an entry from the cache using the entry key
func (c *Cache) Delete(key interface{}) {
	c.data.Delete(key)
}

// Load loads the value from the cache that is stored using key
func (c *Cache) Load(key interface{}) (value interface{}, ok bool) {
	return c.data.Load(key)
}

// Store stores the input value using the input key with a ttl specified as a time.Duration pointer
// a nil duration will store the value forever
func (c *Cache) Store(ctx context.Context, key, value interface{}, ttl *time.Duration) {
	// cancel the old record to replace it
	c.cancel(key)

	c.data.Store(key, value)

	if ttl != nil {
		go c.ttl(ctx, key, ttl)
	}
}

// LoadOrStore stores the value at the key with a cache deletion after the ttl time.Duration expires. If the value for the key
// is already stored though then the boolean return "loaded" will indicate if it was loaded. An overlapping key entry does not
// store the input value if it already exists in the cache
func (c *Cache) LoadOrStore(ctx context.Context, key, value interface{}, ttl *time.Duration) (actual interface{}, loaded bool) {
	actual, loaded = c.data.LoadOrStore(key, value)

	// Create the ttl go routine only if the record is new
	if !loaded {
		go c.ttl(ctx, key, ttl)
	}

	return actual, loaded
}

// cancel executes the cancel function for the specified key, removing it from the cache
func (c *Cache) cancel(key interface{}) {

	// Execute the cancellation function if it exists
	if cancel, exists := c.ctxmap.Load(key); exists {
		if cancel != nil {
			if cfunc, ok := cancel.(context.CancelFunc); ok {
				cfunc()
				c.ctxmap.Delete(key)
			}
		}
	}
}

// ttl uses the context with timeout to clear an entry from the map when the context
// expires for the passed key
func (c *Cache) ttl(ctx context.Context, key interface{}, ttl *time.Duration) {

	if ttl != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *ttl)

		// Store the cancellation
		if _, loaded := c.ctxmap.LoadOrStore(key, cancel); !loaded {

			// Block on the context
			select {
			case <-ctx.Done():
				c.data.Delete(key)
			}
		}
	}
}
