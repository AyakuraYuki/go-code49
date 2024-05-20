package code49

import (
	"slices"
	"testing"
)

func Test_toBytes(t *testing.T) {
	tests := []struct {
		Val  string
		Want []int
	}{
		{"hello world", []int{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}},
		{"hello\\<FNC1>world", []int{104, 101, 108, 108, 111, -1, 119, 111, 114, 108, 100}},
		{"hello\\<FNC2>world", []int{104, 101, 108, 108, 111, -2, 119, 111, 114, 108, 100}},
		{"hello\\<FNC3>world", []int{104, 101, 108, 108, 111, -3, 119, 111, 114, 108, 100}},
	}

	for _, tt := range tests {
		get := toBytes(tt.Val)
		if !slices.Equal(get, tt.Want) {
			t.Fatalf("error getting bytes with %q, got %v but want %v", tt.Val, get, tt.Want)
		}
	}
}
