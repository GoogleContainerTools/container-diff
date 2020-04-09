// +build integration

/*
Copyright 2018 Google, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	diffBase     = "gcr.io/gcp-runtimes/container-diff-tests/diff-base"
	diffModified = "gcr.io/gcp-runtimes/container-diff-tests/diff-modified"

	diffLayerBase     = "gcr.io/gcp-runtimes/container-diff-tests/diff-layer-base"
	diffLayerModified = "gcr.io/gcp-runtimes/container-diff-tests/diff-layer-modified"

	metadataBase     = "gcr.io/gcp-runtimes/container-diff-tests/metadata-base"
	metadataModified = "gcr.io/gcp-runtimes/container-diff-tests/metadata-modified"

	aptBase     = "gcr.io/gcp-runtimes/container-diff-tests/apt-base"
	aptModified = "gcr.io/gcp-runtimes/container-diff-tests/apt-modified"

	rpmBase     = "valentinrothberg/containerdiff:diff-base"
	rpmModified = "valentinrothberg/containerdiff:diff-modified"

	// Why is this node-modified:2.0?
	nodeBase     = "gcr.io/gcp-runtimes/container-diff-tests/node-modified:2.0"
	nodeModified = "gcr.io/gcp-runtimes/container-diff-tests/node-modified"

	pipModified = "gcr.io/gcp-runtimes/container-diff-tests/pip-modified"

	multiBase     = "gcr.io/gcp-runtimes/container-diff-tests/multi-base"
	multiModified = "gcr.io/gcp-runtimes/container-diff-tests/multi-modified"

	multiBaseLocal     = "daemon://gcr.io/gcp-runtimes/container-diff-tests/multi-base"
	multiModifiedLocal = "daemon://gcr.io/gcp-runtimes/container-diff-tests/multi-modified"
)

type ContainerDiffRunner struct {
	t          *testing.T
	binaryPath string
}

func (c *ContainerDiffRunner) Run(command ...string) (string, string, error) {
	path, err := filepath.Abs(c.binaryPath)
	if err != nil {
		c.t.Fatalf("Error finding container-diff binary: %s", err)
	}
	c.t.Logf("Running command: %s %s", path, command)
	cmd := exec.Command(path, command...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("Error running command %s: %s Stderr: %s", command, err, stderr.String())
	}
	return stdout.String(), stderr.String(), nil
}

func TestDiffAndAnalysis(t *testing.T) {
	runner := ContainerDiffRunner{
		t:          t,
		binaryPath: "../out/container-diff",
	}

	var tests = []struct {
		description string
		imageA      string
		imageB      string
		differFlags []string
		subcommand  string

		//TODO: Don't consume a json file
		expectedFile string
	}{
		{
			description:  "file differ",
			subcommand:   "diff",
			imageA:       diffBase,
			imageB:       diffModified,
			differFlags:  []string{"--type=file", "--no-cache"},
			expectedFile: "file_diff_expected.json",
		},
		{
			description:  "file layer differ",
			subcommand:   "diff",
			imageA:       diffLayerBase,
			imageB:       diffLayerModified,
			differFlags:  []string{"--type=layer", "--no-cache"},
			expectedFile: "file_layer_diff_expected.json",
		},
		{
			description:  "size differ",
			subcommand:   "diff",
			imageA:       diffLayerBase,
			imageB:       diffLayerModified,
			differFlags:  []string{"--type=size", "--no-cache"},
			expectedFile: "size_diff_expected.json",
		},
		{
			description:  "size layer differ",
			subcommand:   "diff",
			imageA:       diffLayerBase,
			imageB:       diffLayerModified,
			differFlags:  []string{"--type=sizelayer", "--no-cache"},
			expectedFile: "size_layer_diff_expected.json",
		},
		{
			description:  "apt differ",
			subcommand:   "diff",
			imageA:       aptBase,
			imageB:       aptModified,
			differFlags:  []string{"--type=apt", "--no-cache"},
			expectedFile: "apt_diff_expected.json",
		},
		{
			description:  "node differ",
			subcommand:   "diff",
			imageA:       nodeBase,
			imageB:       nodeModified,
			differFlags:  []string{"--type=node", "--no-cache"},
			expectedFile: "node_diff_order_expected.json",
		},
		{
			description:  "multi differ",
			subcommand:   "diff",
			imageA:       multiBase,
			imageB:       multiModified,
			differFlags:  []string{"--type=node", "--type=pip", "--type=apt", "--no-cache"},
			expectedFile: "multi_diff_expected.json",
		},
		{
			description:  "multi differ local",
			subcommand:   "diff",
			imageA:       multiBaseLocal,
			imageB:       multiModifiedLocal,
			differFlags:  []string{"--type=node", "--type=pip", "--type=apt", "--no-cache"},
			expectedFile: "multi_diff_expected.json",
		},
		{
			description:  "history differ",
			subcommand:   "diff",
			imageA:       diffBase,
			imageB:       diffModified,
			differFlags:  []string{"--type=history", "--no-cache"},
			expectedFile: "hist_diff_expected.json",
		},
		{
			description:  "metadata differ",
			subcommand:   "diff",
			imageA:       metadataBase,
			imageB:       metadataModified,
			differFlags:  []string{"--type=metadata", "--no-cache"},
			expectedFile: "metadata_diff_expected.json",
		},
		{
			description:  "apt sorted differ",
			subcommand:   "diff",
			imageA:       aptBase,
			imageB:       aptModified,
			differFlags:  []string{"--type=apt", "-o", "--no-cache"},
			expectedFile: "apt_sorted_diff_expected.json",
		},
		{
			description:  "apt analysis",
			subcommand:   "analyze",
			imageA:       aptModified,
			differFlags:  []string{"--type=apt", "--no-cache"},
			expectedFile: "apt_analysis_expected.json",
		},
		{
			description:  "size analysis",
			subcommand:   "analyze",
			imageA:       diffBase,
			differFlags:  []string{"--type=size", "--no-cache"},
			expectedFile: "size_analysis_expected.json",
		},
		{
			description:  "size layer analysis",
			subcommand:   "analyze",
			imageA:       diffLayerBase,
			differFlags:  []string{"--type=sizelayer", "--no-cache"},
			expectedFile: "size_layer_analysis_expected.json",
		},
		{
			description:  "pip analysis",
			subcommand:   "analyze",
			imageA:       pipModified,
			differFlags:  []string{"--type=pip", "--no-cache"},
			expectedFile: "pip_analysis_expected.json",
		},
		{
			description:  "node analysis",
			subcommand:   "analyze",
			imageA:       nodeModified,
			differFlags:  []string{"--type=node", "--no-cache"},
			expectedFile: "node_analysis_expected.json",
		},
	}
	for _, test := range tests {
		// Capture the range variable for parallel testing.
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			args := []string{test.subcommand, test.imageA}
			if test.imageB != "" {
				args = append(args, test.imageB)
			}
			args = append(args, test.differFlags...)
			args = append(args, "-j")
			actual, stderr, err := runner.Run(args...)
			if err != nil {
				t.Fatalf("Error running command: %s. Stderr: %s", err, stderr)
			}
			e, err := ioutil.ReadFile(test.expectedFile)
			if err != nil {
				t.Fatalf("Error reading expected file output file: %s", err)
			}
			actual = strings.TrimSpace(actual)
			expected := strings.TrimSpace(string(e))
			if actual != expected {
				t.Errorf("Error actual output does not match expected.  \n\nExpected: %s\n\n Actual: %s\n\n, Stderr: %s", expected, actual, stderr)
			}
		})
	}
}

func newClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("Error getting docker client: %s", err)
	}
	cli.NegotiateAPIVersion(context.Background())

	return cli, nil
}

func TestMain(m *testing.M) {
	// setup
	ctx := context.Background()
	cli, _ := newClient()
	closer, err := cli.ImagePull(ctx, multiBase, types.ImagePullOptions{})
	if err != nil {
		fmt.Printf("Error retrieving docker client: %s", err)
		os.Exit(1)
	}
	io.Copy(os.Stdout, closer)

	closer, err = cli.ImagePull(ctx, multiModified, types.ImagePullOptions{})
	if err != nil {
		fmt.Printf("Error retrieving docker client: %s", err)
		os.Exit(1)
	}
	io.Copy(os.Stdout, closer)

	closer.Close()
	os.Exit(m.Run())
}

func TestConsoleOutput(t *testing.T) {
	runner := ContainerDiffRunner{
		t:          t,
		binaryPath: "../out/container-diff",
	}

	tests := []struct {
		description    string
		subCommand     string
		extraFlag      string
		expectedOutput []string
		producesError  bool
	}{
		{
			description: "analysis --help",
			subCommand:  "analyze",
			extraFlag:   "--help",
			expectedOutput: []string{
				"Analyzes an image using the specifed analyzers as indicated via --type flag(s).",
				"For details on how to specify images, run: container-diff help",
				"container-diff analyze image [flags]",
				"-c, --cache-dir string      cache directory base to create .container-diff (default is $HOME).",
				"-j, --json                  JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).",
				"-w, --output string         output file to write to (default writes to the screen).",
				"-t, --type multiValueFlag   This flag sets the list of analyzer types to use.",
			},
		},
		{
			description: "analysis help",
			subCommand:  "analyze",
			extraFlag:   "help",
			expectedOutput: []string{
				"Analyzes an image using the specifed analyzers as indicated via --type flag(s).",
				"For details on how to specify images, run: container-diff help",
				"container-diff analyze image [flags]",
				"-c, --cache-dir string      cache directory base to create .container-diff (default is $HOME).",
				"-j, --json                  JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).",
				"-w, --output string         output file to write to (default writes to the screen).",
				"-t, --type multiValueFlag   This flag sets the list of analyzer types to use.",
			},
		},
		{
			description: "container-diff --help",
			subCommand:  "--help",
			extraFlag:   "",
			expectedOutput: []string{
				"container-diff is a CLI tool for analyzing and comparing container images.",
				"Images can be specified from either a local Docker daemon, or from a remote registry.",
				"analyze     Analyzes an image: container-diff image",
				"diff        Compare two images: container-diff image1 image2",
				"--format string                             Format to output diff in.",
				"--skip-tls-verify-registry multiValueFlag   Insecure registry ignoring TLS verify to push and pull. Set it repeatedly for multiple registries.",
				"-v, --verbosity string                          This flag controls the verbosity of container-diff. (default \"warning\")",
			},
		},
		{
			description: "container-diff help",
			subCommand:  "help",
			extraFlag:   "",
			expectedOutput: []string{
				"container-diff is a CLI tool for analyzing and comparing container images.",
				"Images can be specified from either a local Docker daemon, or from a remote registry.",
				"analyze     Analyzes an image: container-diff image",
				"diff        Compare two images: container-diff image1 image2",
				"--format string                             Format to output diff in.",
				"--skip-tls-verify-registry multiValueFlag   Insecure registry ignoring TLS verify to push and pull. Set it repeatedly for multiple registries.",
				"-v, --verbosity string                          This flag controls the verbosity of container-diff. (default \"warning\")",
			},
		},
		{
			description: "container-diff diff --help",
			subCommand:  "diff",
			extraFlag:   "--help",
			expectedOutput: []string{
				"Compares two images using the specifed analyzers as indicated via --type flag(s).",
				"For details on how to specify images, run: container-diff help",
				"container-diff diff image1 image2 [flags]",
				"-c, --cache-dir string      cache directory base to create .container-diff (default is $HOME).",
				"-j, --json                  JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).",
				"-w, --output string         output file to write to (default writes to the screen).",
				"--skip-tls-verify-registry multiValueFlag   Insecure registry ignoring TLS verify to push and pull. Set it repeatedly for multiple registries.",
			},
		},
		{
			description: "container-diff diff --help",
			subCommand:  "diff",
			extraFlag:   "help",
			expectedOutput: []string{
				"Error: 'diff' requires two images as arguments: container-diff diff [image1] [image2]",
				"container-diff diff image1 image2 [flags]",
				"-c, --cache-dir string      cache directory base to create .container-diff (default is $HOME).",
				"-j, --json                  JSON Output defines if the diff should be returned in a human readable format (false) or a JSON (true).",
				"-w, --output string         output file to write to (default writes to the screen).",
				"--skip-tls-verify-registry multiValueFlag   Insecure registry ignoring TLS verify to push and pull. Set it repeatedly for multiple registries.",
			},
			producesError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			args := []string{test.subCommand}
			if test.extraFlag != "" {
				args = append(args, test.extraFlag)
			}
			actual, stderr, err := runner.Run(args...)
			if err != nil {
				if !test.producesError {
					t.Fatalf("Error running command: %s. Stderr: %s", err, stderr)
				} else {
					actual = err.Error()
				}
			}
			actual = strings.TrimSpace(actual)
			for _, expectedLine := range test.expectedOutput {
				if !strings.Contains(actual, expectedLine) {
					t.Errorf("Error actual output does not contain expected line.  \n\nExpected: %s\n\n Actual: %s\n\n, Stderr: %s", expectedLine, actual, stderr)
				}
			}

		})
	}
}
