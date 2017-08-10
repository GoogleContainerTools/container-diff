package utils

import (
	"testing"
)

func TestGetImagePullResponse(t *testing.T) {
	for _, test := range []struct {
		image          string
		response       []Event
		expectedOutput string
		expectedError  bool
	}{
		{
			image:          "noimage",
			response:       []Event{},
			expectedOutput: "Could not pull image noimage",
			expectedError:  true,
		},
		{
			image:          "gcr.io/google_containers/nonexistentimage",
			response:       []Event{{Error: "Non-existing image"}},
			expectedOutput: "Error pulling image gcr.io/google_containers/nonexistentimage: Non-existing image",
			expectedError:  true,
		},
		{
			image:          "gcr.io/google_containers/existentimage",
			response:       []Event{{Status: "Digest: sha256:c34ce3c1fcc0c7431e1392cc3abd0dfe2192ffea1898d5250f199d3ac8d8720f"}},
			expectedOutput: "sha256:c34ce3c1fcc0c7431e1392cc3abd0dfe2192ffea1898d5250f199d3ac8d8720f",
			expectedError:  false,
		},
	} {
		output, err := getImagePullResponse(test.image, test.response)
		if err != nil && !test.expectedError {
			t.Errorf("Got unexpected error: %s", err)
		} else if err == nil && test.expectedError {
			t.Error("Expected error but got none")
		} else if err != nil && (test.expectedOutput != err.Error()) {
			t.Error("Had trouble getting error")
		} else if err == nil && test.expectedOutput != output {
			t.Errorf("Expected %s but got %s", test.expectedOutput, output)
		}
	}
}
