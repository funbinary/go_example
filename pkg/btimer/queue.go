package btimer

import (
	"container/heap"
	"go.uber.org/atomic"
	"math"
	"sync"
)

type priorityQueue struct {
	sync.RWMutex
	heap         *priorityQueueHeap
	nextPriority *atomic.Int64
}

func newPriorityQueue() *priorityQueue {
	queue := &priorityQueue{
		heap:         &priorityQueueHeap{array: make([]item, 0)},
		nextPriority: atomic.NewInt64(math.MaxInt64),
	}
	heap.Init(queue.heap)
	return queue
}

func (self *priorityQueue) NextPriority() int64 {
	return self.nextPriority.Load()
}

func (self *priorityQueue) Push(value interface{}, priority int64) {
	self.Lock()
	defer self.Unlock()
	heap.Push(self.heap, item{
		value:    value,
		priority: priority,
	})
	// Update the minimum priority using atomic operation.
	nextPriority := self.nextPriority.Load()
	if priority >= nextPriority {
		return
	}
	self.nextPriority.Store(priority)
}

func (self *priorityQueue) Pop() interface{} {
	self.Lock()
	defer self.Unlock()
	if v := heap.Pop(self.heap); v != nil {
		var nextPriority int64 = math.MaxInt64
		if len(self.heap.array) > 0 {
			nextPriority = self.heap.array[0].priority
		}
		self.nextPriority.Store(nextPriority)
		return v.(item).value
	}
	return nil
}

type priorityQueueHeap struct {
	array []item
}

// Len  sort.Interface.
func (self *priorityQueueHeap) Len() int {
	return len(self.array)
}

// Less sort.Interface.
func (self *priorityQueueHeap) Less(i, j int) bool {
	return self.array[i].priority < self.array[j].priority
}

// Swap sort.Interface.
func (self *priorityQueueHeap) Swap(i, j int) {
	if len(self.array) == 0 {
		return
	}
	self.array[i], self.array[j] = self.array[j], self.array[i]
}

// Push pushes an item to the heap.
func (self *priorityQueueHeap) Push(x interface{}) {
	self.array = append(self.array, x.(item))
}

// Pop retrieves, removes and returns the most high priority item from the heap.
func (self *priorityQueueHeap) Pop() interface{} {
	length := len(self.array)
	if length == 0 {
		return nil
	}
	item := self.array[length-1]
	self.array = self.array[0 : length-1]
	return item
}

type item struct {
	value    interface{}
	priority int64
}
