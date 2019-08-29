## About
A distributed local cache for Go.

When SetObject/DeleteObject is called, the same action will be performed in each cache node.

##Usage:


``` 
package main

import "github.com/markorca/gocache"
import "fmt"

func main() {
    m := map[string]int{"one":1, "two":2, "three":3}
    // m := "hello world 3"

    localCache := cache.Init()

    localCache.SetObject("test", m, 0)

    test, _ := localCache.GetObject("test")
    fmt.Println(test)

    localCache.DeleteObject("test")
}

```