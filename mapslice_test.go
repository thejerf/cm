package cm

import (
	"reflect"
	"testing"
)

func TestMapSlice(t *testing.T) {
	var ms MapSlice[int, int]

	panics(t, "failed on MapSlice.Append",
		func() { ms.Append(1, 2) })
	panics(t, "failed on MapSlice.Set",
		func() { ms.Set(1, []int{2}) })

	ms = MapSlice[int, int]{}

	ms.Append(1, 2)
	if !reflect.DeepEqual(ms,
		MapSlice[int, int]{
			1: {2},
		},
	) {
		t.Fatal("Append doesn't work on unset values")
	}
	ms.Append(1, 2)
	if !reflect.DeepEqual(ms,
		MapSlice[int, int]{
			1: {2, 2},
		},
	) {
		t.Fatal("Append doesn't work on set values")
	}

	ms.Set(1, []int{3})
	if !reflect.DeepEqual(ms,
		MapSlice[int, int]{
			1: {3},
		},
	) {
		t.Fatal("Append doesn't work on unset values")
	}

	ms.Set(1, []int{})
	if !reflect.DeepEqual(ms, MapSlice[int, int]{}) {
		t.Fatal("set doesn't remove empty slices")
	}
}
