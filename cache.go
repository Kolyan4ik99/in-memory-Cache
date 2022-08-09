package inMemoryCache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type CacheInterface interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, error)
	Delete(key string) error
}

// на каждый уникальный ключ создается горутина, которая ждёт ttl
// ctx добавлен чтобы завершить горутину, если нужно обновить или удалить ключ
type valueItem struct {
	ttl       time.Duration
	item      interface{}
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type cache struct {
	values map[string]valueItem
	mu     sync.RWMutex
}

func (c *cache) cleanValues(key string, ttl time.Duration) {
	select {
	case <-c.values[key].ctx.Done():
		return
	case <-time.After(ttl):
		c.mu.Lock()
		delete(c.values, key)
		c.mu.Unlock()
	}
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanIfExist(key)
	ctx, ctxCancel := context.WithCancel(context.TODO())
	c.values[key] = valueItem{
		item:      value,
		ttl:       ttl,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}

	go c.cleanValues(key, ttl)
}

func (c *cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	el, err := c.findEl(key)
	if err == nil {
		return el, err
	}
	return nil, err
}

func (c *cache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err := c.findEl(key)
	if err != nil {
		c.values[key].ctxCancel()
		delete(c.values, key)
		return err
	}

	return err
}

func New() *cache {
	return &cache{
		values: make(map[string]valueItem),
	}
}

func (c *cache) cleanIfExist(key string) {
	item, exist := c.values[key]
	if exist {
		item.ctxCancel()
		time.Sleep(time.Second)
	}
}

// If element exist err == nil
func (c *cache) findEl(item string) (interface{}, error) {
	foundItem, exist := c.values[item]
	if !exist {
		err := errors.New(fmt.Sprintf("item not found, key = [%s]", item))
		return nil, err
	}
	return foundItem.item, nil
}
