package skiplist

import (
	"fmt"
	"math/rand"
	"time"
)

// 每个数据节点的信息
type nodeStructure struct {
	key     int
	value   int
	forward []*nodeStructure
}

// 每一个层的信息
type listStructure struct {
	level  int
	header *nodeStructure
}

const (
	MAX_NUMBER_OF_LEVELS = 11
	MAX_LEVEL            = 10
)

func newList() *listStructure {
	l := &listStructure{}

	l.level = 0
	l.header = newNodeOfLevel(MAX_NUMBER_OF_LEVELS)

	var i int
	for i = 0; i < MAX_NUMBER_OF_LEVELS; i++ {
		l.header.forward[i] = nil
	}
	return l
}

func newNodeOfLevel(level int) *nodeStructure {
	nodeArr := make([]*nodeStructure, level)
	return &nodeStructure{forward: nodeArr}
}

func randomLevel() int {
	seed := rand.New(rand.NewSource(time.Now().Unix()))
	return seed.Intn(MAX_LEVEL)
}

func insert(l *listStructure, key, value int) bool {
	var k int
	var update [MAX_NUMBER_OF_LEVELS]*nodeStructure

	var p, q *nodeStructure

	p = l.header
	k = l.level

	//这来采取 k--是跳表从上往下寻找适合的插入位置
	for ; k >= 0; k-- {
		q = p.forward[k]

		// 查找合适的插入位置
		for q != nil && q.key < key {
			//跳表往前移动一步
			p = q
			//往前移动一步之后指向后面的元素
			q = p.forward[k]
		}
		//记录需要改变指向的节点
		update[k] = p
	}

	//这一步的意义是不插入重复的元素
	if q != nil && q.key == key {
		//更新同一位置的数据
		//也许这里存在一个bug，稍后测试看是否需要和前后的值比较，确保链表是有序的
		q.value = value
		return false
	}

	k = randomLevel()
	fmt.Println("random level: %v", k)

	if k > l.level {
		l.level++
		k = l.level

		//如果随机生成的层数超过了当前的高度，那么header就要保持更新
		update[k] = l.header
	}

	//这里的意义是保持往上增长一层，避免不必要的浪费
	q = newNodeOfLevel(k + 1)
	//写入需要插入的值
	q.key, q.value = key, value

	// 更新节点的指向
	// 从最底层开始更新
	for ; k >= 0; k-- {
		p = update[k]
		q.forward[k] = p.forward[k]
		p.forward[k] = q
	}
	//跳表的插入工作完成
	return true
}
