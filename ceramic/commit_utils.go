package ceramic

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/samber/lo"
)

func MustParseLink(nd datamodel.Node, key string) cid.Cid {
	link := lo.Must(lo.Must(nd.LookupByString(key)).AsLink()).(cidlink.Link)
	return link.Cid
}

func DecodeDagCborNodeDataFromReader(reader io.Reader) (nd datamodel.Node, err error) {
	builder := basicnode.Prototype.Map.NewBuilder()
	if err = dagcbor.Decode(builder, reader); err != nil {
		return
	}
	nd = builder.Build()
	return
}

func ContainField(nd datamodel.Node, key string) bool {
	return lo.T2(nd.LookupByString(key)).A.IsNull()
}
