package cm

// MLM3 is a 3-level multimap, with three key types and a value type.
//
// Note this recursively includes an MultiLevelMap2 as its values, meaning
// that all values fetched from this value has all those methods as well.
type MultiLevelMap3[K1, K2, K3 comparable, V any] map[K1]MultiLevelMap2[K2, K3, V]

// Tuple3 is a three-element tuple struct with a slot for each of the two keys.
type Tuple3[K1, K2, K3 comparable] struct {
	Key1 K1
	Key2 K2
	Key3 K3
}

// Set will set the given value with the given keys.
func (mlm MultiLevelMap3[K1, K2, K3, V]) Set(
	key1 K1,
	key2 K2,
	key3 K3,
	value V,
) {
	l1 := mlm[key1]
	if l1 == nil {
		l1 = MultiLevelMap2[K2, K3, V]{}
		mlm[key1] = l1
	}
	l1.Set(key2, key3, value)
}

// SetByKey sets by the key tuple.
func (mlm MultiLevelMap3[K1, K2, K3, V]) SetByTuple(
	key Tuple3[K1, K2, K3],
	value V,
) {
	mlm.Set(key.Key1, key.Key2, key.Key3, value)
}

// Delete deletes the value from the map.
func (mlm MultiLevelMap3[K1, K2, K3, V]) Delete(
	key1 K1,
	key2 K2,
	key3 K3,
) {
	l1 := mlm[key1]
	if l1 == nil {
		return
	}
	l1.Delete(key2, key3)
	if len(l1) == 0 {
		delete(mlm, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mlm MultiLevelMap3[K1, K2, K3, V]) DeleteByTuple(key Tuple3[K1, K2, K3]) {
	mlm.Delete(key.Key1, key.Key2, key.Key3)
}

// Get retreives by the given key. The second value is true if the key exists,
// false otherwise.
func (mlm MultiLevelMap3[K1, K2, K3, V]) Get(key1 K1, key2 K2, key3 K3) (val V, exists bool) {
	l1 := mlm[key1]
	if l1 == nil {
		exists = false
		return
	}
	return l1.Get(key2, key3)
}

// Get retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
func (mlm MultiLevelMap3[K1, K2, K3, V]) GetByTuple(key Tuple3[K1, K2, K3]) (val V, exists bool) {
	return mlm.Get(key.Key1, key.Key2, key.Key3)
}

// GetThirdLevel safely fetches the third level of the multimap for the
// given first two keys.
//
// Fetching the second-level map can be done by simply looking up the key1
// value in the map. This safely fetches the third-level map, or returns
// nil if either key1 or (key1, key2) don't exist.
func (mlm MultiLevelMap3[K1, K2, K3, V]) GetThirdLevel(key1 K1, key2 K2) map[K3]V {
	l1 := mlm[key1]
	if l1 == nil {
		return nil
	}
	return l1[key2]
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple3 values.
func (mlm MultiLevelMap3[K1, K2, K3, V]) KeySlice() []Tuple3[K1, K2, K3] {
	r := []Tuple3[K1, K2, K3]{}

	for key1, m1 := range mlm {
		for key2, m2 := range m1 {
			for key3 := range m2 {
				r = append(r, Tuple3[K1, K2, K3]{key1, key2, key3})
			}
		}
	}

	return r
}

type KeyTree3[K1, K2, K3 comparable] struct {
	Key1  K1
	Value []KeyTree2[K2, K3]
}

// KeyTree returns the keys of the multimap as a 3-level tree of the
// various keys.
func (mlm MultiLevelMap3[K1, K2, K3, V]) KeyTree() []KeyTree3[K1, K2, K3] {
	r := make([]KeyTree3[K1, K2, K3], 0, len(mlm))
	for key1, m1 := range mlm {
		keyTree := MultiLevelMap2[K2, K3, V](m1).KeyTree()
		r = append(r, KeyTree3[K1, K2, K3]{key1, keyTree})
	}
	return r
}
