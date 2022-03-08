package cm

import (
	"reflect"
	"sort"
	"testing"
)

func TestMLM3(t *testing.T) {
	mlm := MultiLevelMap3[int, int, int, int]{}

	mlm.Set(0, 1, 2, 3)
	mlm.SetByTuple(Tuple3[int, int, int]{4, 5, 6}, 7)

	val, exists := mlm.Get(0, 1, 2)
	if val != 3 || !exists {
		t.Fatal("couldn't get a set value")
	}
	val, exists = mlm.GetByTuple(Tuple3[int, int, int]{4, 5, 6})
	if val != 7 || !exists {
		t.Fatal("couldn't get a set value")
	}
	_, exists = mlm.Get(99, 99, 99)
	if exists {
		t.Fatal("can get things that don't exist")
	}

	thirdLevel := mlm.GetThirdLevel(0, 1)
	if !reflect.DeepEqual(thirdLevel, map[int]int{2: 3}) {
		t.Fatal("couldn't fetch third level")
	}
	thirdLevel = mlm.GetThirdLevel(99, 99)
	if thirdLevel != nil {
		t.Fatal("can get third levels of things that don't exist")
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

	keyTree := mlm.KeyTree()

	if !reflect.DeepEqual(keyTree,
		[]KeyTree3[int, int, int]{
			{
				Key1: 0,
				Value: []KeyTree2[int, int]{
					{
						Key1:  1,
						Key2s: []int{2},
					},
				},
			},
			{
				Key1: 4,
				Value: []KeyTree2[int, int]{
					{
						Key1:  5,
						Key2s: []int{6},
					},
				},
			},
		},
	) {
		t.Fatal("incorrect key tree")
	}

	mlm.Delete(0, 1, 2)
	if !reflect.DeepEqual(mlm,
		MultiLevelMap3[int, int, int, int]{
			4: MultiLevelMap2[int, int, int]{
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
		MultiLevelMap3[int, int, int, int]{
			4: MultiLevelMap2[int, int, int]{
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
