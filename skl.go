package lru

import (
	"math"
	"math/rand"
	"time"
)

const (
	// DefaultMaxLevel for skl
	DefaultMaxLevel int = 18
	// DefaultProbability for skl
	DefaultProbability float64 = 1 / math.E
)

type elementNode struct {
	next []*element
}

type element struct {
	elementNode
	key   int64
	value interface{}
}

type skl struct {
	elementNode
	maxLevel       int
	length         int
	randSource     rand.Source
	probability    float64
	probTable      []float64
	prevNodesCache []*elementNode
}

// NewSkipList creates a new SkipList
func NewSkipList() SkipList {
	return NewSkipListWithMaxLevel(DefaultMaxLevel)
}

// NewSkipListWithMaxLevel creates a new SkipList with specified maxLevel
func NewSkipListWithMaxLevel(maxLevel int) SkipList {
	return &skl{
		elementNode:    elementNode{next: make([]*element, maxLevel)},
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    DefaultProbability,
		probTable:      probabilityTable(DefaultProbability, maxLevel),
		prevNodesCache: make([]*elementNode, maxLevel),
	}
}

func probabilityTable(probability float64, maxLevel int) (table []float64) {
	for i := 0; i < maxLevel; i++ {
		prob := math.Pow(probability, float64(i))
		table = append(table, prob)
	}
	return table
}

func (s *skl) Add(key int64, value interface{}) {
	prevs := s.getPrevElementNodes(key)
	ele := prevs[0].next[0]
	if ele != nil && ele.key <= key {
		ele.value = value
		return
	}

	ele = &element{
		elementNode: elementNode{next: make([]*element, s.randLevel())},
		key:         key,
		value:       value,
	}

	for i := range ele.next {
		ele.next[i] = prevs[i].next[i]
		prevs[i].next[i] = ele
	}

	s.length++
}

func (s *skl) randLevel() (level int) {
	r := float64(s.randSource.Int63()) / (1 << 63)

	level = 1
	for level < s.maxLevel && r < s.probTable[level] {
		level++
	}
	return
}

// 找到每一层上毗邻于该key对应元素之前的节点
func (s *skl) getPrevElementNodes(key int64) []*elementNode {
	var prev = &s.elementNode
	var current *element

	prevs := s.prevNodesCache
	for i := s.maxLevel - 1; i >= 0; i-- {
		current = prev.next[i]

		for current != nil && current.key < key {
			prev = &current.elementNode
			current = current.next[i]
		}

		prevs[i] = prev
	}

	return prevs
}

func (s *skl) Get(key int64) (value interface{}, ok bool) {
	prev := &s.elementNode
	var current *element
	for i := s.maxLevel - 1; i >= 0; i-- {
		current = prev.next[i]
		for current != nil && current.key < key {
			prev = &current.elementNode
			current = current.next[i]
		}
	}

	if current != nil && current.key <= key {
		return current.value, true
	}

	return nil, false
}

func (s *skl) Remove(key int64) {
	prevs := s.getPrevElementNodes(key)
	if ele := prevs[0].next[0]; ele != nil && ele.key <= key {

		for i, iele := range ele.next {
			prevs[i].next[i] = iele
		}
		s.length--
	}
}

func (s *skl) Head() (key int64, value interface{}, ok bool) {
	if s.next[0] != nil {
		key = s.next[0].key
		value = s.next[0].value
		ok = true
	}
	return
}
