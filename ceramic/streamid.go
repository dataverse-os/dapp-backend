package ceramic

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
)

type (
	StreamType uint64
	CommitType uint64
)

// https://github.com/ceramicnetwork/CIPs/blob/main/CIPs/cip-59.md
const (
	StreamTypeTile StreamType = iota
	StreamTypeCaip10Link
	StreamTypeModel
	StreamTypeModelInstanceDocument
	StreamTypeUnloadable
	StreamTypeEventId
)

const (
	CommitTypeGenesis CommitType = iota
	CommitTypeSigned
	CommitTypeAnchor
)

const (
	StreamIdEncoding = multibase.Base36
)

var (
	streamIdEncoder = multibase.MustNewEncoder(StreamIdEncoding)
)

type StreamId struct {
	Type       StreamType
	Cid        cid.Cid
	Log        cid.Cid
	GenesisLog bool
}

func (id StreamId) String() string {
	var buf bytes.Buffer
	buf.Write(binary.AppendUvarint(nil, uint64(multicodec.Streamid)))
	buf.Write(binary.AppendUvarint(nil, uint64(id.Type)))
	buf.Write(id.Cid.Bytes())
	if id.Log.ByteLen() != 0 {
		buf.Write(id.Log.Bytes())
	} else if id.GenesisLog {
		buf.Write([]byte{0})
	}
	return streamIdEncoder.Encode(buf.Bytes())
}

func (id StreamId) Genesis() StreamId {
	id.GenesisLog = true
	return id
}

func (id StreamId) With(str string) StreamId {
	id.GenesisLog = false
	id.Log = cid.MustParse(str)
	return id
}

func (id StreamId) GetStream(ctx context.Context) (stream Stream, err error) {
	return GetStreamId(ctx, id)
}

func NewStreamId(t StreamType, cidStr ...string) (id StreamId, err error) {
	if len(cidStr) != 1 && len(cidStr) != 2 {
		err = fmt.Errorf("unexpect cid length, could only 1 or 2 (genesis cid and log cid)")
	}
	id.Type = t

	return
}

func ParseStreamID(str string) (id StreamId, err error) {
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
	if encoding != StreamIdEncoding {
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
	id.Type = StreamType(streamType)
	buf = buf[idx:]

	var nr int
	if nr, id.Cid, err = cid.CidFromBytes(buf); err != nil {
		return
	}
	buf = buf[nr:]
	if len(buf) != 0 {
		if len(buf) == 1 && buf[0] == 0 {
			id.GenesisLog = true
		} else if _, id.Log, err = cid.CidFromBytes(buf); err != nil {
			return
		}
	}
	return
}

var _ json.Marshaler = (*StreamId)(nil)

func (id *StreamId) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

var _ json.Unmarshaler = (*StreamId)(nil)

func (id *StreamId) UnmarshalJSON(src []byte) (err error) {
	var str string
	if err = json.Unmarshal(src, &str); err != nil {
		return
	}
	if *id, err = ParseStreamID(str); err != nil {
		return
	}
	return
}

func (StreamId) GormDataType() string {
	return "text"
}

var _ sql.Scanner = (*StreamId)(nil)

func (id *StreamId) Scan(src any) (err error) {
	str, ok := src.(string)
	if !ok {
		err = fmt.Errorf("cannot parse %s to string", src)
		return
	}
	*id, err = ParseStreamID(str)
	return
}

var _ driver.Valuer = (*StreamId)(nil)

func (id *StreamId) Value() (driver.Value, error) {
	return id.String, nil
}
