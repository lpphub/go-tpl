package cmd

import (
	"go-tpl/server"

	"github.com/spf13/cobra"
)

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "Main Function",
		Run: func(cmd *cobra.Command, args []string) {
			server.Serve()
		},
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err.Error())
	}
}
