package wnfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

var VerifyGroup *errgroup.Group

func init() {
	VerifyGroup = new(errgroup.Group)
	VerifyGroup.SetLimit(1)
}

var StreamResponsePath = map[string]struct{}{
	"/api/v0/streams": {},
	"/api/v0/commits": {},
}

func AppendToVerifyGroup(buf *bytes.Buffer) {
	VerifyGroup.Go(func() (err error) {
		if err = StoreAndVerifyContentHashFromReader(context.Background(), buf); err != nil {
			log.Println(err)
		}
		return
	})
}

type CommitProofStatus uint64

const (
	CommitProofStatusUnverified CommitProofStatus = iota
	CommitProofStatusLegal
	CommitProofStatusIllegal
)

type CommitProof struct {
	CommitId ceramic.StreamId  `json:"commitId" gorm:"primaryKey"`
	StreamId ceramic.StreamId  `json:"streamId"`
	Cid      string            `json:"cid"`
	Hash     common.Hash       `json:"hash"`
	Content  json.RawMessage   `json:"content" gorm:"type:jsonb"`
	Status   CommitProofStatus `json:"status"`
}

func StoreAndVerifyContentHashFromReader(ctx context.Context, buf *bytes.Buffer) (err error) {
	var stream ceramic.Stream
	if err = json.Unmarshal(buf.Bytes(), &stream); err != nil {
		return
	}
	if err = StoreAndVerifyContentHash(ctx, stream); err != nil {
		return
	}
	return
}

func StoreAndVerifyContentHash(ctx context.Context, stream ceramic.Stream) (err error) {
	var (
		commit              ceramic.Stream
		commitProofs        = make([]CommitProof, len(stream.State.Log))
		commitProofsInDB    []CommitProof
		commitProofsInDBMap = make(map[ceramic.StreamId]CommitProof)
		commitIds           = stream.State.CommitIds()
	)

	defer func() {
		if err != nil {
			fmt.Printf("wnfs check %s with %d commits\n", stream.StreamId, len(commitProofs))
		}
	}()

	if err = db.WithContext(ctx).Where(&CommitProof{
		StreamId: stream.StreamId,
	}).Find(&commitProofsInDB).Error; err != nil {
		return
	}
	for _, proof := range commitProofsInDB {
		commitProofsInDBMap[proof.CommitId] = proof
	}
	for i, v := range stream.State.Log {
		if proof, exists := commitProofsInDBMap[commitIds[i]]; exists {
			commitProofs[i] = proof
			continue
		}
		commitProofs[i] = CommitProof{
			CommitId: commitIds[i],
			StreamId: stream.StreamId,
			Cid:      v.Cid,
		}
		if commit, err = commitIds[i].GetStream(ctx); err != nil {
			return
		}
		commitProofs[i].Content = commit.State.Content
		if commitProofs[i].Hash, err = commit.State.ContentHash(); err != nil {
			return
		}
	}
	shouldInsert := make([]CommitProof, 0)
	for _, proof := range commitProofs {
		if _, exists := commitProofsInDBMap[proof.CommitId]; !exists {
			shouldInsert = append(shouldInsert, proof)
		}
	}
	if len(shouldInsert) != 0 {
		if err = db.Create(shouldInsert).Error; err != nil {
			return
		}
	}
	return
}
