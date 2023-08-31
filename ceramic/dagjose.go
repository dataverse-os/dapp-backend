package ceramic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ceramicnetwork/go-dag-jose/dagjose"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/datamodel"
)

type DagJWS struct {
	Link       cid.Cid
	Payload    cid.Cid
	Signatures []Signature
}

type Signature struct {
	Header    any
	Protected json.RawMessage
	Signature []byte
}

func (sig Signature) String() string {
	return fmt.Sprintf("header: %s\n", sig.Header) +
		fmt.Sprintf("protected: %s\n", string(sig.Protected)) +
		fmt.Sprintf("signature: %s\n", hexutil.Encode(sig.Signature))
}

var (
	dagjsonEncodeOption = dagjson.EncodeOptions{
		EncodeLinks: true,
		EncodeBytes: true,
		MapSortMode: codec.MapSortMode_None,
	}
)

func LoadDagJWS(ctx context.Context, node *rpc.HttpApi, logCid cid.Cid) (dagJws dagjose.DecodedJWS, err error) {
	var blockReader io.Reader
	if blockReader, err = node.Block().Get(ctx, path.IpfsPath(logCid)); err != nil {
		return
	}
	return DecodeDagJWSFromReader(blockReader)
}

func DecodeDagJWSNodeDataFromReader(reader io.Reader) (nd datamodel.Node, err error) {
	builder := dagjose.Type.DecodedJWS.NewBuilder()
	cfg := dagjose.DecodeOptions{}
	if err = cfg.DecodeJWS(builder, reader); err != nil {
		return
	}
	nd = builder.Build()
	return
}

func DecodeDagJWSFromReader(reader io.Reader) (dagJws dagjose.DecodedJWS, err error) {
	var nd datamodel.Node
	if nd, err = DecodeDagJWSNodeDataFromReader(reader); err != nil {
		return
	}
	var ok bool
	if dagJws, ok = nd.(dagjose.DecodedJWS); !ok {
		err = fmt.Errorf("cannot asset %s as dagjose.DecodedJWS", nd)
	}
	return
}

func decodeBase64Url(field dagjose.Base64Url) (data []byte, err error) {
	str, err := field.AsString()
	if err != nil {
		return
	}
	data, err = base64.RawURLEncoding.DecodeString(str)
	return
}
