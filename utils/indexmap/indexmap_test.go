package indexmap_test

import (
	"testing"

	"github.com/howjmay/sui-go/sui_types"
	"github.com/howjmay/sui-go/utils/indexmap"
	"github.com/stretchr/testify/require"
)

func TestIndexMap(t *testing.T) {
	t.Run("primitive type", func(t *testing.T) {
		m := indexmap.NewIndexMap[string, int]()
		m.Insert("first", 1)
		m.Insert("second", 2)

		val, ok := m.Get("first")
		require.True(t, ok)
		require.Equal(t, 1, val)
		idx, ok := m.Find("first")
		require.True(t, ok)
		require.Equal(t, idx, 0)

		m.Insert("first", 3)
		val, ok = m.Get("first")
		require.True(t, ok)
		require.Equal(t, 3, val)
		idx, ok = m.Find("first")
		require.True(t, ok)
		require.Equal(t, idx, 0)

		var targetList []int = []int{6, 4}
		var testList []int
		m.ForEach(func(k string, v int) {
			testList = append(testList, v*2)
		})
		require.Equal(t, targetList, testList)
	})

	t.Run("customized type", func(t *testing.T) {
		m := indexmap.NewIndexMap[sui_types.BuilderArg, sui_types.CallArg]()
		testBytes := [][]byte{
			[]byte{1, 4, 7},
			[]byte{2, 5, 8},
			[]byte{3, 6, 9},
			[]byte{10, 11, 12},
			[]byte{13, 14, 15},
		}
		m.Insert(sui_types.BuilderArg{Pure: &testBytes[0]}, sui_types.CallArg{Pure: &testBytes[0]})
		m.Insert(sui_types.BuilderArg{Pure: &testBytes[1]}, sui_types.CallArg{Pure: &testBytes[1]})

		val, ok := m.Get(sui_types.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, sui_types.CallArg{Pure: &testBytes[0]}, val)
		idx, ok := m.Find(sui_types.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, idx, 0)

		m.Insert(sui_types.BuilderArg{Pure: &testBytes[0]}, sui_types.CallArg{Pure: &testBytes[2]})
		val, ok = m.Get(sui_types.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, sui_types.CallArg{Pure: &testBytes[2]}, val)
		idx, ok = m.Find(sui_types.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, idx, 0)

		var targetList []sui_types.CallArg = []sui_types.CallArg{
			sui_types.CallArg{Pure: &testBytes[3]},
			sui_types.CallArg{Pure: &testBytes[4]},
		}
		var testList []sui_types.CallArg
		i := 3
		m.ForEach(func(k sui_types.BuilderArg, v sui_types.CallArg) {
			testList = append(testList, sui_types.CallArg{Pure: &testBytes[i]})
			i++
		})
		require.Equal(t, targetList, testList)
	})
}
