package lru

type skl struct {
}

// NewSkipList creates a new SkipList
func NewSkipList() SkipList {
	return &skl{}
}

func (s *skl) Add(key Key, value interface{}) {

}

func (s *skl) Get(key Key) (value interface{}, ok bool) {
	return
}

func (s *skl) Remove(key Key) {

}

func (s *skl) Head() (value interface{}, ok bool) {
	return
}

func (s *skl) Tail() (value interface{}, ok bool) {
	return
}
