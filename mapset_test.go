package cm

import (
	"reflect"
	"testing"
)

func TestMapSet(t *testing.T) {
	{
		ms := MapSet[int, int]{}
		ms.Set(1, 3)
		ms.Set(1, 4)
		ms.Set(4, 5)

		if !ms.AllValueSet().Equal(SetFromSlice[int]([]int{3, 4, 5})) {
			t.Fatal("AllValueSet didn't work")
		}
	}

	{
		ms := MapSet[int, int]{}
		ms.Set(1, 2)
		ms.Set(1, 3)

		if !ms.Contains(1, 3) {
			t.Fatal("MapSet fails to contain things")
		}
		if ms.Contains(1, 5) {
			t.Fatal("MapSet contains things it shouldn't")
		}
		if ms.Contains(4, 5) {
			t.Fatal("MapSet contains things it shouldn't")
		}

		ms.Delete(99, 99)
		ms.Delete(1, 2)
		ms.Delete(1, 3)
		if !reflect.DeepEqual(ms, MapSet[int, int]{}) {
			t.Fatal("delete doesn't clean up empty sets")
		}
	}
}

func TestNilMapSet(t *testing.T) {
	var ms MapSet[int, int]

	panics(t, "failed on MapSet.Set", func() { ms.Set(1, 2) })
	if ms.Contains(1, 2) {
		t.Fatal("nil MapSet contains something?")
	}
	ms.Delete(1, 2) // should not panic
}
