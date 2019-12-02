package lru

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestLRUCache(t *testing.T) {

	// without active gc
	c := NewCache(100, 0, nil)
	ret := c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1 && ret)
	ret = c.Add("k1", "v1", 0)
	assert.Assert(t, !ret)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	ret = c.Add("k2", "v2", 1)
	assert.Assert(t, !ret)
	time.Sleep(time.Second * 2)
	_, ok := c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 2)
	ret = c.Add("k2", "v2", 1)
	assert.Assert(t, ret)

	// with active gc every second
	c = NewCache(100, 1, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok = c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 1)

	// test cap
	c = NewCache(2, 1, nil)
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
	// test CompareAndSet
	c.CompareAndSet("k2", func(value interface{}, exists bool) {
		assert.Assert(t, exists && value.(string) == "v2")
		isNew := c.AddLocked("k2", 2, 0)
		assert.Assert(t, !isNew)
	})
	v, ok := c.Get("k2")
	assert.Assert(t, ok && v.(int) == 2)
}
