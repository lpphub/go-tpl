package cmd

import (
	"go-tpl/internal"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "Main Function",
		Run: func(cmd *cobra.Command, args []string) {
			internal.Serve()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
