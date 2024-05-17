package lists

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/gopi-frame/contract/support"
)

var _ support.List[any] = (*List[any])(nil)

// NewList new list
func NewList[E any](values ...E) *List[E] {
	instance := new(List[E])
	instance.Push(values...)
	return instance
}

// List list
type List[E any] struct {
	sync.Mutex
	items []E
}

func (list *List[E]) Count() int64 {
	return int64(len(list.items))
}

func (list *List[E]) IsEmpty() bool {
	return list.Count() == 0
}

func (list *List[E]) IsNotEmpty() bool {
	return !list.IsEmpty()
}

func (list *List[E]) Contains(value E) bool {
	return list.ContainsWhere(func(e E) bool {
		return reflect.DeepEqual(e, value)
	})
}

func (list *List[E]) ContainsWhere(callback func(value E) bool) bool {
	return slices.ContainsFunc(list.items, callback)
}

func (list *List[E]) Push(values ...E) {
	list.items = append(list.items, values...)
}

func (list *List[E]) Remove(value E) {
	list.RemoveWhere(func(item E) bool {
		return reflect.DeepEqual(value, item)
	})
}

func (list *List[E]) RemoveWhere(callback func(item E) bool) {
	list.items = slices.DeleteFunc(list.items, callback)
}

func (list *List[E]) RemoveAt(index int) {
	list.items = slices.Delete(list.items, index, index+1)
}

func (list *List[E]) Clear() {
	list.items = []E{}
}

func (list *List[E]) Get(index int) E {
	return list.items[index]
}

func (list *List[E]) Set(index int, value E) {
	list.items[index] = value
}

func (list *List[E]) First() (E, bool) {
	if len(list.items) == 0 {
		return *new(E), false
	}
	return list.items[0], true
}

func (list *List[E]) FirstOr(value E) E {
	if v, ok := list.First(); ok {
		return v
	}
	return value
}

func (list *List[E]) FirstWhere(callback func(item E) bool) (E, bool) {
	for _, item := range list.items {
		if callback(item) {
			return item, true
		}
	}
	return *new(E), false
}

func (list *List[E]) FirstWhereOr(callback func(item E) bool, value E) E {
	if v, found := list.FirstWhere(callback); found {
		return v
	}
	return value
}

func (list *List[E]) Last() (E, bool) {
	length := len(list.items)
	if length == 0 {
		return *new(E), false
	}
	return list.items[length-1], true
}

func (list *List[E]) LastOr(value E) E {
	if v, ok := list.Last(); ok {
		return v
	}
	return value
}

func (list *List[E]) LastWhere(callback func(item E) bool) (E, bool) {
	length := len(list.items)
	for index := range list.items {
		if value := list.items[length-index-1]; callback(value) {
			return value, true
		}
	}
	return *new(E), false
}

func (list *List[E]) LastWhereOr(callback func(item E) bool, value E) E {
	if v, ok := list.LastWhere(callback); ok {
		return v
	}
	return value
}

func (list *List[E]) Pop() (E, bool) {
	length := len(list.items)
	if length == 0 {
		return *new(E), false
	}
	value := list.items[length-1]
	list.items = list.items[:length-1]
	return value, true
}

func (list *List[E]) Shift() (E, bool) {
	if len(list.items) == 0 {
		return *new(E), false
	}
	value := list.items[0]
	list.items = list.items[1:]
	return value, true
}

func (list *List[E]) Unshift(values ...E) {
	list.items = slices.Insert(list.items, 0, values...)
}

func (list *List[E]) IndexOf(value E) int {
	return list.IndexOfWhere(func(item E) bool {
		return reflect.DeepEqual(value, item)
	})
}

func (list *List[E]) IndexOfWhere(callback func(item E) bool) int {
	return slices.IndexFunc(list.items, callback)
}

func (list *List[E]) Sub(from, to int) *List[E] {
	return &List[E]{items: list.items[from:to]}
}

func (list *List[E]) Where(callback func(item E) bool) *List[E] {
	l := &List[E]{}
	for _, item := range list.items {
		if callback(item) {
			l.items = append(l.items, item)
		}
	}
	return l
}

func (list *List[E]) Compact(callback func(a, b E) bool) {
	if callback == nil {
		callback = func(a, b E) bool {
			return reflect.DeepEqual(a, b)
		}
	}
	list.items = slices.CompactFunc(list.items, callback)
}

func (list *List[E]) Min(callback func(a, b E) int) E {
	return slices.MinFunc(list.items, callback)
}

func (list *List[E]) Max(callback func(a, b E) int) E {
	return slices.MaxFunc(list.items, callback)
}

func (list *List[E]) Sort(callback func(a, b E) int) {
	slices.SortFunc(list.items, callback)
}

func (list *List[E]) Chunk(size int) *List[*List[any]] {
	chunks := NewList[*List[any]]()
	chunk := NewList[any]()
	for _, item := range list.items {
		if len(chunk.items) < size {
			chunk.Push(item)
		} else {
			chunks.Push(chunk)
			chunk = NewList[any](item)
		}
	}
	chunks.Push(chunk)
	return chunks
}

func (list *List[E]) Each(callback func(index int, value E) bool) {
	for index, value := range list.items {
		if !callback(index, value) {
			break
		}
	}
}

func (list *List[E]) Reverse() {
	slices.Reverse(list.items)
}

func (list *List[E]) Clone() *List[E] {
	list.items = slices.Clone(list.items)
	return list
}

func (list *List[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("List[%T](len=%d)", *new(E), list.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, value := range list.items {
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
	if list.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}

func (list *List[E]) ToJSON() ([]byte, error) {
	return json.Marshal(list.items)
}

func (list *List[E]) ToArray() []E {
	return list.items
}

func (list *List[E]) MarshalJSON() ([]byte, error) {
	return list.ToJSON()
}

func (list *List[E]) UnmarshalJSON(data []byte) error {
	items := []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	list.items = items
	return nil
}
