package cache

type LocalStorage interface {
	// set cacheItem in local storage
	Set(cacheItem *CacheItem) (err error)
	
	// use key to get cacheItem in local storage
	Get(key string) (cacheItem *CacheItem, err error)
	
	// // delete cacheItem in local storage
	// Delete(key string) (err error)
}