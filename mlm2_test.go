package cm

import (
	"reflect"
	"sort"
	"testing"
)

func TestMLM2(t *testing.T) {
	mlm := MultiLevelMap2[int, int, int]{}

	mlm.Set(0, 1, 2)

	if !reflect.DeepEqual(mlm,
		MultiLevelMap2[int, int, int]{
			0: map[int]int{
				1: 2,
			},
		},
	) {
		t.Fatal("did not correctly set the values")
	}

	mlm.SetByTuple(Tuple2[int, int]{3, 4}, 5)
	if !reflect.DeepEqual(mlm,
		MultiLevelMap2[int, int, int]{
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

	mlm.Set(0, 1, 6)
	mlm.Set(0, 7, 8)
	if !reflect.DeepEqual(mlm,
		MultiLevelMap2[int, int, int]{
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

	v, exists := mlm.Get(0, 0)
	if exists || v != 0 {
		t.Fatal("could fetch something that didn't exist!")
	}
	v, exists = mlm.Get(-1, -1)
	if exists || v != 0 {
		t.Fatal("could fetch something that didn't exist!")
	}
	v, exists = mlm.Get(0, 1)
	if !exists || v != 6 {
		t.Fatal("couldn't fetch something that exists")
	}
	v, exists = mlm.GetByTuple(Tuple2[int, int]{0, 1})
	if !exists || v != 6 {
		t.Fatal("couldn't fetch something that exists")
	}

	keySlice := mlm.KeySlice()
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

	keyTree := mlm.KeyTree()
	if len(keyTree) != 2 {
		t.Fatal("incorrect key tree")
	}
	if keyTree[0].Key1 == 3 {
		keyTree[0], keyTree[1] = keyTree[1], keyTree[0]
	}
	if keyTree[0].Key2s[0] == 7 {
		keyTree[0].Key2s[0], keyTree[0].Key2s[1] =
			keyTree[0].Key2s[1], keyTree[0].Key2s[0]
	}
	if !reflect.DeepEqual(keyTree,
		[]KeyTree2[int, int]{
			{0, []int{1, 7}},
			{3, []int{4}},
		},
	) {
		t.Fatal("Incorrect key tree")
	}

	// And now, clean things up
	mlm.Delete(3, 4)
	if !reflect.DeepEqual(mlm,
		MultiLevelMap2[int, int, int]{
			0: map[int]int{
				1: 6,
				7: 8,
			},
		},
	) {
		t.Fatal("did not correctly delete the values")
	}

	mlm.DeleteByTuple(Tuple2[int, int]{0, 7})
	if !reflect.DeepEqual(mlm,
		MultiLevelMap2[int, int, int]{
			0: map[int]int{
				1: 6,
			},
		},
	) {
		t.Fatal("did not correctly delete the values")
	}

	mlm.Delete(-1, -1)
	mlm.Delete(0, 1)
	if !reflect.DeepEqual(mlm, MultiLevelMap2[int, int, int]{}) {
		t.Fatal("didn't clear out values")
	}
}

func TestNilMLM2(t *testing.T) {
	var mlm MultiLevelMap2[int, int, int]

	panics(t, "failed on set", func() { mlm.Set(0, 1, 2) })
	panics(t, "failed on set by tuple",
		func() { mlm.SetByTuple(Tuple2[int, int]{0, 1}, 2) })

	// doesn't panic
	mlm.Delete(0, 0)
	mlm.DeleteByTuple(Tuple2[int, int]{0, 0})

	_, exists := mlm.Get(0, 0)
	if exists {
		t.Fatal("was able to retrieve something from nil map")
	}
	_, exists = mlm.GetByTuple(Tuple2[int, int]{0, 0})
	if exists {
		t.Fatal("was able to retrieve something from nil map")
	}

	if mlm.KeySlice() != nil {
		t.Fatal("incorrect KeySlice from nil map")
	}
	if mlm.KeyTree() != nil {
		t.Fatal("incorrect KeyTree from nil map")
	}
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
