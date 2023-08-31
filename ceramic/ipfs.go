package ceramic

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipld/go-ipld-prime/datamodel"
)

func NewIpfsImpl(node *rpc.HttpApi) IpfsImpl {
	return IpfsImpl{node: node}
}

type IpfsImpl struct {
	node *rpc.HttpApi
}

func (impl IpfsImpl) LoadStreamState(ctx context.Context, tip cid.Cid) (state StreamState, err error) {
	var commits []NodeDataDecoder
	if commits, err = impl.LoadCommits(ctx, tip); err != nil {
		return
	}
	for _, commit := range commits {
		if stateBuilder, ok := commit.(CommitPayload); ok {
			if err = stateBuilder.ApplyToStream(&state); err != nil {
				return
			}
		}
	}
	return
}

func (impl IpfsImpl) LoadCommits(ctx context.Context, c cid.Cid) (commits []NodeDataDecoder, err error) {

	return
}

func (impl IpfsImpl) LoadLog(ctx context.Context, c cid.Cid) (commit NodeDataDecoder, err error) {
	var (
		blkReader io.Reader
		nd        datamodel.Node
	)
	if blkReader, err = impl.node.Block().Get(ctx, path.IpfsPath(c)); err != nil {
		return
	}
	if c.Prefix().Codec == cid.DagJOSE {
		commit = &SignedCommit{}
		if nd, err = DecodeDagJWSNodeDataFromReader(blkReader); err != nil {
			return
		}
		var buf bytes.Buffer
		if err = dagjsonEncodeOption.Encode(nd, &buf); err != nil {
			return
		}
		fmt.Println(c.Prefix().Codec, buf.String())
		err = commit.DecodeFromNodeData(nd)
		impl.LoadLog(ctx, commit.(*SignedCommit).Payload)
		return
	}

	if nd, err = DecodeDagCborNodeDataFromReader(blkReader); err != nil {
		return
	}
	var buf bytes.Buffer
	if err = dagjsonEncodeOption.Encode(nd, &buf); err != nil {
		return
	}
	fmt.Println(c.Prefix().Codec, buf.String())

	return nil, nil
}
