package main

import (
	"testing"
)

func TestChooseShell(t *testing.T) {

	cases := []struct {
		Name     string
		Expected string
	}{
		{"Test bash", "bash"},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := tc.Expected
			if actual != chooseShell() {
				t.Fatal("failure")
			}
		})
	}
}
