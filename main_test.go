package main

import (
	"testing"
)

func TestExecutution(t *testing.T) {
	if 2 != 2 {
		t.Error("Expected 2 != 4")
	}
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"Test fist item", 1, 1, 2},
		{"Test second item", 2, 4, 6},
		{"Another one", 2, 8, 10},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := tc.A + tc.B
			if actual != tc.A+tc.B {
				t.Fatal("failure")
			}
		})
	}
}
