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
	c := lru.NewCache(
		100/*最多保存的记录数*/, 
		0/*主动GC的周期，单位秒，0表示不主动GC，到达上限后以一进一出的方式进行淘汰*/, 
		nil/*淘汰时的回调函数，不需要填nil*/)
	c.Add("k1", "v1", 0/*没有过期时间*/)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1/*1秒钟后过期*/)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok := c.Get("k2")
	// 虽然已经过期，但是由于没有主动回收，坑位仍然占着
	assert.Assert(t, !ok && c.Len() == 2)

	// with active gc every second
	c = lru.NewCache(
		100/*最多保存的记录数*/, 
		1/*每隔1秒钟尝试主动GC，释放坑位*/, 
		nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok = c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 1)

	// test cap
	c = lru.NewCache(2/*最多保存2条记录*/, 1, nil)
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 0)
	assert.Assert(t, c.Len() == 2)
	//由于最多保存2条，当插入第3条时会淘汰最老的k1
	c.Add("k3", "v3", 0)
	assert.Assert(t, c.Len() == 2)
	_, ok = c.Get("k1")
	// 确认k1已不存在
	assert.Assert(t, !ok && c.Len() == 2)
	_, ok = c.Get("k2")
	assert.Assert(t, ok)
}

```
