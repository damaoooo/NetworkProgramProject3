package Utils

import (
	"NPProj3/ORM"
	"sync"
)

type ItemQueue struct {
	Items []ORM.Event
	lock  sync.RWMutex
}

//Make a New Queue
func (q *ItemQueue) New() *ItemQueue {
	q.Items = []ORM.Event{}
	return q
}

//Judge this queue is Empty or not
func (q *ItemQueue) IsEmpty() bool {
	return len(q.Items) == 0
}

//Add a new item in this queue
func (q *ItemQueue) Add(t ORM.Event) {
	q.lock.Lock()
	q.Items = append(q.Items, t)
	q.lock.Unlock()
}

//Pop the most front item of this queue and return it.
//If this queue is Empty, return nil and do nothing
func (q *ItemQueue) Dequeue() *ORM.Event {
	if !q.IsEmpty() {
		q.lock.Lock()
		item := q.Items[0]
		q.Items = q.Items[1:len(q.Items)]
		q.lock.Unlock()
		return &item
	} else {
		return nil
	}
}

//Return the Size of this Queue
func (q *ItemQueue) Size() int {
	return len(q.Items)
}

//Get a New Queue
func InitQueue() *ItemQueue {
	queue := ItemQueue{}
	queue.New()
	return &queue
}
