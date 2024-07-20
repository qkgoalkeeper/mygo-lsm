package skipList

import (
	"github.com/whuanle/lsm/kv"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	//跳表索引最⼤层数，可根据实际情况进⾏调整
	maxLevel    int     = 18
	probability float64 = 1 / math.E
)

type (
	// Node 跳表节点
	Node struct {
		next []*Element // next指针，分层存储
	}
	// Element 跳表存储元素
	Element struct {
		Node
		KV kv.Value
		//key []byte // 存储 key
		//value interface{} // 存储 value
	}
	// SkipList 跳表定义
	SkipList struct {
		Node
		maxLevel       int         // 最⼤层数
		Len            int         // 跳表⻓度
		randSource     rand.Source // 随机数⽣成
		probability    float64
		probTable      []float64
		rWLock         *sync.RWMutex
		prevNodesCache []*Node
	}
)

// NewSkipList initializes and returns a new SkipList instance
func NewSkipList(maxLevel int, probability float64) *SkipList {
	// Initialize random source
	randSource := rand.NewSource(time.Now().UnixNano())

	// Initialize probability table based on maxLevel and probability
	probTable := make([]float64, maxLevel)
	for i := range probTable {
		probTable[i] = probability
	}

	// Initialize prevNodesCache
	prevNodesCache := make([]*Node, maxLevel)
	rWLock := &sync.RWMutex{}
	// Initialize SkipList
	return &SkipList{
		Node:           Node{next: make([]*Element, maxLevel)},
		maxLevel:       maxLevel,
		Len:            0,
		randSource:     randSource,
		probability:    probability,
		probTable:      probTable,
		prevNodesCache: prevNodesCache,
		rWLock:         rWLock,
	}
}

// Put 存储⼀个元素⾄跳表中，如果key已经存在，则会更新其对应的value
// 因此此跳表的实现暂不⽀持相同的key
func (t *SkipList) Set(key string, value []byte) *Element {
	var element *Element
	prev := t.backNodes(key) // 找出key节点在每⼀层索引应该放的位置的前⼀个节点
	if element = prev[0].next[0]; element != nil && element.KV.Key <= key {
		element.KV.Value = value // 如果key和prev的下⼀个节点的key相等，说明该key已存在，更新value返回 即可
		return element
	}
	element = &Element{
		Node: Node{
			next: make([]*Element, t.randomLevel()), // 初始化ele的next索引层
		},
		KV: kv.Value{
			Key:   key,
			Value: value,
		},
	}
	// 当前key应该插⼊的位置已经确定，就在prev的下⼀个位置
	for i := range element.next { // 遍历ele的所有索引层，建⽴节点前后联系
		element.next[i] = prev[i].next[i]
		prev[i].next[i] = element
	}
	t.Len++
	return element
}

// 找到key对应的前⼀个节点索引的信息，即key节点在每⼀层索引的前⼀个节点
func (t *SkipList) backNodes(key string) []*Node {
	var prev = &t.Node
	var next *Element
	prevs := t.prevNodesCache
	for i := t.maxLevel - 1; i >= 0; i-- { // 从最⾼层索引开始遍历
		next = prev.next[i]                    // 当前节点在第i层索引上的下⼀个节点
		for next != nil && key > next.KV.Key { // 如果⽬标节点的key⽐next节点的key⼤
			prev = &next.Node   // 将prev放到next节点的位置上
			next = next.next[i] // next通过当前层的索引跳到下⼀个位置
		} // 循环跳出后，key节点应位于pre和next之间
		prevs[i] = prev // 将当前的prev节点缓存到跳表中的对应层上
	} // 到下⼀层继续寻找
	return prevs
}

// Get 根据 key 查找对应的 Element 元素
//未找到则返回nil
func (t *SkipList) Search(key string) *Element {
	var prev = &t.Node
	var next *Element
	for i := t.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && key > next.KV.Key {
			prev = &next.Node
			next = next.next[i]
		}
	}
	if next != nil && next.KV.Key <= key {
		return next
	}
	return nil
}

// ⽣成索引随机层数
func (t *SkipList) randomLevel() (level int) {
	r := float64(t.randSource.Int63()) / (1 << 63) // ⽣成⼀个[0, 1)的概率值
	level = 1
	for level < t.maxLevel && r < t.probTable[level] { // 找到第⼀个prob⼩于 r 的层数
		level++
	}
	return
}
func probabilityTable(probability float64, maxLevel int) (table []float64) {
	for i := 1; i <= maxLevel; i++ { // 将每⼀层的prob值设置为p的(层数减⼀)次⽅
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

// Delete removes the element with the specified key from the skip list.
func (t *SkipList) Delete(key string) bool {
	current := &t.Node
	update := make([]*Node, t.maxLevel)
	// Locate the node to be deleted and record the path.
	for i := t.maxLevel - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].KV.Key < key {
			current = &current.next[i].Node
		}
		update[i] = current
	}
	// The node to delete is the next node at level 0.
	target := current.next[0]
	if target == nil || target.KV.Key != key {
		// Key not found.
		return false
	}
	// Update the next pointers to remove the target node.
	for i := 0; i < len(target.next); i++ {
		if update[i].next[i] == target {
			update[i].next[i] = target.next[i]
		}
	}
	t.Len--
	return true
}

func (t *SkipList) GetCount() int {
	return t.Len
}

func (t *SkipList) GetValues() []kv.Value {
	var values []kv.Value

	// Start at the first element in the bottom layer (level 0).
	current := t.Node.next[0]

	// Traverse the bottom layer and collect all values.
	for current != nil {
		values = append(values, current.KV)
		current = current.next[0]
	}

	return values
}

func (t *SkipList) Swap() *SkipList {
	t.rWLock.Lock()
	defer t.rWLock.Unlock()
	maxLevel := t.maxLevel
	probability := t.probability
	return NewSkipList(maxLevel, probability)

}
