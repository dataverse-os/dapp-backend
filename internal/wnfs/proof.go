package wnfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

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
	"/api/v0/streams/": {},
	"/api/v0/commits/": {},
}

var CollectionResponsePath = map[string]struct{}{
	"/api/v0/collection": {},
}

func AppendToStreamVerifyGroup(buf *bytes.Buffer) {
	var stream ceramic.Stream
	if err := json.Unmarshal(buf.Bytes(), &stream); err != nil {
		fmt.Println(err)
		return
	}
	VerifyGroup.Go(func() (err error) {
		if err = StoreAndVerifyContentHash(context.Background(), stream.State); err != nil {
			return
		}
		return
	})
}

func AppendToCollectionVerifyGroup(buf *bytes.Buffer) {
	var collection ceramic.Collection
	if err := json.Unmarshal(buf.Bytes(), &collection); err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range collection.Edges {
		state := v.Node
		VerifyGroup.Go(func() (err error) {
			if err = StoreAndVerifyContentHash(context.Background(), state); err != nil {
				return
			}
			return
		})
	}
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

func StoreAndVerifyContentHash(ctx context.Context, streamState ceramic.StreamState) (err error) {
	var (
		streamId            = streamState.StreamId()
		commit              ceramic.Stream
		commitProofs        = make([]CommitProof, len(streamState.Log))
		commitProofsInDB    []CommitProof
		commitProofsInDBMap = make(map[ceramic.StreamId]CommitProof)
		commitIds           = streamState.CommitIds()
	)

	defer func() {
		if err != nil {
			fmt.Printf("wnfs check %s with %d commits\n", streamState.StreamId(), len(commitProofs))
		}
	}()

	if err = db.WithContext(ctx).Where(&CommitProof{
		StreamId: streamId,
	}).Find(&commitProofsInDB).Error; err != nil {
		return
	}
	for _, proof := range commitProofsInDB {
		commitProofsInDBMap[proof.CommitId] = proof
	}
	for i, v := range streamState.Log {
		if proof, exists := commitProofsInDBMap[commitIds[i]]; exists {
			commitProofs[i] = proof
			continue
		}
		commitProofs[i] = CommitProof{
			CommitId: commitIds[i],
			StreamId: streamId,
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
