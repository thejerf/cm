package cm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDualMaps(t *testing.T) {
	dm := NewDualMap[int, string, int]()

	dm.Set(0, "a", 100)

	expected := DualMap[int, string, int]{
		MultiLevelMap2[int, string, int]{
			0: map[string]int{
				"a": 100,
			},
		},
		MultiLevelMap2[string, int, int]{
			"a": map[int]int{
				0: 100,
			},
		},
	}
	if !reflect.DeepEqual(dm, expected) {
		t.Fatal("incorrectly-set dual map")
	}

	dm.SetByTuple(Tuple2[int, string]{0, "b"}, 200)

	v, exists := dm.Get(0, "a")
	if v != 100 || !exists {
		t.Fatal("couldn't get value")
	}
	v, exists = dm.GetByTuple(Tuple2[int, string]{0, "b"})
	if v != 200 || !exists {
		t.Fatal("couldn't get value")
	}

	dm.Delete(0, "a")
	_, exists = dm.Get(0, "a")
	if exists {
		t.Fatal("can fetch deleted values")
	}
	v, exists = dm.Get(0, "b")
	if v != 200 || !exists {
		t.Fatal("over-deleted value")
	}

	dm.DeleteByTuple(Tuple2[int, string]{1, "hi"})

	dm = DualMap[int, string, int]{}
	dm.Set(0, "a", 300)
	v, exists = dm.Get(0, "a")
	if v != 300 || !exists {
		fmt.Println(v, exists)
		t.Fatal("couldn't set into zero-value DualMap")
	}
}
