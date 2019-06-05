package commands

import (
	"spark-cluster/cmd/sparkctl/commands/job"

	"github.com/spf13/cobra"
)

const (
	// CLIName is the name of this CLI
	CLIName = "sparkctl"
)

func NewCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   CLIName,
		Short: "sparkctl is the command line interface to submit spark jobs",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(job.NewSubmitCommand())

	return command
}
