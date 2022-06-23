package cm

import (
	"reflect"
	"sort"
	"testing"
)

func TestMapMapMapAny(t *testing.T) {
	mlm := MapMapMapAny[int, int, int, int]{}

	mlm.Set(0, 1, 2, 3)
	mlm.SetByTuple(Tuple3[int, int, int]{4, 5, 6}, 7)

	val, exists := mlm[0][1][2]
	if val != 3 || !exists {
		t.Fatal("couldn't get a set value")
	}
	val, exists = mlm.GetByTuple(Tuple3[int, int, int]{4, 5, 6})
	if val != 7 || !exists {
		t.Fatal("couldn't get a set value")
	}
	_, exists = mlm[99][99][99]
	if exists {
		t.Fatal("can get things that don't exist")
	}

	keySlice := mlm.KeySlice()
	sort.Slice(keySlice, func(i, j int) bool {
		if keySlice[i].Key1 < keySlice[j].Key1 {
			return true
		}
		if keySlice[i].Key1 > keySlice[j].Key1 {
			return false
		}
		if keySlice[i].Key2 < keySlice[j].Key2 {
			return true
		}
		if keySlice[i].Key2 > keySlice[j].Key2 {
			return false
		}
		return keySlice[i].Key3 < keySlice[j].Key3
	})

	if !reflect.DeepEqual(keySlice,
		[]Tuple3[int, int, int]{
			{0, 1, 2},
			{4, 5, 6},
		},
	) {
		t.Fatal("incorrect key slice")
	}

	values := mlm.Values()
	sort.Ints(values)
	if !reflect.DeepEqual(values, []int{3, 7}) {
		t.Fatal("incorrect values for .Value")
	}

	keyTree := mlm.KeyTree()
	if keyTree[0].Key == 4 {
		keyTree[0], keyTree[1] = keyTree[1], keyTree[0]
	}

	if !reflect.DeepEqual(keyTree,
		[]KeyTree[int, KeyTree[int, int]]{
			{
				Key: 0,
				Vals: []KeyTree[int, int]{
					{
						Key:  1,
						Vals: []int{2},
					},
				},
			},
			{
				Key: 4,
				Vals: []KeyTree[int, int]{
					{
						Key:  5,
						Vals: []int{6},
					},
				},
			},
		},
	) {
		t.Fatal("incorrect key tree")
	}

	mlm.Delete(0, 1, 2)
	if !reflect.DeepEqual(mlm,
		MapMapMapAny[int, int, int, int]{
			4: MapMapAny[int, int, int]{
				5: map[int]int{6: 7},
			},
		},
	) {
		t.Fatal("incorrect after deleting")
	}

	mlm.Delete(4, 5, 99)
	mlm.DeleteByTuple(Tuple3[int, int, int]{4, 99, 99})
	mlm.DeleteByTuple(Tuple3[int, int, int]{99, 99, 99})
	if !reflect.DeepEqual(mlm,
		MapMapMapAny[int, int, int, int]{
			4: MapMapAny[int, int, int]{
				5: map[int]int{6: 7},
			},
		},
	) {
		t.Fatal("incorrect after deleting things that don't exist")
	}

	mlm.Delete(4, 5, 6)
	if len(mlm) != 0 {
		t.Fatal("failure to cleanup properly")
	}
}

func TestMapMapMapCloneAndEqual(t *testing.T) {
	{
		mm1 := MapMapMapAny[int, int, int, int]{}
		mm1.Set(0, 1, 2, 3)

		mm2 := mm1.Clone()

		if !mm1.EqualFunc(mm2, func(i, j int) bool { return i == j }) {
			t.Fatal("cloned maps aren't equal")
		}

		// established both the equalFunc worsk, and that the
		// cloned maps are indeed independent.
		mm2.Set(1, 2, 3, 4)
		if mm1.EqualFunc(mm2, func(i, j int) bool { return i == j }) {
			t.Fatal("different maps aren't equal")
		}
	}

}

func TestMapMapMapEqualFunc(t *testing.T) {
	var mm1 MapMapMapAny[int, int, int, int]
	var mm2 MapMapMapAny[int, int, int, int]

	eqFunc := func(i, j int) bool { return i == j }

	if !mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("two nil maps compare unequal with EqualFunc")
	}
	mm1 = MapMapMapAny[int, int, int, int]{}
	if !mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("two zero-length maps compare unequal with EqualFunc")
	}

	mm1.Set(0, 1, 2, 3)
	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("different maps compare equal with EqualFunc")
	}
	mm2 = MapMapMapAny[int, int, int, int]{}
	mm2.Set(3, 4, 5, 6)

	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("unequal maps compare equal with EqualFunc")
	}
	mm2 = MapMapMapAny[int, int, int, int]{}
	mm2.Set(0, 1, 5, 6)
	if mm1.EqualFunc(mm2, eqFunc) {
		t.Fatal("unequal maps compare equal with EqualFunc")
	}
}

func TestMapMapMapDeleteFunc(t *testing.T) {
	var mmm MapMapMapAny[int, int, int, int]

	delZero := func(k1, k2, k3, v int) bool {
		return k1 == 0
	}

	// ensure this doesn't panic, as deleting from a nil map doesn't
	mmm.DeleteFunc(delZero)

	mmm = MapMapMapAny[int, int, int, int]{}

	mmm.Set(0, 1, 2, 3)
	mmm.Set(0, 1, 3, 4)
	mmm.Set(1, 3, 4, 5)

	mmm.DeleteFunc(delZero)

	mmmTarget := MapMapMapAny[int, int, int, int]{}
	mmmTarget.Set(1, 3, 4, 5)

	if !reflect.DeepEqual(mmm, mmmTarget) {
		t.Fatal("DeleteFunc did not operate correctly")
	}
}

func TestMapMapMapEqual(t *testing.T) {
	var mmm1 MapMapMap[int, int, int, int]
	var mmm2 MapMapMap[int, int, int, int]

	if !mmm1.Equal(mmm2) {
		t.Fatal("two nil maps compare unequal with Equal")
	}
	mmm1 = MapMapMap[int, int, int, int]{}
	if !mmm1.Equal(mmm2) {
		t.Fatal("two zero-length maps compare as not equal")
	}

	mmm1.Set(0, 1, 2, 3)
	if mmm1.Equal(mmm2) {
		t.Fatal("different maps compare equal with Equal")
	}
	mmm2 = MapMapMap[int, int, int, int]{}
	mmm2.Set(3, 4, 5, 6)

	if mmm1.Equal(mmm2) {
		t.Fatal("unequal maps compare equal with Equal")
	}
	mmm2 = MapMapMap[int, int, int, int]{}
	mmm2.Set(0, 1, 5, 6)
	if mmm1.Equal(mmm2) {
		t.Fatal("unequal maps compare equal with Equal")
	}
}

// This heavily leans on TestMapMapAny to cover the functionality.
func TestMapMapMap(t *testing.T) {
	mmm := MapMapMap[int, int, int, int]{}

	mmm.Set(0, 1, 2, 3)

	mmmClone := mmm.Clone()
	if !mmmClone.Equal(mmm) {
		t.Fatal("equal doesn't seem to work")
	}

	mmm.SetByTuple(Tuple3[int, int, int]{1, 2, 3}, 4)
	mmm.Set(2, 3, 4, 5)
	mmm.Set(3, 4, 5, 6)

	mmm.Delete(0, 1, 2)
	mmm.DeleteByTuple(Tuple3[int, int, int]{1, 2, 3})
	mmm.DeleteFunc(func(int, int, int, int) bool { return false })

	val, exists := mmm[2][3][4]
	if val != 5 || !exists {
		t.Fatal("Get is wrong")
	}
	val, exists = mmm.GetByTuple(Tuple3[int, int, int]{2, 3, 4})
	if val != 5 || !exists {
		t.Fatal("GetByTuple is wrong")
	}

	mmmClone = mmm.Clone()
	if !mmm.EqualFunc(mmmClone, func(x, y int) bool { return x == y }) {
		t.Fatal("EqualFunc doesn't work.")
	}

	mmm.KeySlice()
	mmm.KeyTree()
}
