package cache

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type cacheItem struct {
	key     string
	val     any
	hit     int
	expires int64
}

type Cache struct {
	ctx  context.Context
	list sync.Map
	ttl  time.Duration
}

func New(ctx context.Context) *Cache {
	c := &Cache{
		ctx:  ctx,
		list: sync.Map{},
		ttl:  time.Hour,
	}

	go c.run()

	return c
}

func (c *Cache) run() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.cleanList()
		}
	}
}

func (c *Cache) Set(key string, val any) {
	c.list.Store(key,
		cacheItem{
			key:     key,
			val:     val,
			hit:     0,
			expires: time.Now().Add(c.ttl).Unix(),
		})
}

func (c *Cache) Get(key string) (any, bool) {
	val, ok := c.list.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(cacheItem)

	if time.Now().Unix() >= item.expires {
		c.list.Delete(key)
		return nil, false
	}

	item.hit += 1
	c.list.Store(key, item)
	return item.val, true
}

func (c *Cache) InvalidateUserGallery(username string) {
	galleryPrefix := fmt.Sprintf("gallery:%s", username)
	statsPrefix := fmt.Sprintf("%s:stats", username)

	c.list.Range(func(key, value any) bool {
		keyStr := key.(string)
		if strings.HasPrefix(keyStr, galleryPrefix) || strings.HasPrefix(keyStr, statsPrefix) {
			c.list.Delete(key)
		}
		return true
	})
}

func (c *Cache) ClearGalleryCache() {
	c.list.Range(func(key, value any) bool {
		if strings.Contains(key.(string), "gallery:") {
			c.list.Delete(key)
		}

		return true
	})
}

func (c *Cache) cleanList() {
	now := time.Now().Unix()

	var items []cacheItem

	c.list.Range(func(key, value any) bool {
		item := value.(cacheItem)

		if now >= item.expires {
			c.list.Delete(key)
		} else {
			items = append(items, item)
		}

		return true
	})

	sort.Slice(items, func(i, j int) bool {
		return items[i].hit < items[j].hit
	})

	var delKeys []string

	for i := range len(items) / 2 {
		delKeys = append(delKeys, items[i].key)

	}

	for _, key := range delKeys {
		c.list.Delete(key)
	}
}
