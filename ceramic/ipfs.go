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
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
)

type Commit struct {
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
		fmt.Sprintf("signature: %s", hexutil.Encode(sig.Signature))
}

func ConvertFrom(dagJws dagjose.DecodedJWS) (commit Commit, err error) {
	payloadData, err := decodeBase64Url(dagJws.FieldPayload())
	if err != nil {
		fmt.Println(err)
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

func LoadDagJWS(ctx context.Context, node *rpc.HttpApi, logCid cid.Cid) (dagJws dagjose.DecodedJWS, err error) {
	var blockReader io.Reader
	if blockReader, err = node.Block().Get(ctx, path.IpfsPath(logCid)); err != nil {
		return
	}
	return BuildDagJWSFromReader(blockReader)
}

func BuildDagJWSFromReader(reader io.Reader) (dagJws dagjose.DecodedJWS, err error) {
	builder := dagjose.Type.DecodedJWS.NewBuilder()
	cfg := dagjose.DecodeOptions{
		AddLink: true,
	}
	if err = cfg.DecodeJWS(builder, reader); err != nil {
		return
	}
	var ok bool
	if dagJws, ok = builder.Build().(dagjose.DecodedJWS); !ok {
		err = fmt.Errorf("cannot asset %s as dagjose.DecodedJWS", builder.Build())
	}
	return
}

func BuildPayloadFromReader(reader io.Reader) {
	builder := basicnode.Prototype.Map.NewBuilder()
	_ = builder
}

//nolint:unused
func lookupIgnoreAbsent(key string, n datamodel.Node) (datamodel.Node, error) {
	value, err := n.LookupByString(key)
	if err != nil {
		if _, notFoundErr := err.(datamodel.ErrNotExists); !notFoundErr {
			return nil, err
		}
	}
	if value == datamodel.Absent {
		value = nil
	}
	return value, nil
}

//nolint:unused
func lookupIgnoreNoSuchField(key string, n datamodel.Node) (datamodel.Node, error) {
	value, err := lookupIgnoreAbsent(key, n)
	if err != nil {
		if _, noSuchFieldErr := err.(schema.ErrNoSuchField); !noSuchFieldErr {
			return nil, err
		}
	}
	return value, nil
}

func decodeBase64Url(field dagjose.Base64Url) (data []byte, err error) {
	str, err := field.AsString()
	if err != nil {
		fmt.Println(err)
	}
	data, err = base64.RawURLEncoding.DecodeString(str)
	return
}

//nolint:unused
func decodeBase64UrlString(field dagjose.Base64Url) (string, error) {
	data, err := decodeBase64Url(field)
	return string(data), err
}
