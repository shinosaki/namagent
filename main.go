package main

import (
	"github.com/shinosaki/namagent/cli"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "namagent",
		Short: "Live streaming alert built with Golang",
	}

	cmd.AddCommand(cli.Alert())
	cmd.AddCommand(cli.Recorder())

	cmd.Execute()
}
