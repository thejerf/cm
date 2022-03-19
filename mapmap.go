package cm

import (
	"golang.org/x/exp/maps"
)

// MapMap is a map of maps that has a comparable Value type, which allows
// for the .Equal method.
//
// AnyMapMap lacks this restriction, and has all the methods MapMap has
// except .Equal.
type MapMap[K1, K2, V comparable] MapMapAny[K1, K2, V]

// Equal returns if this MapMap is equal to the passed-in MapMap.
//
// Two zero-sized maps are considered equal to each other, even if one is
// nil and the other is not. This matches the current behavior of
// maps.Equal. If that changes before release, this will change to match it.
func (mm MapMap[K1, K2, V]) Equal(r MapMap[K1, K2, V]) bool {
	if mm == nil && r == nil {
		return true
	}
	if len(mm) != len(r) {
		return false
	}
	if len(mm) == 0 {
		return true
	}

	for key, submap := range mm {
		rightSubmap, exists := r[key]
		if !exists {
			return false
		}
		if !maps.Equal(submap, rightSubmap) {
			return false
		}
	}

	return true
}

// Clone returns a copy of the MapMap struture, with the keys copied
// over. It's a shallow copy of the full MapMap.
func (mm MapMap[K1, K2, V]) Clone() MapMap[K1, K2, V] {
	newMap := make(MapMap[K1, K2, V], len(mm))

	for key1, submap := range mm {
		newMap[key1] = maps.Clone(submap)
	}

	return newMap
}

// Set will set the given value with the given keys.
//
// This will panic if called on a nil map.
func (mm MapMap[K1, K2, V]) Set(key1 K1, key2 K2, value V) {
	MapMapAny[K1, K2, V](mm).Set(key1, key2, value)
}

// SetByTuple sets by the key tuple.
func (mm MapMap[K1, K2, V]) SetByTuple(key Tuple2[K1, K2], value V) {
	MapMapAny[K1, K2, V](mm).SetByTuple(key, value)
}

// Delete deletes the value from the map.
func (mm MapMap[K1, K2, V]) Delete(key1 K1, key2 K2) {
	MapMapAny[K1, K2, V](mm).Delete(key1, key2)
}

// DeleteByKey deletes by the tuple version of the key.
func (mm MapMap[K1, K2, V]) DeleteByTuple(key Tuple2[K1, K2]) {
	MapMapAny[K1, K2, V](mm).DeleteByTuple(key)
}

// GetByTuple retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
//
// GetByTuple called on a nil MapMap will not panic, and return
// that the value was not found.
func (mm MapMap[K1, K2, V]) GetByTuple(key Tuple2[K1, K2]) (val V, exists bool) {
	return MapMapAny[K1, K2, V](mm).GetByTuple(key)
}

// EqualFunc reimplements maps.EqualFunc on the MapMap.
func (mm MapMap[K1, K2, V]) EqualFunc(
	r MapMap[K1, K2, V],
	eq func(v1, v2 V) bool,
) bool {
	return MapMapAny[K1, K2, V](mm).EqualFunc(
		MapMapAny[K1, K2, V](r),
		eq,
	)
}

// DeleteFunc deletes from the map the values for which the function
// returns true. If all values from a submap are deleted, the submap will
// be deleted from the MapMap.
func (mm MapMap[K1, K2, V]) DeleteFunc(f func(K1, K2, V) bool) {
	MapMapAny[K1, K2, V](mm).DeleteFunc(f)
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple2 values.
//
// A nil map will return a nil slice.
func (mm MapMap[K1, K2, V]) KeySlice() []Tuple2[K1, K2] {
	return MapMapAny[K1, K2, V](mm).KeySlice()
}

// KeyTree returns the keys of the multimap as a 2-level tree of the
// various keys.
//
// A nil map will return a nil slice.
func (mm MapMap[K1, K2, V]) KeyTree() []KeyTree[K1, K2] {
	return MapMapAny[K1, K2, V](mm).KeyTree()
}

// MapMapAny is a 2-level multimap, with two key types and a value type.
type MapMapAny[K1, K2 comparable, V any] map[K1]map[K2]V

// Tuple2 is a two-element tuple struct with a slot for each of the two keys.
type Tuple2[K1, K2 comparable] struct {
	Key1 K1
	Key2 K2
}

// KeyTree is a data type that can represent the keys of a map via a
// tree sort of structure.
type KeyTree[K1, V any] struct {
	Key  K1
	Vals []V
}

// Set will set the given value with the given keys.
//
// This will panic if called on a nil map.
func (mma MapMapAny[K1, K2, V]) Set(
	key1 K1,
	key2 K2,
	value V,
) {
	if mma == nil {
		panic("Set called on a nil MapMap")
	}
	l1 := mma[key1]
	if l1 == nil {
		l1 = map[K2]V{}
		mma[key1] = l1
	}
	l1[key2] = value
}

// SetByTuple sets by the key tuple.
func (mma MapMapAny[K1, K2, V]) SetByTuple(
	key Tuple2[K1, K2],
	value V,
) {
	if mma == nil {
		panic("SetByTuple called on a nil MapMap")
	}
	mma.Set(key.Key1, key.Key2, value)
}

// Delete deletes the value from the map.
func (mma MapMapAny[K1, K2, V]) Delete(
	key1 K1,
	key2 K2,
) {
	if mma == nil {
		return
	}
	l1 := mma[key1]
	if l1 == nil {
		return
	}
	delete(l1, key2)

	if len(l1) == 0 {
		delete(mma, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mma MapMapAny[K1, K2, V]) DeleteByTuple(key Tuple2[K1, K2]) {
	mma.Delete(key.Key1, key.Key2)
}

// GetByTuple retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
//
// GetByTuple called on a nil MapMap will not panic, and return
// that the value was not found.
func (mma MapMapAny[K1, K2, V]) GetByTuple(key Tuple2[K1, K2]) (val V, exists bool) {
	val, exists = mma[key.Key1][key.Key2]
	return val, exists
}

// Clone returns a copy of the MapMap struture, with the keys copied
// over. It's a shallow copy of the full MapMap.
func (mma MapMapAny[K1, K2, V]) Clone() MapMapAny[K1, K2, V] {
	newMap := make(MapMapAny[K1, K2, V], len(mma))

	for key1, submap := range mma {
		newMap[key1] = maps.Clone(submap)
	}

	return newMap
}

// EqualFunc reimplements maps.EqualFunc on the MapMap.
func (mma MapMapAny[K1, K2, V]) EqualFunc(
	r MapMapAny[K1, K2, V],
	eq func(v1, v2 V) bool,
) bool {
	if mma == nil && r == nil {
		return true
	}
	if len(mma) != len(r) {
		return false
	}
	if len(mma) == 0 {
		return true
	}

	for key, submap := range mma {
		rightSubmap, exists := r[key]
		if !exists {
			return false
		}
		if !maps.EqualFunc(submap, rightSubmap, eq) {
			return false
		}
	}
	return true
}

// DeleteFunc deletes from the map the values for which the function
// returns true. If all values from a submap are deleted, the submap will
// be deleted from the MapMap.
func (mma MapMapAny[K1, K2, V]) DeleteFunc(f func(K1, K2, V) bool) {
	if mma == nil {
		return
	}

	for key1, submap := range mma {
		delF := func(key2 K2, val V) bool {
			return f(key1, key2, val)
		}
		maps.DeleteFunc(submap, delF)
		if len(submap) == 0 {
			delete(mma, key1)
		}
	}
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple2 values.
//
// A nil map will return a nil slice.
func (mma MapMapAny[K1, K2, V]) KeySlice() []Tuple2[K1, K2] {
	if mma == nil {
		return nil
	}

	r := []Tuple2[K1, K2]{}

	for key1, m1 := range mma {
		for key2 := range m1 {
			r = append(r, Tuple2[K1, K2]{key1, key2})
		}
	}

	return r
}

// KeyTree returns the keys of the multimap as a 2-level tree of the
// various keys.
//
// A nil map will return a nil slice.
func (mma MapMapAny[K1, K2, V]) KeyTree() []KeyTree[K1, K2] {
	if mma == nil {
		return nil
	}

	r := make([]KeyTree[K1, K2], 0, len(mma))
	for key1, m1 := range mma {
		key2Slice := make([]K2, 0, len(m1))
		for key2 := range m1 {
			key2Slice = append(key2Slice, key2)
		}
		r = append(r, KeyTree[K1, K2]{key1, key2Slice})
	}
	return r
}
