package ceramic

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ComposeDB *gorm.DB

func InitComposeDB() {
	var err error
	if ComposeDB, err = gorm.Open(postgres.Open(os.Getenv("CERAMIC_POSTGRES_DSN")), &gorm.Config{}); err != nil {
		log.Fatalln("failed init ceramic postgres connection:", err)
	}
}
