package store

import (
	"fmt"
	"strings"
	"sync"
)

type Set interface {
	// 添加元素
	Add(items ...interface{})
	// 包含
	Contain(item interface{}) bool
	// 删除
	Remove(item interface{})
	// 集合大小
	Size() int
	// 清空
	Clear()
	// 判断是否为空
	Empty() bool
	// 创建副本
	Duplicate() Set
	// 数组
	ToSlice() []interface{}
}

// threadUnsafeSet 非线程安全的集合，采用元素值作为key，空的struct作为值
type threadUnsafeSet map[interface{}]struct{}

func newThreadUnsafeSet() threadUnsafeSet {
	return make(threadUnsafeSet)
}

func (set *threadUnsafeSet) Add(items ...interface{}) {
	for _, item := range items {
		(*set)[item] = struct{}{}
	}
}

func (set *threadUnsafeSet) Contain(item interface{}) bool {
	_, ok := (*set)[item]
	return ok
}

func (set *threadUnsafeSet) Remove(item interface{}) {
	delete((*set), item)
}

func (set *threadUnsafeSet) Size() int {
	return len((*set))
}

func (set *threadUnsafeSet) Clear() {
	*set = newThreadUnsafeSet()
}

func (set *threadUnsafeSet) Empty() bool {
	return set.Size() == 0
}

func (set *threadUnsafeSet) Duplicate() Set {
	duplicateSet := newThreadUnsafeSet()
	for item, _ := range *set {
		duplicateSet.Add(item)
	}
	return &duplicateSet
}

func (set *threadUnsafeSet) String() string {
	items := make([]string, 0, len(*set))

	for elem := range *set {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ", "))
}

func (set *threadUnsafeSet) ToSlice() []interface{} {
	keys := make([]interface{}, 0, set.Size())
	for elem := range *set {
		keys = append(keys, elem)
	}

	return keys
}

type threadSafeSet struct {
	s threadUnsafeSet
	sync.RWMutex
}

func newThreadSafeSet() *threadSafeSet  {
	return &threadSafeSet{s: newThreadUnsafeSet()}
}

func (set *threadSafeSet) Add(items ...interface{}) {
	set.Lock() //数据新增采用互斥锁
	set.s.Add(items...)
	set.Unlock()
}

func (set *threadSafeSet) Contain(item interface{}) bool {
	set.RLock() // 采用读写锁
	defer set.RUnlock()
	return set.s.Contain(item)
}

func (set *threadSafeSet) Remove(item interface{}) {
	set.Lock()
	set.s.Remove(item)
	set.Unlock()
}

func (set *threadSafeSet) Size() int {
	set.RLock()
	defer set.RUnlock()
	return set.s.Size()
}

func (set *threadSafeSet) Clear() {
	set.Lock()
	set.s.Clear()
	set.Unlock()
}

func (set *threadSafeSet) Empty() bool {
	return set.Size() == 0
}

func (set *threadSafeSet) Duplicate() Set {
	set.RLock()
	defer set.RUnlock()
	s := set.s.Duplicate()
	return &threadSafeSet{s: *(s.(*threadUnsafeSet))}
}

func (set *threadSafeSet) ToSlice() []interface{} {
	set.RLock()
	defer set.RUnlock()
	return set.s.ToSlice()
}

func ThreadSafe(items ...interface{}) Set {
	s := newThreadSafeSet()
	s.Add(items...)
	return s
}
