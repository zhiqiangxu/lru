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
type SkipList interface {
	Add(key Key, value interface{})
	Get(key Key) (value interface{}, ok bool)
	Remove(key Key)
	Head() (value interface{}, ok bool)
	Tail() (value interface{}, ok bool)
}
