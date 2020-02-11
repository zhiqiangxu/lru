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

	View(funcLocked func(rt RTxn))
	Update(funcLocked func(t Txn))
	CompareAndSet(key Key, funcLocked func(value interface{}, exists bool, t Txn))
	Range(funcLocked func(key Key, value interface{}, expireTime int64) bool)
	Reverse(funcLocked func(key Key, value interface{}, expireTime int64) bool)

	Close()
}

// Txn for read/write transaction
type Txn interface {
	Add(key Key, value interface{}, expireSeconds int) (new bool)
	Get(key Key) (value interface{}, ok bool)
	Remove(key Key)
	Len() int
}

// RTxn for read only transaction
type RTxn interface {
	Get(key Key) (value interface{}, ok bool)
	Len() int
}

// SkipList for skl interface
// TODO extend key type when go supports generic
type SkipList interface {
	Add(key int64, value interface{})
	Get(key int64) (value interface{}, ok bool)
	Remove(key int64)
	Head() (key int64, value interface{}, ok bool)
}
