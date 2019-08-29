package cache

type LocalStorage interface {
	// set cacheItem in local storage
	Set(cacheItem *CacheItem) (error)
	
	// use key to get cacheItem in local storage
	Get(key string) (*CacheItem, error)
	
	// delete cacheItem in local storage
	Delete(key string) (error)
}