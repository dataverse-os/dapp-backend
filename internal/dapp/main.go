package dapp

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dataverse-os/dapp-backend/ceramic"
)

var (
	CeramicSession  ceramic.Session
	IsSandbox       = os.Getenv("IS_SANDBOX") == "TRUE"
	ceramicURL      = os.Getenv("CERAMIC_URL")
	ceramicAdminKey = os.Getenv("CERAMIC_ADMIN_KEY")
)

func InitCeramicSession() {
	var err error
	if CeramicSession, err = ceramic.NewSession(ceramicURL, ceramicAdminKey); err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err = ceramic.Default.CheckAdminAccess(ctx, CeramicSession); err != nil {
		log.Fatalf("failed to parse ceramic url with error: %s", err)
	}
}
