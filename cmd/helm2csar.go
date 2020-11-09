package cmd

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "h2c <subcommand>",
		Short: "Simple Helm-2-CSAR generator",
	}
	c.AddCommand(NewGenerateCommand())

	return &c
}
