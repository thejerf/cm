package cm

// MapMapMap is a map of map of maps that has a comparable Value type,
// which allows for the .Equal method.
//
// MapMapMapAny lacks this restriction, and has all the methods MapMapMap
// has except .Equal.
type MapMapMap[K1, K2, K3, V comparable] MapMapMapAny[K1, K2, K3, V]

// Equal returns if this MapMapMap is equal to the passed-in MapMapMap.
//
// Two zero-sized maps are considered equal to each other, even if one is
// nil and the other is not. This matches the current behavior of
// maps.Equal. If that changes before release, this will change to match
// it.
func (mmm MapMapMap[K1, K2, K3, V]) Equal(r MapMapMap[K1, K2, K3, V]) bool {
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
		if !MapMap[K2, K3, V](submap).Equal(
			MapMap[K2, K3, V](rightSubmap)) {
			return false
		}
	}

	return true
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

// Set will set the given value with the given keys.
//
// This will panic if called on a nil map.
func (mmm MapMapMap[K1, K2, K3, V]) Set(key1 K1, key2 K2, key3 K3, value V) {
	MapMapMapAny[K1, K2, K3, V](mmm).Set(key1, key2, key3, value)
}

// SetByTuple sets by the key tuple.
func (mmm MapMapMap[K1, K2, K3, V]) SetByTuple(key Tuple3[K1, K2, K3], value V) {
	MapMapMapAny[K1, K2, K3, V](mmm).SetByTuple(key, value)
}

// Delete deletes the value from the map.
func (mmm MapMapMap[K1, K2, K3, V]) Delete(key1 K1, key2 K2, key3 K3) {
	MapMapMapAny[K1, K2, K3, V](mmm).Delete(key1, key2, key3)
}

// DeleteByKey deletes by the tuple version of the key.
func (mmm MapMapMap[K1, K2, K3, V]) DeleteByTuple(key Tuple3[K1, K2, K3]) {
	MapMapMapAny[K1, K2, K3, V](mmm).DeleteByTuple(key)
}

// GetByTuple retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
//
// GetByTuple called on a nil MapMap will not panic, and return
// that the value was not found.
func (mmm MapMapMap[K1, K2, K3, V]) GetByTuple(key Tuple3[K1, K2, K3]) (val V, exists bool) {
	return MapMapMapAny[K1, K2, K3, V](mmm).GetByTuple(key)
}

// EqualFunc reimplements maps.EqualFunc on the MapMap.
func (mmm MapMapMap[K1, K2, K3, V]) EqualFunc(
	r MapMapMap[K1, K2, K3, V],
	eq func(v1, v2 V) bool,
) bool {
	return MapMapMapAny[K1, K2, K3, V](mmm).EqualFunc(
		MapMapMapAny[K1, K2, K3, V](r),
		eq,
	)
}

// DeleteFunc deletes from the map the values for which the function
// returns true. If all values from a submap are deleted, the submap will
// be deleted from the MapMap.
func (mmm MapMapMap[K1, K2, K3, V]) DeleteFunc(f func(K1, K2, K3, V) bool) {
	MapMapMapAny[K1, K2, K3, V](mmm).DeleteFunc(f)
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple2 values.
//
// A nil map will return a nil slice.
func (mmm MapMapMap[K1, K2, K3, V]) KeySlice() []Tuple3[K1, K2, K3] {
	return MapMapMapAny[K1, K2, K3, V](mmm).KeySlice()
}

// KeyTree returns the keys of the multimap as a 2-level tree of the
// various keys.
//
// A nil map will return a nil slice.
func (mmm MapMapMap[K1, K2, K3, V]) KeyTree() []KeyTree[K1, KeyTree[K2, K3]] {
	return MapMapMapAny[K1, K2, K3, V](mmm).KeyTree()
}

// MapMapMapAny is a three-level map that can contain any value.
type MapMapMapAny[K1, K2, K3 comparable, V any] map[K1]MapMapAny[K2, K3, V]

// Tuple3 is a three-element tuple struct with a slot for each of the two keys.
type Tuple3[K1, K2, K3 comparable] struct {
	Key1 K1
	Key2 K2
	Key3 K3
}

// Set will set the given value with the given keys.
func (mmma MapMapMapAny[K1, K2, K3, V]) Set(
	key1 K1,
	key2 K2,
	key3 K3,
	value V,
) {
	l1 := mmma[key1]
	if l1 == nil {
		l1 = MapMapAny[K2, K3, V]{}
		mmma[key1] = l1
	}
	l1.Set(key2, key3, value)
}

// SetByKey sets by the key tuple.
func (mmma MapMapMapAny[K1, K2, K3, V]) SetByTuple(
	key Tuple3[K1, K2, K3],
	value V,
) {
	mmma.Set(key.Key1, key.Key2, key.Key3, value)
}

// Delete deletes the value from the map.
func (mmma MapMapMapAny[K1, K2, K3, V]) Delete(
	key1 K1,
	key2 K2,
	key3 K3,
) {
	l1 := mmma[key1]
	if l1 == nil {
		return
	}
	l1.Delete(key2, key3)
	if len(l1) == 0 {
		delete(mmma, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mmma MapMapMapAny[K1, K2, K3, V]) DeleteByTuple(key Tuple3[K1, K2, K3]) {
	mmma.Delete(key.Key1, key.Key2, key.Key3)
}

// GetByTuple is a convenience function to retrieve the value out of the
// map by the key tuple returned by KeySlice. Normal usage should just
// index into the map like mmm[key1][key2][key3].
func (mmma MapMapMapAny[K1, K2, K3, V]) GetByTuple(key Tuple3[K1, K2, K3]) (val V, exists bool) {
	val, exists = mmma[key.Key1][key.Key2][key.Key3]
	return val, exists
}

// EqualFunc returns if this MapMapMap is equal to the passed-in MapMapMap,
// using the passed-in function to compare the equality of values,
// re-implementing maps.EqualFunc on a MapMapMap.
func (mmma MapMapMapAny[K1, K2, K3, V]) EqualFunc(
	r MapMapMapAny[K1, K2, K3, V],
	eq func(v1, v2 V) bool,
) bool {
	if mmma == nil && r == nil {
		return true
	}
	if len(mmma) != len(r) {
		return false
	}
	if len(mmma) == 0 {
		return true
	}

	for key, submap := range mmma {
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
func (mmma MapMapMapAny[K1, K2, K3, V]) DeleteFunc(f func(K1, K2, K3, V) bool) {
	if mmma == nil {
		return
	}

	for key1, submap := range mmma {
		delF := func(key2 K2, key3 K3, val V) bool {
			return f(key1, key2, key3, val)
		}
		submap.DeleteFunc(delF)
		if len(submap) == 0 {
			delete(mmma, key1)
		}
	}
}

// Clone yields a shallow copy of the MapMapMapAny, with the values simply
// copied across.
func (mmma MapMapMapAny[K1, K2, K3, V]) Clone() MapMapMapAny[K1, K2, K3, V] {
	newMMM := make(MapMapMapAny[K1, K2, K3, V], len(mmma))

	for key1, mapmap := range mmma {
		newMMM[key1] = mapmap.Clone()
	}

	return newMMM
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple3 values.
func (mmma MapMapMapAny[K1, K2, K3, V]) KeySlice() []Tuple3[K1, K2, K3] {
	r := []Tuple3[K1, K2, K3]{}

	for key1, m1 := range mmma {
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
func (mmma MapMapMapAny[K1, K2, K3, V]) KeyTree() []KeyTree[K1, KeyTree[K2, K3]] {
	r := make([]KeyTree[K1, KeyTree[K2, K3]], 0, len(mmma))
	for key1, m1 := range mmma {
		keyTree := MapMapAny[K2, K3, V](m1).KeyTree()
		r = append(r, KeyTree[K1, KeyTree[K2, K3]]{key1, keyTree})
	}
	return r
}
