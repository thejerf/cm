package cm

// A DualMap is a map that will store values in a way that allows you to
// access them by either key. For any key tuple (K1, K2), you can get the
// set of all values by either K1 or K2. This contrasts with a standard
// multilevel map, which provides no querying capability with just the
// second level of a key (other than scanning the whole thing).
//
// The full key (P, S) must be unique, but there can be any number of
// "primaries" associated with a given "secondary key" and vice versa.
//
// I have found it convenient to remember the map as having one of the
// particular types as "primary", so this map refers to the "Primary"
// mapping and the "Reverse" mapping. This helps keep straight which keys
// are which, even in situations where you have no particular preference.
//
// The zero-value of this struct is safe to use. When Set is first used,
// the maps will be initialized.
//
// Direct read access is permissible. You should not directly write to the
// maps. DualMap makes no guarantees if you directly write to the internal
// maps.
type DualMap[P, S comparable, V any] struct {
	Primary  MapMapAny[P, S, V]
	Reverse MapMapAny[S, P, V]
}

// Set sets the given value with the keys in primary/secondary order.
func (dm *DualMap[P, S, V]) Set(
	l P,
	r S,
	value V,
) {
	if dm.Primary == nil {
		dm.Primary = MapMapAny[P, S, V]{}
		dm.Reverse = MapMapAny[S, P, V]{}
	}
	dm.Primary.Set(l, r, value)
	dm.Reverse.Set(r, l, value)
}

// SetByTuple sets by the tuple returned by the Primary's KeySlice method.
func (dm *DualMap[P, S, V]) SetByTuple(
	key Tuple2[P, S],
	value V,
) {
	dm.Set(key.Key1, key.Key2, value)
}

// GetByTuple retrieves by the tuple returned by the Primary's KeySlice method.
func (dm *DualMap[P, S, V]) GetByTuple(key Tuple2[P, S]) (val V, exists bool) {
	val, exists = dm.Primary[key.Key1][key.Key2]
	return val, exists
}

// Delete deletes by the keys in primary/secondary order.
func (dm *DualMap[P, S, V]) Delete(l P, r S) {
	dm.Primary.Delete(l, r)
	dm.Reverse.Delete(r, l)
}

// Delete deletes by the tuple returned by the Primary's KeySlice method.
func (dm *DualMap[P, S, V]) DeleteByTuple(key Tuple2[P, S]) {
	dm.Delete(key.Key1, key.Key2)
}
