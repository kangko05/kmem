package database

import (
	"context"
	"strings"
	"sync"
	"time"
)

type cacheItem struct {
	key        string
	value      any
	insertedAt int64 // unix time
	permanent  bool
}

func (ci cacheItem) Value() any {
	return ci.value
}

type Cache struct {
	ctx         context.Context
	items       sync.Map
	cleanupTime time.Duration
}

func NewCache(ctx context.Context, cleanupTime time.Duration) *Cache {
	c := &Cache{
		ctx:         ctx,
		items:       sync.Map{},
		cleanupTime: cleanupTime,
	}

	go c.start()

	return c
}

func (c *Cache) start() {
	ticker := time.NewTicker(c.cleanupTime)

	for {
		select {
		case <-c.ctx.Done():
			c.items.Clear() // necessary?
			return
		case <-ticker.C:
			cleanList := []string{}

			c.items.Range(func(key, value any) bool {
				ci, ok := value.(cacheItem)
				if !ok {
					cleanList = append(cleanList, key.(string))
					return true
				}

				if ci.permanent {
					return true
				}

				if time.Duration(time.Now().Unix()-ci.insertedAt) > c.cleanupTime {
					cleanList = append(cleanList, key.(string))
				}

				return true
			})

			for _, key := range cleanList {
				c.Delete(key)
			}
		}
	}
}

func (c *Cache) AddPermanent(key string, val any) {
	c.items.Store(key, cacheItem{
		key:        key,
		value:      val,
		insertedAt: time.Now().Unix(),
		permanent:  true,
	})
}

func (c *Cache) Add(key string, val any) {
	c.items.Store(key, cacheItem{
		key:        key,
		value:      val,
		insertedAt: time.Now().Unix(),
		permanent:  false,
	})
}

func (c *Cache) Get(key any) (any, bool) {
	citem, exists := c.items.Load(key)
	if !exists {
		return nil, exists
	}

	return citem.(cacheItem).Value(), exists
}

func (c *Cache) Delete(key any) {
	c.items.Delete(key)
}

func (c *Cache) InvalidateUserCache(username string) {
	c.items.Range(func(key, value any) bool {
		k, ok := key.(string)
		if !ok {
			return true
		}

		if strings.Contains(k, username) {
			c.items.Delete(key)
		}

		return true
	})
}
