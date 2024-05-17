package queue

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gopi-frame/contract/support"
)

var _ support.BlockingQueue[any] = (*PriorityBlockingQueue[any])(nil)

// NewPriorityBlockingQueue new priority blocking queue
func NewPriorityBlockingQueue[E any](comparator support.Comparator[E], cap int) *PriorityBlockingQueue[E] {
	queue := new(PriorityBlockingQueue[E])
	queue.items = NewPriorityQueue(comparator)
	queue.takeLock = sync.NewCond(queue.items)
	queue.putLock = sync.NewCond(queue.items)
	queue.cap = cap
	return queue
}

// PriorityBlockingQueue priority blocking queue
type PriorityBlockingQueue[E any] struct {
	items    *PriorityQueue[E]
	cap      int
	takeLock *sync.Cond
	putLock  *sync.Cond
}

func (q *PriorityBlockingQueue[E]) Count() int64 {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.Count()
}

func (q *PriorityBlockingQueue[E]) IsEmpty() bool {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.IsEmpty()
}

func (q *PriorityBlockingQueue[E]) IsNotEmpty() bool {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.IsNotEmpty()
}

func (q *PriorityBlockingQueue[E]) Clear() {
	q.items.Lock()
	defer q.items.Unlock()
	q.items.Clear()
}

func (q *PriorityBlockingQueue[E]) Peek() (E, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.Peek()
}

func (q *PriorityBlockingQueue[E]) TryEnqueue(value E) bool {
	q.items.Lock()
	defer q.items.Unlock()
	if q.cap == int(q.items.Count()) {
		return false
	}
	ok := q.items.Enqueue(value)
	q.takeLock.Broadcast()
	return ok
}

func (q *PriorityBlockingQueue[E]) TryDequeue() (E, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	if q.items.Count() == 0 {
		return *new(E), false
	}
	value, ok := q.items.Dequeue()
	q.putLock.Broadcast()
	return value, ok
}

func (q *PriorityBlockingQueue[E]) Enqueue(value E) bool {
	q.items.Lock()
	defer q.items.Unlock()
	for q.cap == int(q.items.Count()) {
		q.putLock.Wait()
	}
	ok := q.items.Enqueue(value)
	q.takeLock.Broadcast()
	return ok
}

func (q *PriorityBlockingQueue[E]) Dequeue() (E, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	for q.items.IsEmpty() {
		q.takeLock.Wait()
	}
	value, ok := q.items.Dequeue()
	q.putLock.Broadcast()
	return value, ok
}

func (q *PriorityBlockingQueue[E]) EnqueueTimeout(value E, duration time.Duration) bool {
	timeout := time.After(duration)
	done := make(chan struct{})
	go func() {
		q.items.Lock()
		defer q.items.Unlock()
		for int64(q.cap) == q.items.Count() {
			q.putLock.Wait()
		}
		close(done)
	}()
	select {
	case <-done:
		ok := q.items.Enqueue(value)
		q.takeLock.Broadcast()
		return ok
	case <-timeout:
		return false
	}
}

func (q *PriorityBlockingQueue[E]) DequeueTimeout(duration time.Duration) (value E, ok bool) {
	timeout := time.After(duration)
	done := make(chan struct{})
	go func() {
		q.items.Lock()
		defer q.items.Unlock()
		for q.items.IsEmpty() {
			q.takeLock.Wait()
		}
		close(done)
	}()
	select {
	case <-done:
		value, ok = q.items.Dequeue()
		q.putLock.Broadcast()
		return
	case <-timeout:
		return
	}
}

func (q *PriorityBlockingQueue[E]) ToArray() []E {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.ToArray()
}

func (q *PriorityBlockingQueue[E]) ToJSON() ([]byte, error) {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.ToJSON()
}

func (q *PriorityBlockingQueue[E]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

func (q *PriorityBlockingQueue[E]) UnmarshalJSON(data []byte) error {
	q.items.Lock()
	defer q.items.Unlock()
	values := make([]E, 0)
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	q.items.Clear()
	for _, value := range values {
		for q.cap == int(q.items.Count()) {
			q.putLock.Wait()
		}
		q.items.Enqueue(value)
		q.takeLock.Broadcast()
	}
	return nil
}

func (q *PriorityBlockingQueue[E]) String() string {
	q.items.Lock()
	defer q.items.Unlock()
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("PriorityBlockingQueue[%T](len=%d)", *new(E), q.items.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, value := range q.items.items {
		str.WriteByte('\t')
		if v, ok := any(value).(support.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		if index >= 4 {
			break
		}
	}
	if q.items.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
