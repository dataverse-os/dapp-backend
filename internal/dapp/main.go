package dapp

import (
	"context"
	"log"
	"os"

	"github.com/dataverse-os/dapp-backend/ceramic"
)

var (
	CeramicSession ceramic.Session
)

func InitCeramicSession() {
	var err error
	if CeramicSession, err = ceramic.NewSession(os.Getenv("CERAMIC_URL"), os.Getenv("DID_PRIVATE_KEY")); err != nil {
		log.Fatalln(err)
	}
	if err = ceramic.Default.CheckAdminAccess(context.Background(), CeramicSession); err != nil {
		log.Fatalf("failed to parse ceramic url with error: %s", err)
	}
}
