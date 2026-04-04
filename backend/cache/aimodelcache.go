package cache

import "sync"

// AIModelCache 用于缓存模型的智能度测试结果
type AIModelCache struct {
	sync.RWMutex
	results map[string]string
}

// GlobalCache 全局缓存实例
var GlobalCache = &AIModelCache{
	results: make(map[string]string),
}

// Get 获取缓存结果
func (c *AIModelCache) Get(key string) (string, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.results[key]
	return val, ok
}

// Set 设置缓存结果
func (c *AIModelCache) Set(key, val string) {
	c.Lock()
	defer c.Unlock()
	c.results[key] = val
}

// GenerateCacheKey 生成唯一的缓存键
func GenerateCacheKey(protocol, apiKey, model string) string {
	return protocol + ":" + apiKey + ":" + model
}
