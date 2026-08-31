package pokecache

import (
    "testing"
    "time"
)

func TestCreateCache(t *testing.T) {
    interval := time.Millisecond * 10
    cache := NewCache(interval)
    if cache.cache == nil {
        t.Error("cache is nil")
    }
}
