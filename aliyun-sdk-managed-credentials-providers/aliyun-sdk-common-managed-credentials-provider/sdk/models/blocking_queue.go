package models

import (
	"container/list"
	"errors"
	"sync"
)

type Collection interface {
	Size() uint64
	Capacity() uint64
	IsEmpty() bool
	Clear()

	Push(item interface{})
	Pop() interface{}
	Offer(item interface{}) bool

	Peek() interface{}
}

type BlockingQueue struct {
	lock     *sync.Mutex
	notEmpty *sync.Cond
	notFull  *sync.Cond
	count    uint64

	queue    *list.List
	capacity uint64
}

func NewBlockingQueue(capacity uint64) *BlockingQueue {
	if capacity < 1 {
		capacity = 1
	}
	lock := new(sync.Mutex)
	return &BlockingQueue{
		lock:     lock,
		notEmpty: sync.NewCond(lock),
		notFull:  sync.NewCond(lock),
		count:    uint64(0),
		queue:    list.New(),
		capacity: capacity,
	}
}

func (q *BlockingQueue) Size() uint64 {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

func (q *BlockingQueue) Capacity() uint64 {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.capacity
}

func (q *BlockingQueue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	for q.count == 0 {
		q.notEmpty.Wait()
	}
	res, _ := q.tryPop()
	return res
}

func (q *BlockingQueue) Push(value interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	for q.count == q.capacity {
		q.notFull.Wait()
	}
	_, _ = q.tryPush(value)
	return
}

func (q *BlockingQueue) Offer(value interface{}) bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	ok, _ := q.tryPush(value)
	return ok
}

func (q *BlockingQueue) tryPush(value interface{}) (bool, error) {
	if q.count == q.capacity {
		return false, errors.New("queue is full")
	}
	q.push(value)
	return true, nil
}

func (q *BlockingQueue) tryPop() (interface{}, error) {
	if q.count == 0 {
		return nil, errors.New("queue is empty")
	}
	return q.pop(), nil
}

func (q *BlockingQueue) push(value interface{}) {
	q.queue.PushBack(value)
	q.count += 1
	q.notEmpty.Signal()
}

func (q *BlockingQueue) pop() interface{} {
	value := q.queue.Remove(q.queue.Front())
	q.count -= 1
	q.notFull.Signal()
	return value
}

func (q *BlockingQueue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.count == 0 {
		return nil
	}
	return q.queue.Front()
}

func (q *BlockingQueue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *BlockingQueue) Clear() {
	q.lock.Lock()
	defer q.lock.Unlock()
	for e := q.queue.Front(); e != nil; e = e.Next() {
		q.queue.Remove(e)
	}
	q.count = uint64(0)
	q.notFull.Broadcast()
}
