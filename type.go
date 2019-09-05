package lru

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

// Cache for lru interface
type Cache interface {
	// expireSeconds 0 means never expire
	Add(key Key, value interface{}, expireSeconds int)
	Get(key Key) (value interface{}, ok bool)
	Remove(key Key)
	Len() int
}

// SkipList for skl interface
// TODO extend key type when go supports generic
type SkipList interface {
	Add(key int64, value interface{})
	Get(key int64) (value interface{}, ok bool)
	Remove(key int64)
	Head() (value interface{}, ok bool)
}
