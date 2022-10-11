package cm

// A Set wraps a map[Type]struct{} with various convenience functions to
// turn it into a set.
//
// As with the other data structures in this package, redundant operations
// with map are not implemented. To get the contents of this set, use
// maps.Keys. to clear a set, use maps.Clear, to get the number of elements
// in a set, use length(), to iterate over the set, use range in a for
// loop, etc. This only implements additional useful set functionality.
//
// This Set offer high-efficiency mutating operations on the map; for
// instance, Union will copy the target Set into the Set the method is
// called on. To get the behavior where a separate Set is returned that is
// a union, add .Clone() to the appropriate place in your code.
//
// It is safe to directly manipulate a set using map syntax in Go. This
// data type does not make any guarantees that would be violated.
//
// The nil Set is mostly legal and functions as an empty set, except you
// can not .Add or .Union it.
type Set[M comparable] map[M]struct{}

var void = struct{}{}

// SetFromSlice loads a set in from a slice.
//
// It is not appropriate or indeed even possible for a library to select
// how to serialize specific sets, but it can be helpful for JSON or YAML
// serialization to embed a specific Set and then use this and AsSlice as a
// pair to serialize a set as a slice instead of a map.
func SetFromSlice[M comparable](l []M) Set[M] {
	s := Set[M]{}

	for _, val := range l {
		s[val] = void
	}

	return s
}

// Add will add the given value in to the set.
func (s Set[M]) Add(v M) {
	if s == nil {
		panic("Add called on nil Set")
	}
	s[v] = void
}

// AsSlice returns the set as a slice, in hash order.
//
// See comment on SetFromSlice.
func (s Set[M]) AsSlice() []M {
	vals := make([]M, 0, len(s))
	for val := range s {
		vals = append(vals, val)
	}
	return vals
}

// Clone returns a copy of the set.
func (s Set[M]) Clone() Set[M] {
	if s == nil {
		return nil
	}
	newSet := make(map[M]struct{}, len(s))
	for value := range s {
		newSet[value] = void
	}
	return newSet
}

// Contains returns true if the set contains the given value.
func (s Set[M]) Contains(v M) bool {
	if s == nil {
		return false
	}
	_, exists := s[v]
	return exists
}

// Equal returns if two sets are equal.
func (s Set[M]) Equal(r Set[M]) bool {
	if s == nil && r == nil {
		return true
	}
	if len(s) != len(r) {
		return false
	}
	if len(s) == 0 {
		return true
	}

	for key := range s {
		if !r.Contains(key) {
			return false
		}
	}

	return true
}

// SubsetOf returns true if the set this is called on is a subset of the
// passed-in set.
func (s Set[M]) SubsetOf(r Set[M]) bool {
	if len(s) > len(r) {
		return false
	}

	for key := range s {
		if !r.Contains(key) {
			return false
		}
	}

	return true
}

// SupersetOf returns true if the set this is called on is a superset of
// the passed-in set.
func (s Set[M]) SupersetOf(r Set[M]) bool {
	if len(s) < len(r) {
		return false
	}

	for key := range r {
		if !s.Contains(key) {
			return false
		}
	}

	return true
}

// Subtract removes all elements from this set that are in the passed-in
// set.
//
// The set this is called on is returned, allowing for chaining.
func (s Set[M]) Subtract(r Set[M]) Set[M] {
	if s == nil {
		return s
	}
	for v := range r {
		delete(s, v)
	}
	return s
}

// Remove will remove the given value from the set if it exists.
func (s Set[M]) Remove(v M) {
	delete(s, v)
}

// Union will add all elements with the passed-in set to this set. This set
// is then returned, allowing chaining.
func (s Set[M]) Union(r Set[M]) Set[M] {
	if s == nil {
		panic("Union called on nil Set")
	}
	for v := range r {
		s[v] = void
	}
	return s
}

// XOR returns the values in either the set this method is called on, or
// the set passed in, but not both.
func (s Set[M]) XOR(r Set[M]) Set[M] {
	// this can be fixed up later to be a bit smarter about building these.
	sButNotR := s.Clone().Subtract(r)
	rButNotS := r.Clone().Subtract(s)

	return sButNotR.Union(rButNotS)
}

// Returns a new set containing only the elements
// that exist only in both sets.
func (s Set[M]) Intersect(r Set[M]) Set[M] {
	intersection := Set[M]{}
	if len(s) < len(r) {
		for elem := range s {
			if r.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range r {
			if s.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}

	return intersection
}
