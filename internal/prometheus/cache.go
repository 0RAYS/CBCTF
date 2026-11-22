package prometheus

import (
	"sync"
	"time"
)

// collectionTTL 是数据库型采集器的缓存窗口。抓取本身不应成为平台的负载源，
// 但指标又必须反映当前状态，因此取一个远小于常见抓取间隔的窗口。
const collectionTTL = 30 * time.Second

// cacheMemo 缓存一次采集的结果。采集失败时保留上一次的成功值，
// 由 collection_success 指标单独反映采集状态。
type cacheMemo[T any] struct {
	mu        sync.Mutex
	value     T
	expiresAt time.Time
}

func (c *cacheMemo[T]) get(collect func() (T, bool)) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if time.Now().Before(c.expiresAt) {
		return c.value, true
	}
	value, ok := collect()
	if !ok {
		return c.value, false
	}
	c.value = value
	c.expiresAt = time.Now().Add(collectionTTL)
	return value, true
}
