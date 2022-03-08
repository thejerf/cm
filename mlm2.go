package cm

// MultiLevelMap2 is a 2-level multimap, with two key types and a value type.
type MultiLevelMap2[K1, K2 comparable, V any] map[K1]map[K2]V

// Tuple2 is a two-element tuple struct with a slot for each of the two keys.
type Tuple2[K1, K2 comparable] struct {
	Key1 K1
	Key2 K2
}

// KeyTree2 is a data type that can represent the keys of a map via a
// tree sort of structure.
type KeyTree2[K1, K2 comparable] struct {
	Key1  K1
	Key2s []K2
}

// Set will set the given value with the given keys.
//
// This will panic if called on a nil map.
func (mlm MultiLevelMap2[K1, K2, V]) Set(
	key1 K1,
	key2 K2,
	value V,
) {
	if mlm == nil {
		panic("Set called on a nil MultiLevelMap2")
	}
	l1 := mlm[key1]
	if l1 == nil {
		l1 = map[K2]V{}
		mlm[key1] = l1
	}
	l1[key2] = value
}

// SetByTuple sets by the key tuple.
func (mlm MultiLevelMap2[K1, K2, V]) SetByTuple(
	key Tuple2[K1, K2],
	value V,
) {
	if mlm == nil {
		panic("SetByTuple called on a nil MultiLevelMap2")
	}
	mlm.Set(key.Key1, key.Key2, value)
}

// Delete deletes the value from the map.
func (mlm MultiLevelMap2[K1, K2, V]) Delete(
	key1 K1,
	key2 K2,
) {
	if mlm == nil {
		return
	}
	l1 := mlm[key1]
	if l1 == nil {
		return
	}
	delete(l1, key2)

	if len(l1) == 0 {
		delete(mlm, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mlm MultiLevelMap2[K1, K2, V]) DeleteByTuple(key Tuple2[K1, K2]) {
	mlm.Delete(key.Key1, key.Key2)
}

// Get retreives by the given key. The second value is true if the key exists,
// false otherwise.
//
// Get called on a nil MultiLevelMap2 will not panic, and return that the
// value was not found.
func (mlm MultiLevelMap2[K1, K2, V]) Get(key1 K1, key2 K2) (val V, exists bool) {
	if mlm == nil {
		exists = false
		return
	}
	l1 := mlm[key1]
	if l1 == nil {
		exists = false
		return
	}
	v, exists := l1[key2]
	return v, exists
}

// GetByTuple retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
//
// GetByTuple called on a nil MultiLevelMap2 will not panic, and return
// that the value was not found.
func (mlm MultiLevelMap2[K1, K2, V]) GetByTuple(key Tuple2[K1, K2]) (val V, exists bool) {
	return mlm.Get(key.Key1, key.Key2)
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple2 values.
//
// A nil MultiLevelMap2 will return a nil slice.
func (mlm MultiLevelMap2[K1, K2, V]) KeySlice() []Tuple2[K1, K2] {
	if mlm == nil {
		return nil
	}

	r := []Tuple2[K1, K2]{}

	for key1, m1 := range mlm {
		for key2 := range m1 {
			r = append(r, Tuple2[K1, K2]{key1, key2})
		}
	}

	return r
}

// KeyTree returns the keys of the multimap as a 2-level tree of the
// various keys.
//
// A nil MultiLevelMap2 will return a nil slice.
func (mlm MultiLevelMap2[K1, K2, V]) KeyTree() []KeyTree2[K1, K2] {
	if mlm == nil {
		return nil
	}

	r := make([]KeyTree2[K1, K2], 0, len(mlm))
	for key1, m1 := range mlm {
		key2Slice := make([]K2, 0, len(m1))
		for key2 := range m1 {
			key2Slice = append(key2Slice, key2)
		}
		r = append(r, KeyTree2[K1, K2]{key1, key2Slice})
	}
	return r
}
