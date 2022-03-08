package cm

// A DualMap is a map that will store values in a way that allows you to
// access them by either key. For any key tuple (K1, K2), you can get the
// set of all values by either K1 or K2. This contrasts with a standard
// multilevel map, which provides no querying capability with just the
// second level of a key (other than scanning the whole thing).
//
// The full key (LK, RK) must be unique, but there can be any number of
// "lefts" associated with a given "right key" and vice versa.
//
// The zero-value of this struct is safe to use. When Set is first used,
// the maps will be initialized. NewDualMap is provided for your
// convenience if you want a DualMap with guaranteed-non-nil internal maps.
type DualMap[LK, RK comparable, V any] struct {
	Left  MultiLevelMap2[LK, RK, V]
	Right MultiLevelMap2[RK, LK, V]
}

// NewDualMap returns a new DualMap with the maps empty instead of nil.
func NewDualMap[LK, RK comparable, V any]() DualMap[LK, RK, V] {
	return DualMap[LK, RK, V]{
		MultiLevelMap2[LK, RK, V]{},
		MultiLevelMap2[RK, LK, V]{},
	}
}

func (dm *DualMap[LK, RK, V]) Set(
	l LK,
	r RK,
	value V,
) {
	if dm.Left == nil {
		dm.Left = MultiLevelMap2[LK, RK, V]{}
		dm.Right = MultiLevelMap2[RK, LK, V]{}
	}
	dm.Left.Set(l, r, value)
	dm.Right.Set(r, l, value)
}

func (dm *DualMap[LK, RK, V]) SetByTuple(
	key Tuple2[LK, RK],
	value V,
) {
	dm.Set(key.Key1, key.Key2, value)
}

func (dm *DualMap[LK, RK, V]) Get(l LK, r RK) (V, bool) {
	return dm.Left.Get(l, r)
}

func (dm *DualMap[LK, RK, V]) GetByTuple(key Tuple2[LK, RK]) (V, bool) {
	return dm.Left.GetByTuple(key)
}

func (dm *DualMap[LK, RK, V]) Delete(l LK, r RK) {
	dm.Left.Delete(l, r)
	dm.Right.Delete(r, l)
}

func (dm *DualMap[LK, RK, V]) DeleteByTuple(key Tuple2[LK, RK]) {
	dm.Delete(key.Key1, key.Key2)
}
