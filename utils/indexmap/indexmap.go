package indexmap

import "github.com/mitchellh/hashstructure/v2"

type IndexMap[K comparable, V any] struct {
	UnindexedMap    map[uint64]V
	InsertOrderList []K
}

func NewIndexMap[K comparable, V any]() *IndexMap[K, V] {
	return &IndexMap[K, V]{
		UnindexedMap:    make(map[uint64]V),
		InsertOrderList: []K{},
	}
}

// Set method to add or update a key-value pair and maintain insertion order
func (m *IndexMap[K, V]) Set(key K, value V) {
	hash := GetHash(key)
	if _, exists := m.UnindexedMap[hash]; !exists {
		m.InsertOrderList = append(m.InsertOrderList, key)
	}
	m.UnindexedMap[hash] = value
}

// Insert method to add or update a key-value pair and maintain insertion order
func (m *IndexMap[K, V]) Insert(key K, value V) {
	m.InsertFull(key, value)
}

// Insert method to add or update a key-value pair and maintain insertion order and
func (m *IndexMap[K, V]) InsertFull(key K, value V) int {
	hash := GetHash(key)
	_, exists := m.UnindexedMap[hash]
	m.UnindexedMap[hash] = value
	if !exists {
		m.InsertOrderList = append(m.InsertOrderList, key)
		return len(m.InsertOrderList) - 1
	}

	idx, exists := m.Find(key)
	if !exists {
		panic("key didn't successfully add")
	}
	return idx
}

// Get method to retrieve the value for a given key
func (m *IndexMap[K, V]) Get(key K) (V, bool) {
	hash := GetHash(key)
	value, exists := m.UnindexedMap[hash]
	return value, exists
}

func (m *IndexMap[K, V]) Find(key K) (int, bool) {
	for i, v := range m.InsertOrderList {
		if v == key {
			return i, true
		}
	}
	return 0, false
}

func (m *IndexMap[K, V]) Index(i int) (*K, bool) {
	if len(m.InsertOrderList)-1 < i {
		k := m.InsertOrderList[i]
		return &k, true
	}
	return nil, false
}

func (m *IndexMap[K, V]) Len() int {
	return len(m.UnindexedMap)
}

// ForEach method to iterate over all key-value pairs in insertion order
func (m *IndexMap[K, V]) ForEach(action func(K, V)) {
	for _, key := range m.InsertOrderList {
		hash := GetHash(key)
		action(key, m.UnindexedMap[hash])
	}
}

// return a deterministic hash of struct/any datatype
func GetHash(obj any) uint64 {
	// TODO we can implement our own hash func for go structs
	hash, err := hashstructure.Hash(obj, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return hash
}
