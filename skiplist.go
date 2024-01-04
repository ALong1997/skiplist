package skiplist

import (
	"math/rand"
	"time"

	"golang.org/x/exp/constraints"
)

type (
	SkipList[O constraints.Ordered, T any] struct {
		maxLevel int
		head     *node[O, T]
		r        *rand.Rand
	}
)

func NewSkipList[O constraints.Ordered, T any](maxLevel int) *SkipList[O, T] {
	if maxLevel <= 0 {
		return nil
	}

	return &SkipList[O, T]{
		head:     &node[O, T]{},
		maxLevel: maxLevel,
		r:        rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// Level returns level of the *SkipList.
func (sl *SkipList[O, T]) Level() int {
	if sl.head == nil || sl.head.nextNodes == nil {
		return 0
	}
	return len(sl.head.nextNodes)
}

func (sl *SkipList[O, T]) Get(key O) (val T, exist bool) {
	if sl.Level() == 0 {
		return
	}

	if n := sl.get(key); n != nil {
		return n.val, true
	}
	return
}

func (sl *SkipList[O, T]) Put(key O, val T) {
	if sl.Level() == 0 {
		return
	}

	n := sl.get(key)
	if n != nil {
		// update
		n.val = val
		return
	}

	// randomly determined level
	var randL = sl.randLevel()

	// grow
	for sl.Level()-1 < randL {
		sl.head.nextNodes = append(sl.head.nextNodes, nil)
	}

	n = newNode(key, val, make([]*node[O, T], randL+1))

	head := sl.head
	for l := len(head.nextNodes) - 1; l >= 0; l-- {
		for head.nextNodes[l] != nil && head.key < key {
			// search to the right
			head = head.nextNodes[l]
		}

		// insert
		n.nextNodes[l] = head.nextNodes[l]
		head.nextNodes[l] = n

		// search down
	}
}

func (sl *SkipList[O, T]) Del(key O) {
	if sl.Level() == 0 {
		return
	}

	n := sl.get(key)
	if n != nil {
		// not exist
		return
	}

	head := sl.head
	for l := len(head.nextNodes) - 1; l >= 0; l-- {
		for head.nextNodes[l] != nil && head.nextNodes[l].key < key {
			// search to the right
			head = head.nextNodes[l]
		}

		if head.nextNodes[l] != nil && head.nextNodes[l].key == key {
			// delete
			head.nextNodes[l] = head.nextNodes[l].nextNodes[l]
		}

		// search down
	}

	// cut
	var dif int
	for l := sl.Level() - 1; l >= 0; l-- {
		if sl.head.nextNodes[l] != nil {
			break
		}
		dif++
	}
	sl.head.nextNodes = sl.head.nextNodes[:sl.Level()-dif]
}

// Range searches the *KvPair of key in [start, end].
func (sl *SkipList[O, T]) Range(start, end O) []*KvPair[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	var res = make([]*KvPair[O, T], 0)

	// starting point
	ceilingNode := sl.ceil(start)
	if ceilingNode == nil {
		return res
	}

	// range
	for n := ceilingNode; n != nil && n.key <= end; n = n.nextNodes[0] {
		res = append(res, newKvPair(n.key, n.val))
	}
	return res
}

// Ceil returns *KvPair of the least key greater than or equal to target.
func (sl *SkipList[O, T]) Ceil(target O) (*KvPair[O, T], bool) {
	if sl.Level() == 0 {
		return nil, false
	}

	if ceilingNode := sl.ceil(target); ceilingNode != nil {
		return newKvPair(ceilingNode.key, ceilingNode.val), true
	}
	return nil, false
}

// Floor returns *KvPair of the greatest key less than or equal to target.
func (sl *SkipList[O, T]) Floor(target O) (*KvPair[O, T], bool) {
	if sl.Level() == 0 {
		return nil, false
	}

	if floorNode := sl.floor(target); floorNode != sl.head.nextNodes[0] {
		return newKvPair(floorNode.key, floorNode.val), true
	}
	return nil, false
}

func (sl *SkipList[O, T]) get(key O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	head := sl.head
	for l := len(head.nextNodes) - 1; l >= 0; l-- {
		for head.nextNodes[l] != nil && head.nextNodes[l].key < key {
			// search to the right
			head = head.nextNodes[l]
		}

		if head.nextNodes[l] != nil && head.nextNodes[l].key == key {
			// exist
			return head.nextNodes[l]
		}

		// search down
	}
	// not exist
	return nil
}

func (sl *SkipList[O, T]) ceil(target O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	head := sl.head
	for l := len(head.nextNodes) - 1; l >= 0; l-- {
		for head.nextNodes[l] != nil && head.nextNodes[l].key < target {
			// search to the right
			head = head.nextNodes[l]
		}

		if head.nextNodes[l] != nil && head.nextNodes[l].key == target {
			// equal
			return head.nextNodes[l]
		}

		// search down
	}
	// head.nextNodes[0] is ceil || head.nextNodes[0] == nil(ceil is not exist)
	return head.nextNodes[0]
}

func (sl *SkipList[O, T]) floor(target O) *node[O, T] {
	if sl.Level() == 0 {
		return nil
	}

	head := sl.head
	for l := len(head.nextNodes) - 1; l >= 0; l-- {
		for head.nextNodes[l] != nil && head.nextNodes[l].key < target {
			// search to the right
			head = head.nextNodes[l]
		}

		if head.nextNodes[l] != nil && head.nextNodes[l].key == target {
			// equal
			return head.nextNodes[l]
		}

		// search down
	}
	// head is floor || head == sl.head.nextNodes[0](floor is not exist)
	return head
}

func (sl *SkipList[O, T]) randLevel() int {
	var randL int
	for rand.Intn(2) == 0 && randL < sl.maxLevel {
		randL++
	}
	return randL
}
