package main

import (
	"os"

	"github.com/rm2013/kubectl-clogs/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var version = "0.0.1"

func main() {
	cmd.SetVersion(version)

	cLogsCmd := cmd.NewCmdClogs(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := cLogsCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
