package ceramic

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
)

type StreamIdType uint64

// https://github.com/ceramicnetwork/CIPs/blob/main/CIPs/cip-59.md
const (
	StreamIdTypeTile StreamIdType = iota
	StreamIdTypeCaip10Link
	StreamIdTypeModel
	StreamIdTypeModelInstanceDocument
	StreamIdTypeUnloadable
	StreamIdTypeEventId

	StreamIDEncoding = multibase.Base36
)

type StreamID struct {
	Type StreamIdType
	Cid  cid.Cid
	Log  cid.Cid
}

var (
	streamIdEncoder = multibase.MustNewEncoder(StreamIDEncoding)
)

func (id StreamID) String() string {
	var buf bytes.Buffer
	buf.Write(binary.AppendUvarint(nil, uint64(multicodec.Streamid)))
	buf.Write(binary.AppendUvarint(nil, uint64(id.Type)))
	buf.Write(id.Cid.Bytes())
	if id.Log.ByteLen() != 0 {
		buf.Write(id.Log.Bytes())
	}
	return streamIdEncoder.Encode(buf.Bytes())
}

func NewStreamID(t StreamIdType, cidStr ...string) (id StreamID, err error) {
	if len(cidStr) != 1 && len(cidStr) != 2 {
		err = fmt.Errorf("unexpect cid length, could only 1 or 2 (genesis cid and log cid)")
	}
	id.Type = t

	return
}

func ParseStreamID(str string) (id StreamID, err error) {
	var (
		buf         []byte
		encoding    multibase.Encoding
		streamCodec uint64
		streamType  uint64
		idx         int
	)

	if encoding, buf, err = multibase.Decode(str); err != nil {
		return
	}
	if encoding != StreamIDEncoding {
		err = fmt.Errorf("unexpected encoding id %c with input %s", encoding, str)
		return
	}
	// check <multicodec-streamCodec>
	if streamCodec, idx = binary.Uvarint(buf); idx <= 0 {
		err = fmt.Errorf("unable to unpack stream codec %v", buf)
		return
	}
	if multicodec.Code(streamCodec) != multicodec.Streamid {
		err = fmt.Errorf("unexpected multicodec %x != 0xce", buf[0])
		return
	}
	buf = buf[idx:]

	// check <stream-type>
	if streamType, idx = binary.Uvarint(buf); idx <= 0 {
		err = fmt.Errorf("unable to unpack stream type %v", buf)
		return
	}
	id.Type = StreamIdType(streamType)
	buf = buf[idx:]

	var nr int
	if nr, id.Cid, err = cid.CidFromBytes(buf); err != nil {
		return
	}
	if len(buf) != nr {
		if _, id.Log, err = cid.CidFromBytes(buf[nr:]); err != nil {
			return
		}
	}
	return
}
