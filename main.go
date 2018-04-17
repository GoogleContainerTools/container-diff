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

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoogleContainerTools/container-diff/cmd"
	"github.com/pkg/profile"
)

const containerDiffEnvPrefix = "CONTAINER_DIFF_ENABLE_PROFILING"

func main() {
	flag.Parse()
	if os.Getenv(containerDiffEnvPrefix) == "1" {
		defer profile.Start(profile.TraceProfile).Stop()
	}
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
