package cm

import (
	"reflect"
	"testing"
)

func TestMapSet(t *testing.T) {
	{
		ms := MapSet[int, int]{}
		ms.Add(1, 3)
		ms.Add(1, 4)
		ms.Add(4, 5)

		if !ms.AllValueSet().Equal(SetFromSlice([]int{3, 4, 5})) {
			t.Fatal("AllValueSet didn't work")
		}
	}

	{
		ms := MapSet[int, int]{}
		ms.Add(1, 2)
		ms.Add(1, 3)

		ms.Delete(99, 99)
		ms.Delete(1, 2)
		ms.Delete(1, 3)
		if !reflect.DeepEqual(ms, MapSet[int, int]{}) {
			t.Fatal("delete doesn't clean up empty sets")
		}
	}

	{
		ms := MapSet[int, int]{}
		s1 := SetFromSlice([]int{1})
		s2 := SetFromSlice([]int{2})

		ms.Union(0, s1)
		if !ms[0].Contains(1) {
			t.Fatal("union didn't work")
		}

		ms.Union(0, s2)
		if !ms[0].Contains(2) {
			t.Fatal("union didn't work")
		}
	}
}

func TestNilMapSet(t *testing.T) {
	var ms MapSet[int, int]

	panics(t, "failed on MapSet.Add", func() { ms.Add(1, 2) })
	if ms[0].Contains(2) {
		t.Fatal("nil MapSet contains something?")
	}
	ms.Delete(1, 2) // should not panic
	panics(t, "failed on MapSet.Union",
		func() { ms.Union(0, SetFromSlice([]int{1})) })
}
