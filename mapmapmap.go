package cm

// MLM3 is a 3-level multimap, with three key types and a value type.
//
// Note this recursively includes an MultiLevelMap2 as its values, meaning
// that all values fetched from this value has all those methods as well.
type MapMapMap[K1, K2, K3 comparable, V any] map[K1]MapMap[K2, K3, V]

// Tuple3 is a three-element tuple struct with a slot for each of the two keys.
type Tuple3[K1, K2, K3 comparable] struct {
	Key1 K1
	Key2 K2
	Key3 K3
}

// Set will set the given value with the given keys.
func (mmm MapMapMap[K1, K2, K3, V]) Set(
	key1 K1,
	key2 K2,
	key3 K3,
	value V,
) {
	l1 := mmm[key1]
	if l1 == nil {
		l1 = MapMap[K2, K3, V]{}
		mmm[key1] = l1
	}
	l1.Set(key2, key3, value)
}

// SetByKey sets by the key tuple.
func (mmm MapMapMap[K1, K2, K3, V]) SetByTuple(
	key Tuple3[K1, K2, K3],
	value V,
) {
	mmm.Set(key.Key1, key.Key2, key.Key3, value)
}

// Delete deletes the value from the map.
func (mmm MapMapMap[K1, K2, K3, V]) Delete(
	key1 K1,
	key2 K2,
	key3 K3,
) {
	l1 := mmm[key1]
	if l1 == nil {
		return
	}
	l1.Delete(key2, key3)
	if len(l1) == 0 {
		delete(mmm, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mmm MapMapMap[K1, K2, K3, V]) DeleteByTuple(key Tuple3[K1, K2, K3]) {
	mmm.Delete(key.Key1, key.Key2, key.Key3)
}

// Get retreives by the given key. The second value is true if the key exists,
// false otherwise.
func (mmm MapMapMap[K1, K2, K3, V]) Get(key1 K1, key2 K2, key3 K3) (val V, exists bool) {
	l1 := mmm[key1]
	if l1 == nil {
		exists = false
		return
	}
	return l1.Get(key2, key3)
}

// Get retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
func (mmm MapMapMap[K1, K2, K3, V]) GetByTuple(key Tuple3[K1, K2, K3]) (val V, exists bool) {
	return mmm.Get(key.Key1, key.Key2, key.Key3)
}

// EqualFunc returns if this MapMapMap is equal to the passed-in MapMapMap,
// using the passed-in function to compare the equality of values,
// re-implementing maps.EqualFunc on a MapMapMap.
func (mmm MapMapMap[K1, K2, K3, V]) EqualFunc(
	r MapMapMap[K1, K2, K3, V],
	eq func(v1, v2 V) bool,
) bool {
	if mmm == nil && r == nil {
		return true
	}
	if len(mmm) != len(r) {
		return false
	}
	if len(mmm) == 0 {
		return true
	}

	for key, submap := range mmm {
		rightSubmap, exists := r[key]
		if !exists {
			return false
		}
		if !submap.EqualFunc(rightSubmap, eq) {
			return false
		}
	}

	return true
}

// DeleteFunc deletes from the map the values for which the function
// returns true. If all values from a submap are deleted, the submap will
// be deleted from the MapMapMap.
func (mmm MapMapMap[K1, K2, K3, V]) DeleteFunc(f func(K1, K2, K3, V) bool) {
	if mmm == nil {
		return
	}

	for key1, submap := range mmm {
		delF := func(key2 K2, key3 K3, val V) bool {
			return f(key1, key2, key3, val)
		}
		submap.DeleteFunc(delF)
		if len(submap) == 0 {
			delete(mmm, key1)
		}
	}
}

// GetThirdLevel safely fetches the third level of the multimap for the
// given first two keys.
//
// Fetching the "second-level" map can be done by simply looking up the key1
// value in the map, yielding either the corresponding MapMap or nil. This
// method safely fetches the third-level map, or returns nil if either key1 or
// (key1, key2) don't exist.
func (mmm MapMapMap[K1, K2, K3, V]) GetThirdLevel(key1 K1, key2 K2) map[K3]V {
	l1 := mmm[key1]
	if l1 == nil {
		return nil
	}
	return l1[key2]
}

// Clone yields a shallow copy of the MapMapMap, with the values simply
// copied across.
func (mmm MapMapMap[K1, K2, K3, V]) Clone() MapMapMap[K1, K2, K3, V] {
	newMMM := make(MapMapMap[K1, K2, K3, V], len(mmm))

	for key1, mapmap := range mmm {
		newMMM[key1] = mapmap.Clone()
	}

	return newMMM
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple3 values.
func (mmm MapMapMap[K1, K2, K3, V]) KeySlice() []Tuple3[K1, K2, K3] {
	r := []Tuple3[K1, K2, K3]{}

	for key1, m1 := range mmm {
		for key2, m2 := range m1 {
			for key3 := range m2 {
				r = append(r, Tuple3[K1, K2, K3]{key1, key2, key3})
			}
		}
	}

	return r
}

// KeyTree returns the keys of the multimap as a 3-level tree of the
// various keys.
func (mmm MapMapMap[K1, K2, K3, V]) KeyTree() []KeyTree[K1, KeyTree[K2, K3]] {
	r := make([]KeyTree[K1, KeyTree[K2, K3]], 0, len(mmm))
	for key1, m1 := range mmm {
		keyTree := MapMap[K2, K3, V](m1).KeyTree()
		r = append(r, KeyTree[K1, KeyTree[K2, K3]]{key1, keyTree})
	}
	return r
}
