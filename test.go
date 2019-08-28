package main

import "./cache"
import "fmt"

func main() {
	localCache := cache.Init()
	localCache.SetObject("test", 123, 0)
	test := localCache.GetObject("test")
	fmt.Println(test)
}