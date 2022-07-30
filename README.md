
# inMemoryCache

Module for cache your program


## Installation

Write commands in your go project

```bash
  go get -u github.com/Kolyan4ik99/inMemoryCache
  go mod download
```

## Example 1

    cache.Set("userId", 42, time.Second * 5)

    userId, err := cache.Get("userId")
    if err != nil { // err == nil
        log.Fatal(err)
    }
    fmt.Println(userId) // Output: 42
    
    time.Sleep(time.Second * 6) // прошло 5 секунд
    
    userId = cache.Get("userId")
    userId, err := cache.Get("userId")
    if err != nil { // err != nil
        log.Fatal(err) // сработает этот код
    }

## Example 2

    cache := inMemoryCache.New()

    cache.Set("userId", 42, time.Second*5)
    userId, err := cache.Get("userId")
    if err != nil { // err == nil
        log.Fatal(err)
    }
    fmt.Println(userId) // Output: 42

    cache.Set("userId", 52, time.Second*8)

    time.Sleep(time.Second * 7) // прошло 6 секунд
    userId, err = cache.Get("userId")
    if err != nil { // err == nil
        log.Fatal(err)
    }
    fmt.Println(userId) // Output: 52
    