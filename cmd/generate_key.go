package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var generateKeyCmd = &cobra.Command{
	Use: "generaye-key",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func checkPrivateKey(cmd *cobra.Command, args []string) error {
	if _, ok := os.LookupEnv("DID_PRIVATE_KEY"); !ok {
		err := fmt.Errorf("environment DID_PRIVATE_KEY not found, Generate a private key with generate-key command")
		return err
	}
	return nil
}
