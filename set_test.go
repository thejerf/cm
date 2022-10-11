package cm

import (
	"reflect"
	"sort"
	"testing"
)

func TestSetConstruct(t *testing.T) {
	s := SetFromSlice([]int{1, 2, 3})
	s.Add(4)

	if !s.Contains(1) {
		t.Fatal("set doesn't contain properly")
	}
	if s.Contains(99) {
		t.Fatal("set doesn't contain properly")
	}

	s2 := s.Clone()
	if !s.Equal(s2) {
		t.Fatal("clone didn't clone correctly")
	}
	s2.Remove(1)
	if s.Equal(s2) {
		t.Fatal("clone is not correctly independent")
	}

	s2Vals := s2.AsSlice()
	sort.Ints(s2Vals)
	if !reflect.DeepEqual(s2Vals, []int{2, 3, 4}) {
		t.Fatal("AsSlice doesn't work")
	}
}

func TestSetEquality(t *testing.T) {
	var s1 Set[int]
	var s2 Set[int]

	if !s1.Equal(s2) {
		t.Fatal("equal doesn't work on nil Sets")
	}
	s1 = Set[int]{}
	if !s1.Equal(s2) {
		t.Fatal("equal doesn't work on zero-length sets")
	}

	s2 = Set[int]{}

	s1.Add(1)
	s2.Add(2)

	if s1.Equal(s2) {
		t.Fatal("unequal sets are equal")
	}
}

func TestSetSubset(t *testing.T) {
	s1 := SetFromSlice([]int{1, 2, 3})
	s2 := SetFromSlice([]int{1, 2})

	if s1.SubsetOf(s2) {
		t.Fatal("subset fails")
	}
	if !s2.SubsetOf(s1) {
		t.Fatal("subset fails")
	}
	if !s1.SubsetOf(s1) {
		t.Fatal("subset really fails")
	}

	if !s1.SupersetOf(s2) {
		t.Fatal("superset fails")
	}
	if s2.SupersetOf(s1) {
		t.Fatal("superset fails")
	}
	if !s1.SupersetOf(s1) {
		t.Fatal("superset really fails")
	}

	s1 = SetFromSlice([]int{1})
	s2 = SetFromSlice([]int{2})

	if s1.SubsetOf(s2) {
		t.Fatal("subset fails")
	}
	if s1.SupersetOf(s2) {
		t.Fatal("superset fails")
	}
}

func TestSetIntersect(t *testing.T) {
	var s1 = Set[int]{}
	var s2 = Set[int]{}

	s1.Add(1)
	s1.Add(2)
	s1.Add(3)
	s2.Add(2)
	s2.Add(4)
	s2.Add(6)

	s3 := s1.Intersect(s2)

	if !s3.Contains(2) {
		t.Fatal("set doesn't intersect properly")
	}
}

func TestSetOperations(t *testing.T) {
	s1 := SetFromSlice([]int{1, 2, 3})

	s1.Subtract(SetFromSlice([]int{2, 3}))

	if !s1.Equal(SetFromSlice([]int{1})) {
		t.Fatal("subtract doesn't work")
	}

	s1.Union(SetFromSlice([]int{2}))
	if !s1.Equal(SetFromSlice([]int{1, 2})) {
		t.Fatal("subtract doesn't work")
	}

	s1 = SetFromSlice([]int{1, 2})
	s2 := SetFromSlice([]int{2, 3})

	if !(s1.XOR(s2).Equal(SetFromSlice([]int{1, 3}))) {
		t.Fatal("XOR doesn't work")
	}
}

func TestNilSet(t *testing.T) {
	var s Set[int]

	panics(t, "failed on Set.Add", func() { s.Add(1) })
	panics(t, "failed on Set.Union", func() { s.Union(s) })

	if s.Clone() != nil {
		t.Fatal("nil set doesn't clone as a nil set")
	}
	if s.Contains(1) {
		t.Fatal("nil set contains things")
	}
	if s.Subtract(SetFromSlice([]int{1})) != nil {
		t.Fatal("subtract from a nil set is not nil")
	}
}
