package admin

import (
	"reflect"
	"testing"
)

func TestNormalizePeerIDs(t *testing.T) {
	ids := []uint{12, 8, 0, 8, 42, 12, 99, 42}

	got := NormalizePeerIDs(ids)
	want := []uint{12, 8, 42, 99}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizePeerIDs() = %#v, want %#v", got, want)
	}
}

func TestNormalizePeerIDsEmpty(t *testing.T) {
	if got := NormalizePeerIDs(nil); len(got) != 0 {
		t.Fatalf("NormalizePeerIDs(nil) = %#v, want empty slice", got)
	}
}
