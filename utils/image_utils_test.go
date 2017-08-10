package utils

import (
	"testing"
)

type imageTestPair struct {
	input          string
	expectedOutput bool
}

func TestCheckImageID(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "123456789012", expectedOutput: true},
		{input: "gcr.io/repo/image", expectedOutput: false},
		{input: "testTars/la-croix1.tar", expectedOutput: false},
	} {
		output := CheckImageID(test.input)
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be image ID but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be an image ID but %s tested true", test.input)
			}
		}
	}
}

func TestCheckImageTar(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "123456789012", expectedOutput: false},
		{input: "gcr.io/repo/image", expectedOutput: false},
		{input: "testTars/la-croix1.tar", expectedOutput: true},
	} {
		output := CheckTar(test.input)
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be a tar file but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be a tar file but %s tested true", test.input)
			}
		}
	}
}

func TestCheckImageURL(t *testing.T) {
	for _, test := range []imageTestPair{
		{input: "123456789012", expectedOutput: false},
		{input: "gcr.io/repo/image", expectedOutput: true},
		{input: "testTars/la-croix1.tar", expectedOutput: false},
	} {
		output := CheckImageURL(test.input)
		if output != test.expectedOutput {
			if test.expectedOutput {
				t.Errorf("Expected input to be a tar file but %s tested false", test.input)
			} else {
				t.Errorf("Didn't expect input to be a tar file but %s tested true", test.input)
			}
		}
	}
}
