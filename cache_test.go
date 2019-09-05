package lru

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestLRUCache(t *testing.T) {

	c := NewCache(100, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok := c.Get("k2")
	assert.Assert(t, !ok)
}
