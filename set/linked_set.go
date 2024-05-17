package set

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopi-frame/contract/support"
	"github.com/gopi-frame/support/lists"
)

var _ support.Set[support.Comparable] = (*LinkedSet[support.Comparable])(nil)

// NewLinkedSet creates a new linked hash set
func NewLinkedSet[E support.Comparable](values ...E) *LinkedSet[E] {
	set := new(LinkedSet[E])
	set.items = lists.NewLinkedList[E]()
	set.Push(values...)
	return set
}

// LinkedSet linked hash set
type LinkedSet[E support.Comparable] struct {
	items *lists.LinkedList[E]
}

func (s *LinkedSet[E]) Count() int64 {
	return s.items.Count()
}

func (s *LinkedSet[E]) IsEmpty() bool {
	return s.Count() == 0
}

func (s *LinkedSet[E]) IsNotEmpty() bool {
	return !s.IsEmpty()
}

func (s *LinkedSet[E]) Contains(value E) bool {
	return s.ContainsWhere(func(e E) bool {
		return e.Equals(value)
	})
}

func (s *LinkedSet[E]) ContainsWhere(callback func(E) bool) bool {
	var result bool
	s.items.Each(func(index int, value E) bool {
		if callback(value) {
			result = true
		}
		return !result
	})
	return result
}

func (s *LinkedSet[E]) Push(values ...E) {
	for _, value := range values {
		if s.Contains(value) {
			continue
		}
		s.items.Push(value)
	}
}

func (s *LinkedSet[E]) Remove(value E) {
	s.RemoveWhere(func(e E) bool {
		return e.Equals(value)
	})
}

func (s *LinkedSet[E]) RemoveWhere(callback func(E) bool) {
	items := s.items.Where(func(item E) bool {
		return !callback(item)
	})
	s.items = items
}

func (s *LinkedSet[E]) Clear() {
	s.items.Clear()
}

func (s *LinkedSet[E]) Each(callback func(int, E) bool) {
	s.items.Each(callback)
}

func (s *LinkedSet[E]) Clone() *LinkedSet[E] {
	return NewLinkedSet(s.ToArray()...)
}

func (s *LinkedSet[E]) ToArray() []E {
	return s.items.ToArray()
}

func (s *LinkedSet[E]) ToJSON() ([]byte, error) {
	return json.Marshal(s.ToArray())
}

func (s *LinkedSet[E]) MarshalJSON() ([]byte, error) {
	return s.ToJSON()
}

func (s *LinkedSet[E]) UnmarshalJSON(data []byte) error {
	items := lists.NewLinkedList[E]()
	err := json.Unmarshal(data, items)
	if err != nil {
		return err
	}
	s.items = items
	return nil
}

func (s *LinkedSet[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("LinkedSet[%T](len=%d)", *new(E), s.Count()))
	str.WriteByte('{')
	str.WriteByte('\n')
	s.items.Each(func(index int, value E) bool {
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
	if s.Count() > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
