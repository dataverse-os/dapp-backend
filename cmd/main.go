package main

import (
	"dapp-backend/ceramic"
	"dapp-backend/internal/routers"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

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
		}
		routers.InitRouter()
		routers.Start()
	},
}

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

func ExecuteCommand(name string, subname string, args ...string) (string, error) {
	args = append([]string{subname}, args...)

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()

	return string(bytes), err
}

func Error(cmd *cobra.Command, args []string, err error) {
	fmt.Fprintf(os.Stderr, "execute %s args:%v error:%v\n", cmd.Name(), args, err)
	os.Exit(1)
}

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
