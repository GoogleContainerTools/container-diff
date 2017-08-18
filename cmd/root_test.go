package cmd

import (
	"testing"
)

type testpair struct {
	input           []string
	expected_output bool
}

var argNumTests = []testpair{
	{[]string{}, false},
	{[]string{"one"}, true},
	{[]string{"one", "two"}, true},
	{[]string{"one", "two", "three"}, false},
}

var argTypeTests = []testpair{
	{[]string{"badID", "badID"}, false},
	{[]string{"123456789012", "badID"}, false},
	{[]string{"123456789012", "123456789012"}, true},
	{[]string{"?!badDiffer71", "123456789012"}, false},
	{[]string{"123456789012", "gcr.io/repo/image"}, true},
}

func TestArgNum(t *testing.T) {
	for _, test := range argNumTests {
		valid, err := checkArgNum(test.input)
		if valid != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}

func TestArgType(t *testing.T) {
	for _, test := range argTypeTests {
		valid, err := checkArgType(test.input)
		if valid != test.expected_output {
			if test.expected_output {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				t.Errorf("Expected error but got none")
			}
		}
	}
}
