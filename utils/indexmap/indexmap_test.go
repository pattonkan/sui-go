package indexmap_test

import (
	"testing"

	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/utils/indexmap"
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
		m := indexmap.NewIndexMap[suiptb.BuilderArg, suiptb.CallArg]()
		testBytes := [][]byte{
			[]byte{1, 4, 7},
			[]byte{2, 5, 8},
			[]byte{3, 6, 9},
			[]byte{10, 11, 12},
			[]byte{13, 14, 15},
		}
		m.Insert(suiptb.BuilderArg{Pure: &testBytes[0]}, suiptb.CallArg{Pure: &testBytes[0]})
		m.Insert(suiptb.BuilderArg{Pure: &testBytes[1]}, suiptb.CallArg{Pure: &testBytes[1]})

		val, ok := m.Get(suiptb.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, suiptb.CallArg{Pure: &testBytes[0]}, val)
		idx, ok := m.Find(suiptb.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, idx, 0)

		m.Insert(suiptb.BuilderArg{Pure: &testBytes[0]}, suiptb.CallArg{Pure: &testBytes[2]})
		val, ok = m.Get(suiptb.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, suiptb.CallArg{Pure: &testBytes[2]}, val)
		idx, ok = m.Find(suiptb.BuilderArg{Pure: &testBytes[0]})
		require.True(t, ok)
		require.Equal(t, idx, 0)

		var targetList []suiptb.CallArg = []suiptb.CallArg{
			suiptb.CallArg{Pure: &testBytes[3]},
			suiptb.CallArg{Pure: &testBytes[4]},
		}
		var testList []suiptb.CallArg
		i := 3
		m.ForEach(func(k suiptb.BuilderArg, v suiptb.CallArg) {
			testList = append(testList, suiptb.CallArg{Pure: &testBytes[i]})
			i++
		})
		require.Equal(t, targetList, testList)
	})

	t.Run("InsertFull returns correct index", func(t *testing.T) {
		m := indexmap.NewIndexMap[string, int]()

		// First insertion should return index 0
		idx := m.InsertFull("key1", 100)
		require.Equal(t, 0, idx)

		// Second insertion should return index 1
		idx = m.InsertFull("key2", 200)
		require.Equal(t, 1, idx)

		// Update of existing key should return original index
		idx = m.InsertFull("key1", 300)
		require.Equal(t, 0, idx)

		// Verify the value was updated
		val, ok := m.Get("key1")
		require.True(t, ok)
		require.Equal(t, 300, val)
	})

	t.Run("InsertFull with equivalent structs", func(t *testing.T) {
		type TestStruct struct {
			Field1 string
			Field2 int
		}

		m := indexmap.NewIndexMap[TestStruct, string]()

		struct1 := TestStruct{"test", 123}
		struct2 := TestStruct{"test", 123} // Same values but different instance

		// First insertion
		idx1 := m.InsertFull(struct1, "value1")
		require.Equal(t, 0, idx1)

		// Update with equivalent struct shouldn't panic and should return same index
		idx2 := m.InsertFull(struct2, "value2")
		require.Equal(t, 0, idx2, "InsertFull should return the same index for equivalent structs")

		// Verify updated value
		val, ok := m.Get(struct1)
		require.True(t, ok)
		require.Equal(t, "value2", val)
	})

	t.Run("same value different byte arrays", func(t *testing.T) {
		m := indexmap.NewIndexMap[*[]byte, string]()

		bytes1 := []byte{1, 2, 3}
		bytes2 := []byte{1, 2, 3}

		m.InsertFull(&bytes1, "value1")

		idx, exists := m.Find(&bytes2)
		require.True(t, exists, "Should find the key with same values")
		require.Equal(t, 0, idx, "Index should be 0")

		val, exists := m.Get(&bytes2)
		require.True(t, exists, "Should get the value with same values")
		require.Equal(t, "value1", val)
	})

	t.Run("same value different structs", func(t *testing.T) {
		type TestStruct struct {
			Name string
			ID   int
		}

		m := indexmap.NewIndexMap[TestStruct, string]()

		struct1 := TestStruct{Name: "test", ID: 123}
		struct2 := TestStruct{Name: "test", ID: 123}

		m.InsertFull(struct1, "struct-value")

		idx, exists := m.Find(struct2)
		require.True(t, exists, "Should find the key with same values")
		require.Equal(t, 0, idx, "Index should be 0")

		val, exists := m.Get(struct2)
		require.True(t, exists, "Should get the value with same values")
		require.Equal(t, "struct-value", val)
	})

	t.Run("InsertFull with different byte arrays with same content", func(t *testing.T) {
		m := indexmap.NewIndexMap[*[]byte, int]()

		bytes1 := []byte{1, 2, 3, 4, 5}
		bytes2 := []byte{1, 2, 3, 4, 5}

		idx1 := m.InsertFull(&bytes1, 100)
		require.Equal(t, 0, idx1)

		idx2 := m.InsertFull(&bytes2, 200)
		require.Equal(t, 0, idx2, "Should return the same index for equivalent byte arrays")

		val, ok := m.Get(&bytes1)
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = m.Get(&bytes2)
		require.True(t, ok)
		require.Equal(t, 200, val)
	})
}
