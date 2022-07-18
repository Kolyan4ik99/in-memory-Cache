package inMemoryCache

import (
	"errors"
	"fmt"
	"math"
)

type CacheInterface interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, error)
	Delete(key string) error
}

func (c *Cache) Set(key string, value interface{}) {
	_, err := c.findEl(key)
	if err != nil {
		if c.len == c.capacity {
			c.flush()
		}
		c.values[key] = value
		c.rate[key]++
		c.len++
	} else {
		c.values[key] = value
	}

}

type Cache struct {
	values   map[string]interface{}
	rate     map[string]int
	capacity int
	len      int
}

func (c *Cache) Get(key string) (interface{}, error) {
	el, err := c.findEl(key)
	if err == nil {
		c.rate[key] = c.rate[key] + 1
		return el, err
	}
	return nil, err
}

func (c *Cache) Delete(key string) error {
	_, err := c.findEl(key)
	if err == nil {
		c.len--
		delete(c.values, key)
		delete(c.rate, key)
	}
	return err
}

func New() *Cache {
	return &Cache{
		values:   make(map[string]interface{}, 50),
		capacity: 50,
		len:      0,
		rate:     make(map[string]int, 50),
	}
}

// If element exist err == nil
func (c *Cache) findEl(item string) (foundItem interface{}, err error) {
	foundItem, exist := c.values[item]
	if !exist {
		err = errors.New(fmt.Sprintf("item not found, key = [%s]", item))
		return nil, err
	}
	return foundItem, nil
}

func (c *Cache) flush() {
	minRate := math.MaxInt32
	for _, valueRate := range c.rate {
		if valueRate < minRate {
			minRate = valueRate
		}
	}

	for key, valueRate := range c.rate {
		if valueRate == minRate {
			delete(c.values, key)
			delete(c.rate, key)
			c.len--
			return
		}
	}

	fmt.Println(len(c.values), c.len)
}
