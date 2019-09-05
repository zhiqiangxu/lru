package lru

import (
	"container/list"
	"sync"
	"time"
)

type cache struct {
	rwLock           sync.RWMutex
	maxEntries       int
	gcIntervalSecond int
	onEvicted        func(key Key, value interface{})
	ll               *list.List
	cache            map[interface{}]*list.Element
	toExpire         SkipList
}

type entry struct {
	key       Key
	value     interface{}
	timeoutTS int64
}

type sklEntry struct {
}

// NewCache creates a new Cache
func NewCache(maxEntries, gcIntervalSecond int, onEvicted func(key Key, value interface{})) Cache {
	c := &cache{
		maxEntries:       maxEntries,
		gcIntervalSecond: gcIntervalSecond,
		onEvicted:        onEvicted,
		ll:               list.New(),
		cache:            make(map[interface{}]*list.Element),
		toExpire:         NewSkipList()}

	if gcIntervalSecond > 0 {
		go c.gc()
	}
	return c
}

func (c *cache) gc() {
	ticker := time.NewTimer(time.Duration(c.gcIntervalSecond) * time.Second)
	for {
		select {
		case <-ticker.C:
			nowSecond := time.Now().Unix()
			c.rwLock.RLock()
			timeoutTS, expireValues, ok := c.toExpire.Head()
			if !ok || timeoutTS > nowSecond {
				c.rwLock.RUnlock()
				break
			}
			c.rwLock.RUnlock()
			c.rwLock.Lock()
			for {
				timeoutTS, expireValues, ok = c.toExpire.Head()
				if !ok || timeoutTS > nowSecond {
					break
				}
				for key := range expireValues.(map[interface{}]struct{}) {
					if ele, hit := c.cache[key]; hit {
						c.removeElement(ele)
					} else {
						c.rwLock.Unlock()
						panic("bug in cache3")
					}

				}
			}
			c.rwLock.Unlock()
		}
	}
}

func (c *cache) Add(key Key, value interface{}, expireSeconds int) {

	// expireSeconds < 0
	if expireSeconds < 0 {
		c.Remove(key)
		return
	}

	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	var oldTimeoutTS, timeoutTS int64
	if expireSeconds > 0 {
		timeoutTS = int64(expireSeconds) + time.Now().Unix()
	}

	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		oldTimeoutTS = ee.Value.(*entry).timeoutTS
		ee.Value.(*entry).timeoutTS = timeoutTS
	} else {
		ele := c.ll.PushFront(&entry{key, value, timeoutTS})
		c.cache[key] = ele
		if c.maxEntries != 0 && c.ll.Len() > c.maxEntries {
			c.remove1ExpiredOrOldest()
		}
	}

	if oldTimeoutTS == timeoutTS {
		return
	}

	if oldTimeoutTS != 0 {
		c.removeKeyFromExpire(key, oldTimeoutTS)
	}

	if timeoutTS == 0 {
		return
	}

	c.addKeyToExpire(key, timeoutTS)
}

func (c *cache) Get(key Key) (value interface{}, ok bool) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	if ele, hit := c.cache[key]; hit {
		timeoutTS := ele.Value.(*entry).timeoutTS
		if timeoutTS == 0 {
			c.ll.MoveToFront(ele)
			return ele.Value.(*entry).value, true
		}
		if timeoutTS >= time.Now().Unix() {
			c.ll.MoveToFront(ele)
			return ele.Value.(*entry).value, true
		}
	}
	return nil, false
}

func (c *cache) Len() int {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()

	return c.ll.Len()
}

func (c *cache) Remove(key Key) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *cache) remove1ExpiredOrOldest() {
	timeoutTS, expireValues, ok := c.toExpire.Head()
	if !ok || timeoutTS > time.Now().Unix() {
		ele := c.ll.Back()
		if ele != nil {
			c.removeElement(ele)
		}
		return
	}
	for key := range expireValues.(map[interface{}]struct{}) {
		if ele, hit := c.cache[key]; hit {
			c.removeElement(ele)
			return
		}
		panic("bug in cache2")
	}
}

func (c *cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	entry := e.Value.(*entry)
	delete(c.cache, entry.key)
	if entry.timeoutTS != 0 {
		c.removeKeyFromExpire(entry.key, entry.timeoutTS)
	}
	if c.onEvicted != nil {
		c.onEvicted(entry.key, entry.value)
	}
}

func (c *cache) removeKeyFromExpire(key Key, timeoutTS int64) {
	value, ok := c.toExpire.Get(timeoutTS)
	if !ok {
		panic("bug in cache")
	}
	oldExpireValues := value.(map[interface{}]struct{})
	delete(oldExpireValues, key)
	if len(oldExpireValues) == 0 {
		c.toExpire.Remove(timeoutTS)
	}
}

func (c *cache) addKeyToExpire(key Key, timeoutTS int64) {
	var expireValues map[interface{}]struct{}
	value, ok := c.toExpire.Get(timeoutTS)
	if ok {
		expireValues = value.(map[interface{}]struct{})
	} else {
		expireValues = make(map[interface{}]struct{})
	}

	expireValues[key] = struct{}{}

	if !ok {
		c.toExpire.Add(timeoutTS, expireValues)
	}
}
