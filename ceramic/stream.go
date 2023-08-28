package ceramic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ipfs/go-cid"
)

type Stream struct {
	StreamId StreamId    `json:"streamId"`
	State    StreamState `json:"state"`
}

type StreamState struct {
	Type         StreamType      `json:"type"`
	Content      json.RawMessage `json:"content"`
	Metadata     json.RawMessage `json:"metadata"`
	Signature    int             `json:"signature"`
	AnchorStatus string          `json:"anchorStatus"`
	Log          []StreamLog     `json:"log"`
	AnchorProof  *AnchorProof    `json:"anchorProof,omitempty"`
	DocType      string          `json:"docType"`
}

type StreamLog struct {
	Cid            string     `json:"cid"`
	Type           CommitType `json:"type"`
	ExpirationTime *uint64    `json:"expirationTime,omitempty"`
	Timestamp      uint64     `json:"timestamp"`
}

type AnchorProof struct {
	Root    string `json:"root"`
	TxHash  string `json:"txHash"`
	TxType  string `json:"txType"`
	ChainId string `json:"chainId"`
}

func (state StreamState) StreamId() StreamId {
	return StreamId{
		Type: state.Type,
		Cid:  cid.MustParse(state.Log[0].Cid),
	}
}

func (state StreamState) CommitIds() (commitIds []StreamId) {
	streamId := state.StreamId()
	commitIds = append(commitIds, streamId.Genesis())
	for _, v := range state.Log[1:] {
		commitIds = append(commitIds, streamId.With(v.Cid))
	}
	return
}

func (state StreamState) ContentHash() (sum [32]byte, err error) {
	// var data any
	// if err = json.Unmarshal(state.Content, &data); err != nil {
	// 	return
	// }
	// var buf bytes.Buffer
	// if err = json.NewEncoder(&buf).Encode(data); err != nil {
	// 	return
	// }

	// sum = sha256.Sum256(buf.Bytes())

	sum = sha256.Sum256(state.Content)
	hex.EncodeToString(sum[0:])
	return
}

func GetStreamId(ctx context.Context, streamId StreamId) (stream Stream, err error) {
	url := fmt.Sprintf("%s/api/v0/streams/%s", os.Getenv("CERAMIC_URL"), streamId)
	var req *http.Request
	var resp *http.Response
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	if err = json.NewDecoder(resp.Body).Decode(&stream); err != nil {
		return
	}
	return
}
