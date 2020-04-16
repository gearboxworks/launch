package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"launch/only"
	"strings"
)

func showArgs(cmd *cobra.Command, args []string) {
	// var err error

	for range only.Once {
		flargs := cmd.Flags().Args()
		if flargs != nil {
			fmt.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(flargs, " "))
			break
		}

		fmt.Printf("'%s' called with '%s'\n", cmd.CommandPath(), strings.Join(args, " "))
		break
	}

	fmt.Println("")
}
