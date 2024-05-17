package queue

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gopi-frame/contract/support"
)

var _ support.BlockingQueue[support.Delayable[any]] = (*DelayedQueue[support.Delayable[any], any])(nil)

// NewDelayedQueue new delayed queue
func NewDelayedQueue[Q support.Delayable[T], T any]() *DelayedQueue[Q, T] {
	queue := new(DelayedQueue[Q, T])
	queue.items = NewPriorityQueue(queue)
	queue.takeLock = sync.NewCond(queue.items)
	return queue
}

// DelayedQueue delayed queue
type DelayedQueue[Q support.Delayable[T], T any] struct {
	items    *PriorityQueue[Q]
	takeLock *sync.Cond
}

func (q *DelayedQueue[Q, T]) Compare(a, b Q) int {
	if a.Until().Before(b.Until()) {
		return -1
	} else if a.Until().After(b.Until()) {
		return 1
	} else {
		return 0
	}
}

func (q *DelayedQueue[Q, T]) Count() int64 {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.Count()
}

func (q *DelayedQueue[Q, T]) IsEmpty() bool {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.IsEmpty()
}

func (q *DelayedQueue[Q, T]) IsNotEmpty() bool {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.IsNotEmpty()
}

func (q *DelayedQueue[Q, T]) Clear() {
	q.items.Lock()
	defer q.items.Unlock()
	q.items.Clear()
}

func (q *DelayedQueue[Q, T]) Peek() (Q, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.Peek()
}

func (q *DelayedQueue[Q, T]) TryEnqueue(value Q) bool {
	return q.Enqueue(value)
}

func (q *DelayedQueue[Q, T]) Enqueue(value Q) bool {
	q.items.Lock()
	defer q.items.Unlock()
	ok := q.items.Enqueue(value)
	q.takeLock.Broadcast()
	return ok
}

func (q *DelayedQueue[Q, T]) EnqueueTimeout(value Q, duration time.Duration) bool {
	return q.Enqueue(value)
}

func (q *DelayedQueue[Q, T]) TryDequeue() (Q, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	if v, ok := q.items.Peek(); ok && v.Until().Before(time.Now()) {
		return q.items.Dequeue()
	}
	return *new(Q), false
}

func (q *DelayedQueue[Q, T]) Dequeue() (Q, bool) {
	q.items.Lock()
	defer q.items.Unlock()
	for q.items.IsEmpty() {
		q.takeLock.Wait()
	}
	v, _ := q.items.Peek()
	timer := time.NewTimer(time.Until(v.Until()))
	defer timer.Stop()
	<-timer.C
	return q.items.Dequeue()
}

func (q *DelayedQueue[Q, T]) DequeueTimeout(duration time.Duration) (support.Delayable[T], bool) {
	timeout := time.After(duration)
	done := make(chan struct{})
	go func() {
		q.items.Lock()
		defer q.items.Unlock()
		for q.items.IsEmpty() {
			q.takeLock.Wait()
		}
		if v, ok := q.items.Peek(); ok {
			timer := time.NewTimer(time.Until(v.Until()))
			defer timer.Stop()
			<-timer.C
			close(done)
		}

	}()
	select {
	case <-timeout:
		return *new(support.Delayable[T]), false
	case <-done:
		return q.items.Dequeue()
	}
}

func (q *DelayedQueue[Q, T]) ToArray() []Q {
	q.items.Lock()
	defer q.items.Unlock()
	return q.items.ToArray()
}

func (q *DelayedQueue[Q, T]) ToJSON() ([]byte, error) {
	q.items.Lock()
	defer q.items.Unlock()
	return json.Marshal(q.items.ToArray())
}

func (q *DelayedQueue[Q, T]) MarshalJSON() ([]byte, error) {
	return q.ToJSON()
}

func (q *DelayedQueue[Q, T]) UnmarshalJSON(data []byte) error {
	q.items.Lock()
	defer q.items.Unlock()
	items := []Q{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		q.Enqueue(item)
	}
	return nil
}

func (q *DelayedQueue[Q, T]) String() string {
	q.items.Lock()
	defer q.items.Unlock()
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("DelayedQueue[%T](len=%d)", *new(T), q.items.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	items := q.ToArray()
	for index, item := range items {
		str.WriteByte('\t')
		if v, ok := any(item).(support.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("value: %v, unitl: %v", item.Value(), item.Until().Format("2006-01-02 15:04:05")))
		}
		if index >= 4 {
			break
		}
	}
	if q.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
