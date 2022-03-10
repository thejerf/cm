package cm

// A MapSlice is a map that contains slices.
//
// Because the nil slice is broadly functional, there aren't that many
// methods that are helpful beyond direct access, but I find Append comes
// up a lot in my code.
type MapSlice[K comparable, V any] map[K][]V

// Append will append the given values to the slice in the map, or create
// one if needed.
func (mp MapSlice[K, V]) Append(key K, vals ...V) {
	if mp == nil {
		panic("Append called on nil MapSlice")
	}
	slice := mp[key]
	if slice == nil {
		mp[key] = vals
		return
	}
	slice = append(slice, vals...)
	mp[key] = slice
}

// Set will set the given slice to the be the value in the map, except if
// the slice is of length zero, in which case the key will be cleared.
func (mp MapSlice[K, V]) Set(key K, s []V) {
	if mp == nil {
		panic("Set called on nil MapSlice")
	}
	if len(s) == 0 {
		delete(mp, key)
		return
	}
	mp[key] = s
}
