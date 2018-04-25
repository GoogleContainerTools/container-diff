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
	diffBase     = "gcr.io/gcp-runtimes/diff-base"
	diffModified = "gcr.io/gcp-runtimes/diff-modified"

	metadataBase     = "gcr.io/gcp-runtimes/metadata-base"
	metadataModified = "gcr.io/gcp-runtimes/metadata-modified"

	aptBase     = "gcr.io/gcp-runtimes/apt-base"
	aptModified = "gcr.io/gcp-runtimes/apt-modified"

	rpmBase     = "valentinrothberg/containerdiff:diff-base"
	rpmModified = "valentinrothberg/containerdiff:diff-modified"

	// Why is this node-modified:2.0?
	nodeBase     = "gcr.io/gcp-runtimes/node-modified:2.0"
	nodeModified = "gcr.io/gcp-runtimes/node-modified"

	pipModified = "gcr.io/gcp-runtimes/pip-modified"

	multiBase     = "gcr.io/gcp-runtimes/multi-base"
	multiModified = "gcr.io/gcp-runtimes/multi-modified"

	multiBaseLocal     = "daemon://gcr.io/gcp-runtimes/multi-base"
	multiModifiedLocal = "daemon://gcr.io/gcp-runtimes/multi-modified"
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
			differFlags:  []string{"--type=file"},
			expectedFile: "file_diff_expected.json",
		},
		{
			description:  "apt differ",
			subcommand:   "diff",
			imageA:       aptBase,
			imageB:       aptModified,
			differFlags:  []string{"--type=apt"},
			expectedFile: "apt_diff_expected.json",
		},
		// {
		// 	description:  "rpm differ",
		// 	subcommand:   "diff",
		// 	imageA:       rpmBase,
		// 	imageB:       rpmModified,
		// 	differFlags:  []string{"--type=rpm"},
		// 	expectedFile: "rpm_diff_expected.json",
		// },
		{
			description:  "node differ",
			subcommand:   "diff",
			imageA:       nodeBase,
			imageB:       nodeModified,
			differFlags:  []string{"--type=node"},
			expectedFile: "node_diff_order_expected.json",
		},
		{
			description:  "multi differ",
			subcommand:   "diff",
			imageA:       multiBase,
			imageB:       multiModified,
			differFlags:  []string{"--type=node", "--type=pip", "--type=apt"},
			expectedFile: "multi_diff_expected.json",
		},
		{
			description:  "multi differ local",
			subcommand:   "diff",
			imageA:       multiBaseLocal,
			imageB:       multiModifiedLocal,
			differFlags:  []string{"--type=node", "--type=pip", "--type=apt"},
			expectedFile: "multi_diff_expected.json",
		},
		{
			description:  "history differ",
			subcommand:   "diff",
			imageA:       diffBase,
			imageB:       diffModified,
			differFlags:  []string{"--type=history"},
			expectedFile: "hist_diff_expected.json",
		},
		{
			description:  "metadata differ",
			subcommand:   "diff",
			imageA:       metadataBase,
			imageB:       metadataModified,
			differFlags:  []string{"--type=metadata"},
			expectedFile: "metadata_diff_expected.json",
		},
		{
			description:  "apt sorted differ",
			subcommand:   "diff",
			imageA:       aptBase,
			imageB:       aptModified,
			differFlags:  []string{"--type=apt", "-o"},
			expectedFile: "apt_sorted_diff_expected.json",
		},
		{
			description:  "apt analysis",
			subcommand:   "analyze",
			imageA:       aptModified,
			differFlags:  []string{"--type=apt"},
			expectedFile: "apt_analysis_expected.json",
		},
		// {
		// 	description:  "rpm analysis",
		// 	subcommand:   "analyze",
		// 	imageA:       rpmModified,
		// 	differFlags:  []string{"--type=rpm"},
		// 	expectedFile: "rpm_analysis_expected.json",
		// },
		{
			description:  "file sorted analysis",
			subcommand:   "analyze",
			imageA:       diffModified,
			differFlags:  []string{"--type=file", "-o"},
			expectedFile: "file_sorted_analysis_expected.json",
		},
		{
			description:  "pip analysis",
			subcommand:   "analyze",
			imageA:       pipModified,
			differFlags:  []string{"--type=pip"},
			expectedFile: "pip_analysis_expected.json",
		},
		{
			description:  "node analysis",
			subcommand:   "analyze",
			imageA:       nodeModified,
			differFlags:  []string{"--type=node"},
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
