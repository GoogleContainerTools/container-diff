package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/container-diff/cmd"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	glog.Flush()
}
