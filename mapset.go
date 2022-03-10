package cm

// A MapSet is a map that contains sets. For example, a map of users to the
// set of resources they are allowed to access.
//
// Direct access is legal for most operations, except that if you delete
// the last element out of the set with Delete, the set
// will be entirely removed from the MapSet, but if you delete manually it
// won't unless you clean it up yourself. Reading and writing directly is
// fine.
//
// Perhaps the most subtle is, if you want something like "KeyTree" but
// that returns sets, the MapSet is itself that thing. Just:
//
//   for key, set := range MyMapSet {
//       // and here you have the key and the set accessible
//   }
//
// Since the nil set is mostly functional, many operations can be performed
// naturally simply by indexing the top-level map, even if that results in
// a nil set. This package only implements methods that are value-adds on
// top of direct map access, such as Set, which creates the map entry to a
// valid set if necessary.
type MapSet[K, V comparable] map[K]Set[V]

// AllValueSet returns a single set containing all values in the MapSet.
func (ms MapSet[K, V]) AllValueSet() Set[V] {
	var retSet Set[V]

	for _, set := range ms {
		if retSet == nil {
			retSet = set.Clone()
		} else {
			retSet.Union(set)
		}
	}

	return retSet
}

func (ms MapSet[K, V]) Set(key K, val V) {
	if ms == nil {
		panic("Set called on a nil MapSet")
	}
	s := ms[key]
	if s == nil {
		s = Set[V]{}
		ms[key] = s
	}
	s[val] = void
}

func (ms MapSet[K, V]) Contains(key K, val V) bool {
	if ms == nil {
		return false
	}
	s := ms[key]
	if s == nil {
		return false
	}
	return s.Contains(val)
}

func (ms MapSet[K, V]) Delete(key K, val V) {
	if ms == nil {
		return
	}
	s := ms[key]
	if s == nil {
		return
	}
	delete(s, val)
	if len(s) == 0 {
		delete(ms, key)
	}
}
