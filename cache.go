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
type cache struct {
	values map[string]interface{}
	ttl    map[string]ctxWithCancel
	mu     sync.RWMutex
	ctx    context.Context
}

type ctxWithCancel struct {
	context.Context
	context.CancelFunc
}

func (c *cache) cleanValues(key string, ttl time.Duration) {
	select {
	case <-c.ttl[key].Done():
		return
	case <-time.After(ttl):
		c.mu.Lock()
		delete(c.values, key)
		delete(c.ttl, key)
		c.mu.Unlock()
	}
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
	ctxCancel, exist := c.ttl[key]
	if exist {
		ctxCancel.CancelFunc()
		time.Sleep(time.Second)
	}
	ctx, cancel := context.WithCancel(c.ctx)
	c.ttl[key] = ctxWithCancel{ctx, cancel}
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
	if err == nil {
		c.ttl[key].CancelFunc()
		delete(c.values, key)
		// Ожидаем когда контекст завершится
		time.Sleep(time.Millisecond)
		delete(c.ttl, key)
	}
	return err
}

func New() *cache {
	return &cache{
		values: make(map[string]interface{}),
		ttl:    make(map[string]ctxWithCancel),
		ctx:    context.Background(),
	}
}

// If element exist err == nil
func (c *cache) findEl(item string) (foundItem interface{}, err error) {
	foundItem, exist := c.values[item]
	if !exist {
		err = errors.New(fmt.Sprintf("item not found, key = [%s]", item))
		return nil, err
	}
	return foundItem, nil
}
