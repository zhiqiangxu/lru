package lru

import (
	"reflect"
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestLRUCache(t *testing.T) {

	// without active gc
	c := NewCache(100)
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
	c = NewCache(100, WithGCInterval(1))
	c.Add("k1", "v1", 0)
	assert.Assert(t, c.Len() == 1)
	c.Add("k2", "v2", 1)
	assert.Assert(t, c.Len() == 2)
	time.Sleep(time.Second * 2)
	_, ok = c.Get("k2")
	assert.Assert(t, !ok && c.Len() == 1)
	c.Close()

	// test cap
	c = NewCache(2, WithGCInterval(1))
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
	c.CompareAndSet("k2", func(value interface{}, exists bool, txn Txn) {
		assert.Assert(t, exists && value.(string) == "v2")
		isNew := txn.Add("k2", 2, 0)
		assert.Assert(t, !isNew)
	})
	v, ok := c.Get("k2")
	assert.Assert(t, ok && v.(int) == 2)

	// test Range
	m := map[interface{}]interface{}{1: 2, 3: 4, 5: 6}
	c = NewCache(20, WithGCInterval(1))
	var keys []Key
	for k, v := range m {
		keys = append(keys, k)
		c.Add(k, v, 0)
	}
	resultMap := make(map[interface{}]interface{})
	var resultKeys []Key
	c.Range(func(k Key, v interface{}, exipreTime int64) bool {
		resultMap[k] = v
		resultKeys = append(resultKeys, k)
		return true
	})
	assert.Assert(t, reflect.DeepEqual(m, resultMap))
	reverseSlice := func(keys []Key) []Key {
		result := make([]Key, len(keys))
		for i := 0; i < len(keys); i++ {
			result[i] = keys[len(keys)-1-i]
		}
		return result
	}
	assert.Assert(t, reflect.DeepEqual(keys, reverseSlice(resultKeys)))

	// Test Reverse
	resultMap = make(map[interface{}]interface{})
	resultKeys = nil
	c.Reverse(func(k Key, v interface{}, exipreTime int64) bool {
		resultMap[k] = v
		resultKeys = append(resultKeys, k)
		return true
	})
	assert.Assert(t, reflect.DeepEqual(m, resultMap))
	assert.Assert(t, reflect.DeepEqual(keys, resultKeys))

}

func TestCacheExpire(t *testing.T) {
	c := NewCache(10).(*cache)

	for i := 0; i < 2000; i++ {
		c.debug()
		c.Get(i)
		c.debug()
		c.Add(i, i, 3)
		c.debug()
		c.Get(i)
		c.debug()
	}
}
