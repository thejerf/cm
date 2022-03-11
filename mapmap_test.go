package cm

import (
	"reflect"
	"sort"
	"testing"
)

func TestMapMapAny(t *testing.T) {
	mm := MapMapAny[int, int, int]{}

	mm.Set(0, 1, 2)

	if !reflect.DeepEqual(mm,
		MapMapAny[int, int, int]{
			0: map[int]int{
				1: 2,
			},
		},
	) {
		t.Fatal("did not correctly set the values")
	}

	mm.SetByTuple(Tuple2[int, int]{3, 4}, 5)
	if !reflect.DeepEqual(mm,
		MapMapAny[int, int, int]{
			0: map[int]int{
				1: 2,
			},
			3: map[int]int{
				4: 5,
			},
		},
	) {
		t.Fatal("did not correctly set the values")
	}

	mm.Set(0, 1, 6)
	mm.Set(0, 7, 8)
	if !reflect.DeepEqual(mm,
		MapMapAny[int, int, int]{
			0: map[int]int{
				1: 6,
				7: 8,
			},
			3: map[int]int{
				4: 5,
			},
		},
	) {
		t.Fatal("did not correctly set the values")
	}

	v, exists := mm.Get(0, 0)
	if exists || v != 0 {
		t.Fatal("could fetch something that didn't exist!")
	}
	v, exists = mm.Get(-1, -1)
	if exists || v != 0 {
		t.Fatal("could fetch something that didn't exist!")
	}
	v, exists = mm.Get(0, 1)
	if !exists || v != 6 {
		t.Fatal("couldn't fetch something that exists")
	}
	v, exists = mm.GetByTuple(Tuple2[int, int]{0, 1})
	if !exists || v != 6 {
		t.Fatal("couldn't fetch something that exists")
	}

	keySlice := mm.KeySlice()
	sort.Slice(keySlice, func(i, j int) bool {
		if keySlice[i].Key1 < keySlice[j].Key1 {
			return true
		}
		if keySlice[i].Key1 > keySlice[j].Key1 {
			return false
		}
		return keySlice[i].Key2 < keySlice[j].Key2
	})
	if !reflect.DeepEqual(keySlice,
		[]Tuple2[int, int]{
			{0, 1},
			{0, 7},
			{3, 4},
		},
	) {
		t.Fatal("incorrect keys coming out")
	}

	keyTree := mm.KeyTree()
	if len(keyTree) != 2 {
		t.Fatal("incorrect key tree")
	}
	if keyTree[0].Key == 3 {
		keyTree[0], keyTree[1] = keyTree[1], keyTree[0]
	}
	if keyTree[0].Vals[0] == 7 {
		keyTree[0].Vals[0], keyTree[0].Vals[1] =
			keyTree[0].Vals[1], keyTree[0].Vals[0]
	}
	if !reflect.DeepEqual(keyTree,
		[]KeyTree[int, int]{
			{0, []int{1, 7}},
			{3, []int{4}},
		},
	) {
		t.Fatal("Incorrect key tree")
	}

	// And now, clean things up
	mm.Delete(3, 4)
	if !reflect.DeepEqual(mm,
		MapMapAny[int, int, int]{
			0: map[int]int{
				1: 6,
				7: 8,
			},
		},
	) {
		t.Fatal("did not correctly delete the values")
	}

	mm.DeleteByTuple(Tuple2[int, int]{0, 7})
	if !reflect.DeepEqual(mm,
		MapMapAny[int, int, int]{
			0: map[int]int{
				1: 6,
			},
		},
	) {
		t.Fatal("did not correctly delete the values")
	}

	mm.Delete(-1, -1)
	mm.Delete(0, 1)
	if !reflect.DeepEqual(mm, MapMapAny[int, int, int]{}) {
		t.Fatal("didn't clear out values")
	}
}

func TestMapMapCloneAndEqual(t *testing.T) {
	{
		mm1 := MapMapAny[int, int, int]{}
		mm1.Set(0, 1, 2)

		mm2 := mm1.Clone()

		if !(MapMap[int, int, int](mm1).Equal(MapMap[int, int, int](mm2))) {
			t.Fatal("two equal maps aren't equal")
		}
		mm2.Set(0, 3, 4)
		// This both established that they are correctly not equal, and
		// that the clone really is independent
		if MapMap[int, int, int](mm1).Equal(MapMap[int, int, int](mm2)) {
			t.Fatal("two unequal maps aren't unequal")
		}
	}

	{
		// Coverage; test the various quick-escapes for .Equals.
		var cmm1 MapMap[int, int, int]
		var cmm2 MapMap[int, int, int]
		if !cmm1.Equal(cmm2) {
			t.Fatal("two nil CompMapMaps not equal")
		}

		mm1 := MapMapAny[int, int, int]{}
		mm1.Set(0, 1, 2)

		cmm1 = MapMap[int, int, int](mm1)

		if cmm1.Equal(cmm2) {
			t.Fatal("unequal compare the same")
		}

		cmm1 = MapMap[int, int, int]{}
		if !cmm1.Equal(cmm2) {
			t.Fatal("nil and empty maps don't say equal")
		}

		cmm1 = MapMap[int, int, int](mm1)
		mm2 := MapMapAny[int, int, int]{}
		mm2.Set(1, 2, 3)

		if cmm1.Equal(MapMap[int, int, int](mm2)) {
			t.Fatal("unequal maps compare equal")
		}
	}

	strLenEq := func(a, b string) bool {
		return len(a) == len(b)
	}

	{
		// Test that nil and the zero map are equal with EqualFunc
		var mms1 MapMapAny[int, int, string]
		mms2 := MapMapAny[int, int, string]{}
		if !mms1.EqualFunc(mms2, strLenEq) || !mms2.EqualFunc(mms1, strLenEq) {
			t.Fatal("nil != zero MapMap")
		}
	}

	{
		// Test EqualsFunc with something that calls strings equal if they
		// have the same length.
		mms1 := MapMapAny[int, int, string]{}
		mms2 := MapMapAny[int, int, string]{}
		mms1.Set(0, 1, "one")
		mms2.Set(0, 1, "two")

		if !mms1.EqualFunc(mms2, strLenEq) {
			t.Fatal("EqualFunc doesn't work")
		}
	}
}

func TestMapMapEqualFunc(t *testing.T) {
	var mm1 MapMapAny[int, int, int]
	var mm2 MapMapAny[int, int, int]

	eqFunc := func(i, j int) bool { return i == j }

	if !mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("two nil maps compare unequal with EqualFunc")
	}
	mm1 = MapMapAny[int, int, int]{}
	mm1.Set(0, 1, 2)
	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("different maps compare equal with EqualFunc")
	}
	mm2 = MapMapAny[int, int, int]{}
	mm2.Set(3, 4, 5)

	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("unequal maps compare equal with EqualFunc")
	}
	mm2 = MapMapAny[int, int, int]{}
	mm2.Set(0, 1, 5)
	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("unequal maps compare equal with EqualFunc")
	}
}

func TestMapMapDeleteFunc(t *testing.T) {
	var mm MapMapAny[int, int, int]

	// Delete zeros in the first key slot
	delZero := func(k1, k2, v int) bool {
		return k1 == 0
	}

	// ensure this doesn't panic, as deleting from a nil map doesn't
	mm.DeleteFunc(delZero)

	mm = MapMapAny[int, int, int]{}

	mm.Set(0, 1, 2)
	mm.Set(0, 2, 4)
	mm.Set(1, 4, 5)

	mm.DeleteFunc(delZero)

	mmTarget := MapMapAny[int, int, int]{}
	mmTarget.Set(1, 4, 5)

	if !reflect.DeepEqual(mm, mmTarget) {
		t.Fatal("DeleteFunc did not operate correctly")
	}
}

func TestNilMapMap(t *testing.T) {
	var mm MapMapAny[int, int, int]

	panics(t, "failed on set", func() { mm.Set(0, 1, 2) })
	panics(t, "failed on set by tuple",
		func() { mm.SetByTuple(Tuple2[int, int]{0, 1}, 2) })

	// doesn't panic
	mm.Delete(0, 0)
	mm.DeleteByTuple(Tuple2[int, int]{0, 0})

	panics(t, "failed on Get", func() { mm.Get(0, 0) })
	panics(t, "failed on GetByTuple", func() { mm.GetByTuple(Tuple2[int, int]{0, 0}) })

	if mm.KeySlice() != nil {
		t.Fatal("incorrect KeySlice from nil map")
	}
	if mm.KeyTree() != nil {
		t.Fatal("incorrect KeyTree from nil map")
	}
}

// This heavily leans on TestMapMapAny to cover the functionality.
func TestMapMap(t *testing.T) {
	mm := MapMap[int, int, int]{}

	mm.Set(0, 1, 2)

	mmClone := mm.Clone()
	if !mmClone.Equal(mm) {
		t.Fatal("equal doesn't seem to work")
	}

	mm.SetByTuple(Tuple2[int, int]{1, 2}, 3)
	mm.Set(2, 3, 4)
	mm.Set(3, 4, 5)

	mm.Delete(0, 1)
	mm.DeleteByTuple(Tuple2[int, int]{1, 2})
	mm.DeleteFunc(func(int, int, int) bool { return false })

	val, exists := mm.Get(2, 3)
	if val != 4 || !exists {
		t.Fatal("Get is wrong")
	}
	val, exists = mm.GetByTuple(Tuple2[int, int]{2, 3})
	if val != 4 || !exists {
		t.Fatal("GetByTuple is wrong")
	}

	mmClone = mm.Clone()
	if !mm.EqualFunc(mmClone, func(x, y int) bool { return x == y }) {
		t.Fatal("EqualFunc doesn't work.")
	}

	mm.KeySlice()
	mm.KeyTree()
}

func panics(t *testing.T, message string, f func()) {
	panicked := false
	defer func() {
		r := recover()
		if r != nil {
			panicked = true
		}

		if !panicked {
			t.Fatal(message)
		}
	}()
	f()
}
