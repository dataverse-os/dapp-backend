package ceramic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ceramicnetwork/go-dag-jose/dagjose"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/samber/lo"
)

type NodeDataDecoder interface {
	DecodeFromNodeData(nd datamodel.Node) (err error)
}

var CommitTypeDelecters = []CommitTypeDelecter{
	&AnchorCommit{},
}

type CommitWithPayload interface {
	LoadPayload() (payload CommitPayload, err error)
}

type CommitPayload interface {
	NodeDataDecoder
	DelectType(nd datamodel.Node) bool
	ApplyToStream(state *StreamState) (err error)
}

type CommitTypeDelecter interface {
	DelectType(nd datamodel.Node) bool
}

type CommitHeader struct {
	Model       StreamId // raw as StreamID encoded as byte array
	Controllers []string
	Unique      []byte
}

func (header *CommitHeader) DecodeFromNodeData(nd datamodel.Node) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered in %s", r)
			return
		}
	}()
	header.Model = lo.Must(CastStreamID(
		lo.Must(lo.Must(nd.LookupByString("model")).AsBytes()),
	))
	header.Unique = lo.Must(lo.Must(nd.LookupByString("unique")).AsBytes())
	iter := lo.Must(nd.LookupByString("controllers")).ListIterator()
	for !iter.Done() {
		_, n, _ := iter.Next()
		header.Controllers = append(header.Controllers, lo.Must(n.AsString()))
	}
	return
}

var _ CommitPayload = (*GenesisCommitPayload)(nil)

type GenesisCommitPayload struct {
	Header CommitHeader
	Data   json.RawMessage
}

func (commit *GenesisCommitPayload) ApplyToStream(state *StreamState) (err error) {
	state.Content = commit.Data
	return
}

func (*GenesisCommitPayload) DelectType(nd datamodel.Node) bool {
	return ContainField(nd, "header") && !ContainField(nd, "id")
}

func (commit *GenesisCommitPayload) DecodeFromNodeData(nd datamodel.Node) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered in %s", r)
			return
		}
	}()
	var buf bytes.Buffer
	if err = dagjsonEncodeOption.Encode(lo.Must(nd.LookupByString("data")), &buf); err != nil {
		return
	}
	commit.Data = buf.Bytes()
	if headerNode, e := nd.LookupByString("header"); e == nil && !headerNode.IsNull() {
		if err = commit.Header.DecodeFromNodeData(headerNode); err != nil {
			return
		}
	}
	return
}

var _ CommitPayload = (*DataCommitPayload)(nil)

type DataCommitPayload struct {
	ID     cid.Cid // link to init event
	Prev   cid.Cid
	Header *CommitHeader
	Data   json.RawMessage
}

func (commit *DataCommitPayload) DecodeFromNodeData(nd datamodel.Node) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered in %s", r)
			return
		}
	}()
	commit.ID = MustParseLink(nd, "id")
	commit.Prev = MustParseLink(nd, "prev")

	var buf bytes.Buffer
	if err = dagjsonEncodeOption.Encode(lo.Must(nd.LookupByString("data")), &buf); err != nil {
		return
	}
	commit.Data = buf.Bytes()
	return
}

func (DataCommitPayload) DelectType(nd datamodel.Node) bool {
	return ContainField(nd, "prev") && !ContainField(nd, "proof")
}

func (commit *DataCommitPayload) ApplyToStream(state *StreamState) (err error) {
	//TODO
	return
}

var _ NodeDataDecoder = (*AnchorCommit)(nil)

type AnchorCommit struct {
	ID    cid.Cid // link to init event
	Prev  cid.Cid
	Proof cid.Cid
	Path  string
}

func (*AnchorCommit) ApplyToStream(state *StreamState) (err error) {
	return
}

func (*AnchorCommit) DelectType(nd datamodel.Node) bool {
	return ContainField(nd, "proof")
}

func (commit *AnchorCommit) DecodeFromNodeData(nd datamodel.Node) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered in %s", r)
			return
		}
	}()
	commit.ID = MustParseLink(nd, "id")
	commit.Prev = MustParseLink(nd, "prev")
	commit.Proof = MustParseLink(nd, "proof")
	commit.Path = lo.Must(lo.Must(nd.LookupByString("path")).AsString())
	return
}

var _ NodeDataDecoder = (*SignedCommit)(nil)

type SignedCommit struct {
	linkedCommit CommitPayload
	Link         cid.Cid
	Payload      cid.Cid
	Signatures   []Signature
}

func (commit *SignedCommit) GetPayload() CommitPayload {
	return commit.linkedCommit
}

func (commit *SignedCommit) LoadPayload(ctx context.Context, impl IpfsImpl) (err error) {
	var (
		blkReader io.Reader
		nd        datamodel.Node
	)
	if blkReader, err = impl.node.Block().Get(ctx, path.IpfsPath(commit.Payload)); err != nil {
		return
	}
	if nd, err = DecodeDagCborNodeDataFromReader(blkReader); err != nil {
		return
	}
	if ContainField(nd, "prev") {
		commit.linkedCommit = &DataCommitPayload{}
	} else {
		commit.linkedCommit = &GenesisCommitPayload{}
	}
	if err = commit.DecodeFromNodeData(nd); err != nil {
		return
	}
	return
}

func (commit *SignedCommit) DecodeFromNodeData(nd datamodel.Node) (err error) {
	dagJws := nd.(dagjose.DecodedJWS)
	payloadData, err := decodeBase64Url(dagJws.FieldPayload())
	if err != nil {
		return
	}
	if commit.Payload, err = cid.Cast(payloadData); err != nil {
		return
	}
	if link := dagJws.FieldLink(); link.Exists() {
		if link, ok := link.Must().Link().(cidlink.Link); ok {
			commit.Link = link.Cid
		}
	}
	signatureIter := dagJws.FieldSignatures().Must().Iterator()
	for !signatureIter.Done() {
		_, item := signatureIter.Next()
		sig := Signature{}
		if protected := item.FieldProtected(); protected.Exists() {
			if sig.Protected, err = decodeBase64Url(protected.Must()); err != nil {
				return
			}
		}
		if sig.Signature, err = decodeBase64Url(item.FieldSignature()); err != nil {
			return
		}

		commit.Signatures = append(commit.Signatures, sig)
	}
	return
}
