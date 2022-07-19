
# inMemoryCache

Module for cache your program


## Installation

Write commands in your go project

```bash
  go get -u github.com/Kolyan4ik99/inMemoryCache
  go mod download
```

## Example


	package main

	import (
	  "fmt"
	  "github.com/Kolyan4ik99/inMemoryCache"
	)

	func main() {

    cache := inMemoryCache.New()
    
    cache.Set("User1", 5000)
    
    value, err := cache.Get("User1")
    
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(value)
    
    cache.Set("Vasya", 1592)
    
    value, err = cache.Get("Vasya")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(value)
    
    value, err = cache.Get("User1")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(value)

    err = cache.Delete("User1")
    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Item successful delete")
    }

    value, err = cache.Get("User1")
    if err != nil {
      fmt.Println(err)
    }
    fmt.Println(value)

}