package set

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/gopi-frame/types"
	"github.com/stretchr/testify/assert"
)

func TestSet_Count(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.Equal(t, int64(3), set.Count())
}

func TestSet_IsEmpty(t *testing.T) {
	set := NewSet[types.Int]()
	assert.True(t, set.IsEmpty())
}

func TestSet_IsNotEmpty(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.IsNotEmpty())
}

func TestSet_Contains(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.Contains(1))
}

func TestSet_ContainsWhere(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.ContainsWhere(func(i types.Int) bool {
		return i == 2
	}))
}

func TestSet_Remove(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.Contains(1))
	set.Remove(1)
	assert.False(t, set.Contains(1))
}

func TestSet_RemoveWhere(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.ContainsWhere(func(i types.Int) bool {
		return i == 2
	}))
	set.RemoveWhere(func(i types.Int) bool {
		return i == 2
	})
	assert.False(t, set.ContainsWhere(func(i types.Int) bool {
		return i == 2
	}))
}

func TestSet_Each(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	items := []int{}
	set.Each(func(_ int, item types.Int) bool {
		items = append(items, int(item))
		return true
	})
	assert.Equal(t, []int{1, 2, 3}, items)
}

func TestSet_Cleaar(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.True(t, set.IsNotEmpty())
	set.Clear()
	assert.True(t, set.IsEmpty())
}

func TestSet_Clone(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	set2 := set.Clone()
	assert.Equal(t, set.items, set2.items)
}

func TestSet_ToArray(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	assert.Equal(t, []types.Int{1, 2, 3}, set.ToArray())
}

func TestSet_ToJSON(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	jsonBytes, err := set.ToJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestSet_MarshalJSON(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	jsonBytes, err := set.MarshalJSON()
	assert.Nil(t, err)
	assert.JSONEq(t, `[1,2,3]`, string(jsonBytes))
}

func TestSet_UnmarshalJSON(t *testing.T) {
	set := NewSet[types.Int]()
	err := json.Unmarshal([]byte(`[1,2,3]`), set)
	assert.Nil(t, err)
	assert.Equal(t, []types.Int{1, 2, 3}, set.ToArray())
}

func TestSet_String(t *testing.T) {
	set := NewSet[types.Int](1, 2, 3)
	str := set.String()
	pattern := regexp.MustCompile(fmt.Sprintf(`Set\[types\.Int\]\(len=%d\)\{\n(\t\d+,\n){3}\}`, set.size))
	assert.True(t, pattern.MatchString(str))
}
