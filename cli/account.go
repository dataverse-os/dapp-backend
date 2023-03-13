package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Create account.",

	Run: func(cmd *cobra.Command, args []string) {
		// output, err := ExecuteCommand("git", "rev-parse", args...)
		privateKey, err := crypto.GenerateKey()
		publicKey := privateKey.Public()
		publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

		if err != nil {
			Error(cmd, args, err)
		}

		fmt.Fprintln(os.Stdout, "[Private Key]    >>> ", hex.EncodeToString(crypto.FromECDSA(privateKey)))
		fmt.Fprintln(os.Stdout, "[Public Address] >>> ", fromAddress.Hex())
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
