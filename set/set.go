package set

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gopi-frame/contract/support"
)

var _ support.Set[support.Comparable] = (*Set[support.Comparable])(nil)

// NewSet new set
func NewSet[E support.Comparable](values ...E) *Set[E] {
	set := new(Set[E])
	set.Push(values...)
	return set
}

// Set hash set
type Set[E support.Comparable] struct {
	sync.RWMutex
	items []E
	size  int64
}

// Count count
func (s *Set[E]) Count() int64 {
	return s.size
}

// IsEmpty is empty
func (s *Set[E]) IsEmpty() bool {
	return s.Count() == 0
}

// IsNotEmpty is not empty
func (s *Set[E]) IsNotEmpty() bool {
	return !s.IsEmpty()
}

// Contains contains
func (s *Set[E]) Contains(value E) bool {
	return s.ContainsWhere(func(e E) bool {
		return e.Equals(value)
	})
}

// ContainsWhere comtains where
func (s *Set[E]) ContainsWhere(callback func(E) bool) bool {
	for _, item := range s.items {
		if callback(item) {
			return true
		}
	}
	return false
}

// Push push
func (s *Set[E]) Push(values ...E) {
	for _, value := range values {
		if s.Contains(value) {
			continue
		}
		s.items = append(s.items, value)
		s.size++
	}
}

// Remove remove
func (s *Set[E]) Remove(value E) {
	s.RemoveWhere(func(e E) bool {
		return e.Equals(value)
	})
}

// RemoveWhere remove where
func (s *Set[E]) RemoveWhere(callback func(E) bool) {
	items := []E{}
	size := int64(0)
	for _, item := range s.items {
		if callback(item) {
			continue
		}
		items = append(items, item)
		size++
	}
	s.items = items
	s.size = size
}

// Each each
func (s *Set[E]) Each(callback func(_ int, item E) bool) {
	for index, item := range s.items {
		if !callback(index, item) {
			break
		}
	}
}

// Clear clear
func (s *Set[E]) Clear() {
	s.items = make([]E, 0)
	s.size = 0
}

// Clone clone
func (s *Set[E]) Clone() *Set[E] {
	return NewSet(s.items...)
}

// ToArray to array
func (s *Set[E]) ToArray() []E {
	return s.items
}

// ToJSON to json
func (s *Set[E]) ToJSON() ([]byte, error) {
	return json.Marshal(s.items)
}

func (s *Set[E]) MarshalJSON() ([]byte, error) {
	return s.ToJSON()
}

func (s *Set[E]) UnmarshalJSON(data []byte) error {
	var items = []E{}
	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	s.items = items
	s.size = int64(len(items))
	return nil
}

func (s *Set[E]) String() string {
	str := new(strings.Builder)
	str.WriteString(fmt.Sprintf("Set[%T](len=%d)", *new(E), s.size))
	str.WriteByte('{')
	str.WriteByte('\n')
	for index, item := range s.items {
		str.WriteByte('\t')
		if v, ok := any(item).(support.Stringable); ok {
			str.WriteString(v.String())
		} else {
			str.WriteString(fmt.Sprintf("%v", item))
		}
		str.WriteByte(',')
		str.WriteByte('\n')
		if index >= 4 {
			break
		}
	}
	if s.size > 5 {
		str.WriteString("\t...\n")
	}
	str.WriteByte('}')
	return str.String()
}
