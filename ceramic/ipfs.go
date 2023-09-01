package ceramic

import (
	"bytes"
	"context"
	"fmt"
	"io"

	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipld/go-ipld-prime/datamodel"
)

func NewIpfsImpl(node *rpc.HttpApi) IpfsImpl {
	return IpfsImpl{
		blockAPI:  node.Block(),
		pubSubAPI: node.PubSub(),
	}
}

type IpfsImpl struct {
	network   string
	blockAPI  iface.BlockAPI
	pubSubAPI iface.PubSubAPI
}

func (impl IpfsImpl) LoadStreamSate(ctx context.Context, streamId StreamId, rules ...CommitVerifier) (state StreamState, err error) {
	if streamId.Log.ByteLen() == 0 {
		if streamId.Log, err = impl.QueryStream(ctx, streamId); err != nil {
			return
		}
	}
	if state, err = impl.LoadStreamCommits(ctx, streamId.Log, rules...); err != nil {
		return
	}
	state.Type = streamId.Type
	return
}

func (impl IpfsImpl) LoadStreamCommits(ctx context.Context, tip cid.Cid, rules ...CommitVerifier) (state StreamState, err error) {
	var commit Commit
	if commit, err = impl.LoadCommit(ctx, tip); err != nil {
		return
	}
	commitLog := StreamLog{Cid: tip.String()}
	var payload CommitPayload
	if signedCommit, ok := commit.(*SignedCommit); ok {
		if payload, err = signedCommit.LoadPayload(ctx, impl); err != nil {
			return
		}
		switch t := payload.(type) {
		case *DataCommitPayload:
			commitLog.Type = CommitTypeSigned
			if err = ValidatePatches(t.Pathes, rules); err != nil {
				return
			}
			if state, err = impl.LoadStreamCommits(ctx, t.Prev, rules...); err != nil {
				return
			}
		case *GenesisCommitPayload:
			commitLog.Type = CommitTypeGenesis
			if err = ValidateData(t.Data, rules); err != nil {
				return
			}
		}
		payload.ApplyToStream(&state)
		state.Log = append(state.Log, commitLog)
	} else if anchorCommit, ok := commit.(*AnchorCommit); ok {
		if state, err = impl.LoadStreamCommits(ctx, anchorCommit.Prev, rules...); err != nil {
			return
		}
		commitLog.Type = CommitTypeAnchor
		state.Log = append(state.Log, commitLog)
	} else {
		err = fmt.Errorf("unreconized commit event")
	}
	return
}

func (impl IpfsImpl) LoadCommit(ctx context.Context, tip cid.Cid) (commit Commit, err error) {
	var (
		blkReader io.Reader
		nd        datamodel.Node
	)
	if blkReader, err = impl.blockAPI.Get(ctx, path.IpfsPath(tip)); err != nil {
		return
	}
	if tip.Prefix().Codec == cid.DagJOSE {
		commit = &SignedCommit{}
		if nd, err = DecodeDagJWSNodeDataFromReader(blkReader); err != nil {
			return
		}
		err = commit.DecodeFromNodeData(nd)
		return
	}

	commit = &AnchorCommit{}
	if nd, err = DecodeDagCborNodeDataFromReader(blkReader); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = dagjsonEncodeOption.Encode(nd, &buf); err != nil {
		return
	}
	err = commit.DecodeFromNodeData(nd)
	return
}
