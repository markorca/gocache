package cache

type CacheItem struct {
	// Key is the CacheItem's key (250 bytes maximum because of memcached).
	Key string

	// Value is the CacheItem's value.
	Value []byte

	// Expiration is the cache expiration time, in seconds
	// Zero means no expiration time.
	Expiration int32
}