package lru

// A Key may be any value that is comparable. See http://golang.org/ref/spec#Comparison_operators
type Key interface{}

// Cache for lru interface
type Cache interface {
	// expireSeconds 0 means never expire
	Add(key Key, value interface{}, expireSeconds int) (new bool)
	Get(key Key) (value interface{}, ok bool)
	Remove(key Key)
	Len() int
	// 在funcLocked回调内只能调各种Locked方法，否则将死锁
	CompareAndSet(key Key, funcLocked func(value interface{}, exists bool))
	// below are paired with CompareAndSet, use carefully
	GetLocked(key Key) (value interface{}, ok bool)
	AddLocked(key Key, value interface{}, expireSeconds int) (new bool)
	RemoveLocked(key Key)
	LenLocked() int
}

// SkipList for skl interface
// TODO extend key type when go supports generic
type SkipList interface {
	Add(key int64, value interface{})
	Get(key int64) (value interface{}, ok bool)
	Remove(key int64)
	Head() (key int64, value interface{}, ok bool)
}
