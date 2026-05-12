package tests

import (
	"reflect"
	"testing"
	"tetris-optimizer/libs"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{10, 10},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tt := range tests {
		result := libs.Abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    []libs.Point
		expected []libs.Point
	}{
		{
			[]libs.Point{{2, 3}, {3, 3}, {2, 4}, {3, 4}},
			[]libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		},
		{
			[]libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
			[]libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, // already normalized
		},
		{
			[]libs.Point{{5, 5}, {6, 5}, {5, 6}, {6, 6}},
			[]libs.Point{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		},
	}

	for _, tt := range tests {
		result := libs.Normalize(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("Normalize(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
