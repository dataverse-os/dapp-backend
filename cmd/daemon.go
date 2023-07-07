package main

import (
	"log"
	"os"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/dataverse-os/dapp-backend/internal/routers"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon",
	PreRunE: checkPrivateKey,
	Run: func(cmd *cobra.Command, args []string) {
		if _, ok := os.LookupEnv("DID_PRIVATE_KEY"); !ok {
			log.Printf("Environment DID_PRIVATE_KEY not found, Generate a private key now\n\n")
			ceramic.GenerateKey()
			return
		} else {
			routers.InitRouter()
			routers.Start()
		}
	},
}
