package lists

import (
	listlib "container/list"
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/exception"
)

var _ support.List[any] = (*LinkedList[any])(nil)

// NewLinkedList new linked list
func NewLinkedList[E any](values ...E) *LinkedList[E] {
	instance := new(LinkedList[E])
	instance.Push(values...)
	return instance
}

// LinkedList linked list
type LinkedList[E any] struct {
	sync.Mutex
	list *listlib.List
}

func (list *LinkedList[E]) init() {
	if list.list == nil {
		list.list = listlib.New()
	}
}

func (list *LinkedList[E]) Count() int64 {
	list.init()
	return int64(list.list.Len())
}

func (list *LinkedList[E]) IsEmpty() bool {
	list.init()
	return list.Count() == 0
}

func (list *LinkedList[E]) IsNotEmpty() bool {
	list.init()
	return !list.IsEmpty()
}

func (list *LinkedList[E]) Contains(value E) bool {
	list.init()
	return list.ContainsWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

func (list *LinkedList[E]) ContainsWhere(callback func(value E) bool) bool {
	list.init()
	for e := list.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return true
		}
	}
	return false
}

func (list *LinkedList[E]) Push(values ...E) {
	list.init()
	for _, value := range values {
		list.list.PushBack(value)
	}
}

func (list *LinkedList[E]) Remove(value E) {
	list.RemoveWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

func (list *LinkedList[E]) RemoveWhere(callback func(item E) bool) {
	list.init()
	var next *listlib.Element
	for e := list.list.Front(); e != nil; e = next {
		next = e.Next()
		if callback(e.Value.(E)) {
			list.list.Remove(e)
		}
	}
}

func (list *LinkedList[E]) RemoveAt(index int) {
	list.init()
	var next *listlib.Element
	for e, i := list.list.Front(), 0; e != nil; e, i = next, i+1 {
		next = e.Next()
		if i == index {
			list.list.Remove(e)
			break
		}
	}
}

func (list *LinkedList[E]) Clear() {
	list.init()
	list.list.Init()
}

func (list *LinkedList[E]) Get(index int) E {
	list.init()
	if index < 0 || index >= list.list.Len() {
		panic(exception.NewRangeException(0, list.list.Len()-1))
	}
	for i, e := 0, list.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i == index {
			return e.Value.(E)
		}
	}
	return *new(E)
}

func (list *LinkedList[E]) Set(index int, value E) {
	list.init()
	for i, e := 0, list.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i == index {
			e.Value = value
		}
	}
}

func (list *LinkedList[E]) First() (E, bool) {
	list.init()
	if list.list.Len() == 0 {
		return *new(E), false
	}
	return list.list.Front().Value.(E), true
}

func (list *LinkedList[E]) FirstOr(value E) E {
	list.init()
	if list.list.Len() == 0 {
		return value
	}
	return list.list.Front().Value.(E)
}

func (list *LinkedList[E]) FirstWhere(callback func(item E) bool) (E, bool) {
	list.init()
	for e := list.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return e.Value.(E), true
		}
	}
	return *new(E), false
}

func (list *LinkedList[E]) FirstWhereOr(callback func(item E) bool, value E) E {
	list.init()
	for e := list.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			return e.Value.(E)
		}
	}
	return value
}

func (list *LinkedList[E]) Last() (E, bool) {
	list.init()
	if list.list.Len() == 0 {
		return *new(E), false
	}
	return list.list.Back().Value.(E), true
}

func (list *LinkedList[E]) LastOr(value E) E {
	list.init()
	if list.list.Back() == nil {
		return value
	}
	return list.list.Back().Value.(E)
}

func (list *LinkedList[E]) LastWhere(callback func(item E) bool) (E, bool) {
	list.init()
	for e := list.list.Back(); e != nil; e = e.Prev() {
		if callback(e.Value.(E)) {
			return e.Value.(E), true
		}
	}
	return *new(E), false
}

func (list *LinkedList[E]) LastWhereOr(callback func(item E) bool, value E) E {
	list.init()
	if v, ok := list.LastWhere(callback); ok {
		return v
	}
	return value
}

func (list *LinkedList[E]) Pop() (E, bool) {
	list.init()
	if list.list.Len() == 0 {
		return *new(E), false
	}
	item := list.list.Back()
	list.list.Remove(item)
	return item.Value.(E), true
}

func (list *LinkedList[E]) Shift() (E, bool) {
	list.init()
	if list.list.Len() == 0 {
		return *new(E), false
	}
	item := list.list.Front()
	list.list.Remove(item)
	return item.Value.(E), true
}

func (list *LinkedList[E]) Unshift(values ...E) {
	list.init()
	for _, value := range values {
		list.list.PushFront(value)
	}
}

func (list *LinkedList[E]) IndexOf(value E) int {
	list.init()
	return list.IndexOfWhere(func(item E) bool {
		return reflect.DeepEqual(item, value)
	})
}

func (list *LinkedList[E]) IndexOfWhere(callback func(item E) bool) int {
	list.init()
	for i, e := 0, list.list.Front(); e != nil; i, e = i+1, e.Next() {
		if callback(e.Value.(E)) {
			return i
		}
	}
	return -1
}

func (list *LinkedList[E]) Sub(from, to int) *LinkedList[E] {
	list.init()
	linked := NewLinkedList[E]()
	for i, e := 0, list.list.Front(); e != nil; i, e = i+1, e.Next() {
		if i < from {
			continue
		} else if i >= from && i < to {
			linked.Push(e.Value.(E))
		} else {
			break
		}
	}
	return linked
}

func (list *LinkedList[E]) Where(callback func(item E) bool) *LinkedList[E] {
	list.init()
	linked := &LinkedList[E]{}
	for e := list.list.Front(); e != nil; e = e.Next() {
		if callback(e.Value.(E)) {
			linked.Push(e.Value.(E))
		}
	}
	return linked
}

func (list *LinkedList[E]) Compact(callback func(a, b E) bool) {
	list.init()
	if list.list.Len() < 2 {
		return
	}
	if callback == nil {
		callback = func(a, b E) bool {
			return reflect.DeepEqual(a, b)
		}
	}
	var next *listlib.Element
	for e := list.list.Front().Next(); e != nil; e = next {
		next = e.Next()
		if callback(e.Value.(E), e.Prev().Value.(E)) {
			list.list.Remove(e)
		}
	}
}

func (list *LinkedList[E]) Min(callback func(a, b E) int) E {
	list.init()
	return slices.MinFunc(list.ToArray(), callback)
}

func (list *LinkedList[E]) Max(callback func(a, b E) int) E {
	list.init()
	return slices.MaxFunc(list.ToArray(), callback)
}

func (list *LinkedList[E]) Sort(callback func(a, b E) int) {
	list.init()
	var newList = listlib.New()
	for e := list.list.Front(); e != nil; e = e.Next() {
		node := newList.Front()
		for node != nil {
			if callback(e.Value.(E), node.Value.(E)) < 0 {
				newList.InsertBefore(e.Value, node)
				break
			}
			node = node.Next()
		}
		newList.PushBack(e.Value)
	}
	list.list = newList
}

func (list *LinkedList[E]) Chunk(size int) *LinkedList[*LinkedList[any]] {
	list.init()
	chunks := NewLinkedList[*LinkedList[any]]()
	chunk := NewLinkedList[any]()
	for e := list.list.Front(); e != nil; e = e.Next() {
		if chunk.list.Len() < size {
			chunk.Push(e.Value.(E))
		} else {
			chunks.Push(chunk)
			chunk = NewLinkedList(e.Value)
		}
	}
	chunks.Push(chunk)
	return chunks
}

func (list *LinkedList[E]) Each(callback func(index int, value E) bool) {
	list.init()
	for e, i := list.list.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		if !callback(i, e.Value.(E)) {
			break
		}
	}
}

func (list *LinkedList[E]) Reverse() {
	list.init()
	var next *listlib.Element
	for e := list.list.Front(); e != nil; e = next {
		next = e.Next()
		list.list.PushFront(e.Value)
		list.list.Remove(e)
	}
}

func (list *LinkedList[E]) Clone() *LinkedList[E] {
	list.init()
	linked := &LinkedList[E]{}
	for e := list.list.Front(); e != nil; e = e.Next() {
		linked.Push(e.Value.(E))
	}
	return linked
}

func (list *LinkedList[E]) String() string {
	list.init()
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedList[%T](len=%d)", *new(E), list.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	list.Each(func(index int, value E) bool {
		str.WriteByte('\t')
		if v, ok := any(value).(support.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", value))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		return index < 4
	})
	if list.list.Len() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}

func (list *LinkedList[E]) ToJSON() ([]byte, error) {
	list.init()
	return json.Marshal(list.ToArray())
}

func (list *LinkedList[E]) ToArray() []E {
	list.init()
	var items []E
	for e := list.list.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(E))
	}
	return items
}

func (list *LinkedList[E]) MarshalJSON() ([]byte, error) {
	list.init()
	return list.ToJSON()
}

func (list *LinkedList[E]) UnmarshalJSON(data []byte) error {
	list.init()
	items := []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		list.list.PushBack(item)
	}
	return nil
}
