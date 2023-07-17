package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var generateKeyCmd = &cobra.Command{
	Use: "generate-key",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func checkPrivateKey(cmd *cobra.Command, args []string) error {
	if _, ok := os.LookupEnv("CERAMIC_ADMIN_KEY"); !ok {
		err := fmt.Errorf("environment CERAMIC_ADMIN_KEY not found, Generate a private key with generate-key command")
		return err
	}
	return nil
}
