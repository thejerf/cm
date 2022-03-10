package cm

import (
	"golang.org/x/exp/maps"
)

// MapMap is a 2-level multimap, with two key types and a value type.
type MapMap[K1, K2 comparable, V any] map[K1]map[K2]V

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
func (mm MapMap[K1, K2, V]) Set(
	key1 K1,
	key2 K2,
	value V,
) {
	if mm == nil {
		panic("Set called on a nil MapMap")
	}
	l1 := mm[key1]
	if l1 == nil {
		l1 = map[K2]V{}
		mm[key1] = l1
	}
	l1[key2] = value
}

// SetByTuple sets by the key tuple.
func (mm MapMap[K1, K2, V]) SetByTuple(
	key Tuple2[K1, K2],
	value V,
) {
	if mm == nil {
		panic("SetByTuple called on a nil MapMap")
	}
	mm.Set(key.Key1, key.Key2, value)
}

// Delete deletes the value from the map.
func (mm MapMap[K1, K2, V]) Delete(
	key1 K1,
	key2 K2,
) {
	if mm == nil {
		return
	}
	l1 := mm[key1]
	if l1 == nil {
		return
	}
	delete(l1, key2)

	if len(l1) == 0 {
		delete(mm, key1)
	}
}

// DeleteByKey deletes by the tuple version of the key.
func (mm MapMap[K1, K2, V]) DeleteByTuple(key Tuple2[K1, K2]) {
	mm.Delete(key.Key1, key.Key2)
}

// Get retreives by the given key. The second value is true if the key exists,
// false otherwise.
func (mm MapMap[K1, K2, V]) Get(key1 K1, key2 K2) (val V, exists bool) {
	if mm == nil {
		panic("Get called on a nil MapMap")
	}
	l1 := mm[key1]
	if l1 == nil {
		exists = false
		return
	}
	val, exists = l1[key2]
	return
}

// GetByTuple retreives by the given tuple. The second value is true if the key
// exists, false otherwise.
//
// GetByTuple called on a nil MapMap will not panic, and return
// that the value was not found.
func (mm MapMap[K1, K2, V]) GetByTuple(key Tuple2[K1, K2]) (val V, exists bool) {
	return mm.Get(key.Key1, key.Key2)
}

// Clone returns a shallow copy of this MapMap, where the MapMap structure
// is retained and the values are just copied over.
func (mm MapMap[K1, K2, V]) Clone() MapMap[K1, K2, V] {
	newMap := make(MapMap[K1, K2, V], len(mm))

	for key1, submap := range mm {
		newMap[key1] = maps.Clone(submap)
	}

	return newMap
}

// CompMapMap is a type that a MapMap that has a comparable Value type can
// be converted into to allow you to run the .Equal operation.
//
// Because the MapMap's Value type does not implement comparable (to not
// limit what can be stored in it), it can not implement the Equal method
// as the value of a MapMap is not comparable, nor can a MapMap
// automatically coerce itself into this without having to instantiate an
// illegal type. (If you can figure out how to do it, let me know in a
// PR. I'm avoiding reflection, too slow.)
//
// This is logically equivalant to calling MapMap.EqualFunc(otherMap,
// func (a, b V) bool { return a == b }), except that MapMap.EqualFunc will
// result in a potentially large number of calls to that function, whereas
// this is able to more tightly loop through the map without an additional
// non-inlinable function call per value that has to be checked. For small
// maps the convenience of the EqualsFunc call may be worth it, but this is
// included for places where the performance is advantageous.
//
// This allows invoking the Equal method on a MapMap via the somewhat
// circuitous:
//
//   mm1 := MapMap[int, int, int]{1: 2}
//   mm2 := MapMap[int, int, int]{1: 2}
//
//   CompMapMap[int, int, int](mm1).Equal(CompMapMap[int, int, int](mm2))
//
// If you do this often, local code aware of the concrete type can abstract
// this, because your local package can see the value type is
// comparable. The cm package can not.
type CompMapMap[K1, K2, V comparable] map[K1]map[K2]V

// Equal returns if this MapMap is equal to the passed-in MapMap.
//
// Two zero-sized maps are considered equal to each other, even if one is
// nil and the other is not. This matches the current behavior of
// maps.Equal. If that changes before release, this will change to match it.
func (cmm CompMapMap[K1, K2, V]) Equal(r CompMapMap[K1, K2, V]) bool {
	if cmm == nil && r == nil {
		return true
	}
	if len(cmm) != len(r) {
		return false
	}
	if len(cmm) == 0 {
		return true
	}

	for key, submap := range cmm {
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

// EqualFunc reimplements maps.EqualFunc on the MapMap.
func (mm MapMap[K1, K2, V]) EqualFunc(
	r MapMap[K1, K2, V],
	eq func(v1, v2 V) bool,
) bool {
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
		if !maps.EqualFunc(submap, rightSubmap, eq) {
			return false
		}
	}
	return true
}

// DeleteFunc deletes from the map the values for which the function
// returns true. If all values from a submap are deleted, the submap will
// be deleted from the MapMap.
func (mm MapMap[K1, K2, V]) DeleteFunc(f func(K1, K2, V) bool) {
	if mm == nil {
		return
	}

	for key1, submap := range mm {
		delF := func(key2 K2, val V) bool {
			return f(key1, key2, val)
		}
		maps.DeleteFunc(submap, delF)
		if len(submap) == 0 {
			delete(mm, key1)
		}
	}
}

// KeySlice returns the keys of the mulitmap as a slice of Tuple2 values.
//
// A nil MultiLevelMap2 will return a nil slice.
func (mm MapMap[K1, K2, V]) KeySlice() []Tuple2[K1, K2] {
	if mm == nil {
		return nil
	}

	r := []Tuple2[K1, K2]{}

	for key1, m1 := range mm {
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
func (mm MapMap[K1, K2, V]) KeyTree() []KeyTree[K1, K2] {
	if mm == nil {
		return nil
	}

	r := make([]KeyTree[K1, K2], 0, len(mm))
	for key1, m1 := range mm {
		key2Slice := make([]K2, 0, len(m1))
		for key2 := range m1 {
			key2Slice = append(key2Slice, key2)
		}
		r = append(r, KeyTree[K1, K2]{key1, key2Slice})
	}
	return r
}
