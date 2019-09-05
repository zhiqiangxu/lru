# lru

an easy to use concurrent safe lru cache with expire feature backed up by home made skiplist!


```golang

package lru

import (
	"testing"
	"time"

	"github.com/zhiqiangxu/lru"
	"gotest.tools/assert"
)

func TestLRUCache(t *testing.T) {

	// without active gc
	c := lru.NewCache(100, 0, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok := c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 2)

	// with active gc every second
	c = lru.NewCache(100, 1, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok = c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 1)

	// test cap
	c = lru.NewCache(2, 1, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 0)
	assert.Assert(t, c.Len() == 2)
	c.Add("k3", "v3", 0)
	assert.Assert(t, c.Len() == 2)
	_, ok = c.Get("k1")
	assert.Assert(t, !ok && c.Len() == 2)
	_, ok = c.Get("k2")
	assert.Assert(t, ok)
}

```
