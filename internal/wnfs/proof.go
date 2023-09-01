package wnfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var VerifyGroup *errgroup.Group

func init() {
	VerifyGroup = new(errgroup.Group)
	VerifyGroup.SetLimit(1)
}

func IsGetStreamRequest(ctx *gin.Context) bool {
	if ctx.Writer.Status() != http.StatusOK {
		return false
	}
	for k := range StreamResponsePath {
		if strings.HasPrefix(ctx.Request.URL.Path, k) {
			return true
		}
	}
	return false
}

var StreamResponsePath = map[string]struct{}{
	"/api/v0/streams/": {},
	"/api/v0/commits/": {},
}

func HandleStreamVerify(buf *bytes.Buffer) {
	var stream ceramic.Stream
	if err := json.Unmarshal(buf.Bytes(), &stream); err != nil {
		log.Println(err)
		return
	}
	VerifyGroup.Go(func() (err error) {
		if err = StoreAndVerifyContentHash(context.Background(), stream.State); err != nil {
			return
		}
		return
	})
}

func IsGetCollectionRequest(ctx *gin.Context) bool {
	if ctx.Writer.Status() != http.StatusOK {
		return false
	}
	for k := range CollectionResponsePath {
		if strings.HasPrefix(ctx.Request.URL.Path, k) {
			return true
		}
	}
	return false
}

var CollectionResponsePath = map[string]struct{}{
	"/api/v0/collection": {},
}

func HandleCollectionVerify(buf *bytes.Buffer) {
	var collection ceramic.Collection
	if err := json.Unmarshal(buf.Bytes(), &collection); err != nil {
		log.Println(err)
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
	CommitId  ceramic.StreamId  `json:"commitId" gorm:"primaryKey"`
	Timestamp time.Time         `json:"timestamp"`
	StreamId  ceramic.StreamId  `json:"streamId"`
	LogIndex  uint64            `json:"logIndex"`
	Cid       string            `json:"cid"`
	Hash      common.Hash       `json:"hash"`
	Content   json.RawMessage   `json:"content" gorm:"type:jsonb"`
	Status    CommitProofStatus `json:"status"`
}

func StoreAndVerifyContentHash(ctx context.Context, streamState ceramic.StreamState) (err error) {
	var (
		streamId            ceramic.StreamId
		commit              ceramic.Stream
		commitProofs        = make([]CommitProof, len(streamState.Log))
		commitProofsInDB    []CommitProof
		commitProofsInDBMap = make(map[ceramic.StreamId]CommitProof)
		commitIds           []ceramic.StreamId
	)

	if streamId, err = streamState.StreamId(); err != nil {
		return
	}
	if commitIds, err = streamState.CommitIds(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			fmt.Printf("wnfs check %s with %d commits\n", streamId, len(commitProofs))
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
			CommitId:  commitIds[i],
			Timestamp: time.Unix(int64(v.Timestamp), 0),
			StreamId:  streamId,
			LogIndex:  uint64(i),
			Cid:       v.Cid,
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
