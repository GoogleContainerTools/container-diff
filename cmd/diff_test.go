package cmd

import (
	"testing"
)

var diffArgNumTests = []testpair{
	{[]string{}, false},
	{[]string{"one"}, false},
	{[]string{"one", "two"}, true},
	{[]string{"one", "two", "three"}, false},
}

func TestDiffArgNum(t *testing.T) {
	for _, test := range diffArgNumTests {
		valid, err := checkDiffArgNum(test.input)
		if valid != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
