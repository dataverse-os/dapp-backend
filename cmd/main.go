package main

import (
	"errors"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "dapp-backend",
	Run: func(cmd *cobra.Command, args []string) {
		Error(cmd, args, errors.New("unrecognized command"))
	},
}

func main() {
	rootCmd.AddCommand(daemonCmd, generateKeyCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
