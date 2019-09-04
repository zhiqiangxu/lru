package lru

type skl struct {
}

// NewSkipList creates a new SkipList
func NewSkipList() SkipList {
	return &skl{}
}

func (s *skl) Add(key int64, value interface{}) {

}

func (s *skl) Get(key int64) (value interface{}, ok bool) {
	return
}

func (s *skl) Remove(key int64) {

}

func (s *skl) Head() (value interface{}, ok bool) {
	return
}

func (s *skl) Tail() (value interface{}, ok bool) {
	return
}
