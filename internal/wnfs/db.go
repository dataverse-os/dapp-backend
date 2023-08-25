package wnfs

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func init() {
	var err error
	if db, err = gorm.Open(postgres.Open(os.Getenv("WNFS_POSTGRES_DSN")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "wnfs.",
			SingularTable: false,
		},
	}); err != nil {
		log.Fatalln("failed init wnfs postgres connection:", err)
	}

	db.AutoMigrate(
		&CommitProof{},
	)
}
