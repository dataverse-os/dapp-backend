package dapp

import (
	"context"
	"encoding/binary"
	"log"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ethereum/go-ethereum/common"
	bolt "go.etcd.io/bbolt"
)

var (
	BoltDB        *bolt.DB
	IndexedModels = map[string]struct{}{}
)

func InitIndexedModels() (err error) {
	modelIds, err := ceramic.Default.GetIndexedModels(context.Background(), CeramicSession)
	if err != nil {
		return
	}
	for _, v := range modelIds {
		IndexedModels[v] = struct{}{}
	}
	return
}

var (
	bucketModelVersion = []byte("version")
)

func InitBolt() {
	var err error
	if err = EnsureDir(AppBaseDir()); err != nil {
		log.Fatalln(err)
	}
	BoltDB, err = bolt.Open(
		ModelsBoltPath(),
		0640, bolt.DefaultOptions)
	if err != nil {
		log.Fatalln(err)
	}
	if err = BoltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketModelVersion)
		return err
	}); err != nil {
		log.Fatalln(err)
	}
}

func LookupUserModelVersion(address common.Address, modelName string) (version int64, err error) {
	if err = BoltDB.View(func(tx *bolt.Tx) error {
		result := tx.Bucket(bucketModelVersion).Get(
			append(address.Bytes(),
				[]byte(modelName)...,
			),
		)
		if result != nil {
			version = int64(binary.BigEndian.Uint64(result))
		} else {
			version = -1
		}
		return nil
	}); err != nil {
		return
	}
	return
}

func UpdateUserModelVersion(address common.Address, modelName string, version uint64) (err error) {
	if err = BoltDB.Update(func(tx *bolt.Tx) error {
		var versionBinary []byte = make([]byte, 8)
		binary.BigEndian.PutUint64(versionBinary, version)
		return tx.Bucket(bucketModelVersion).Put(
			append(
				address.Bytes(),
				[]byte(modelName)...,
			),
			versionBinary,
		)
	}); err != nil {
		return
	}
	return
}
