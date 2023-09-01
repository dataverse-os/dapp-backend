package ceramic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/tidwall/gjson"
)

type MessageType uint64

const (
	MessageTypeUpdate MessageType = iota
	MessageTypeQuery
	MessageTypeResponse
	MessageTypeKeepalive
)

type PubSubUpdateMessage struct {
	Type     MessageType `json:"typ"`
	StreamId StreamId    `json:"stream"`
	Tip      cid.Cid     `json:"tip"`   // the CID of the latest commit
	Model    *StreamId   `json:"model"` // optional
}

type PubSubQueryMessage struct {
	Type     MessageType `json:"typ"`
	ID       string      `json:"id"` //query id
	StreamId StreamId    `json:"stream"`
}

type PubSubResponseMessage struct {
	Type MessageType          `json:"typ"`
	ID   string               `json:"id"` //query id
	Tips map[StreamId]cid.Cid `json:"tips"`
}

// All nodes will always ignore this message
type PubSubKeepaliveMessage struct {
	Type        MessageType `json:"typ"`
	Timestamp   uint64      `json:"ts"`      // current time in milliseconds since epoch
	Version     string      `json:"version"` // current ceramic version
	IPFSVersion string      `json:"ipfsVer"`
}

func (impl IpfsImpl) QueryStream(ctx context.Context, streamId StreamId) (tip cid.Cid, err error) {
	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(map[string]any{
		"typ":    MessageTypeQuery,
		"stream": streamId,
	}); err != nil {
		return
	}
	queryId := MessageHash(buf.Bytes())
	fmt.Println(queryId)
	queryMsg := PubSubQueryMessage{
		Type:     MessageTypeQuery,
		ID:       queryId,
		StreamId: streamId,
	}
	buf = bytes.Buffer{}
	if err = json.NewEncoder(&buf).Encode(queryMsg); err != nil {
		return
	}
	fmt.Println(buf.String())
	if err = impl.pubSubAPI.Publish(ctx, impl.network, buf.Bytes()); err != nil {
		return
	}
	var sub iface.PubSubSubscription
	if sub, err = impl.pubSubAPI.Subscribe(ctx, impl.network); err != nil {
		return
	}
	defer sub.Close()
	for {
		if msg, err := sub.Next(ctx); err != nil {
			log.Println(err)
		} else {
			if typ := gjson.GetBytes(msg.Data(), "typ"); typ.Exists() && typ.Type == gjson.Number && typ.Num == float64(MessageTypeResponse) {
				fmt.Println(string(msg.Data()))
				if gjson.GetBytes(msg.Data(), "id").String() == queryId {
					fmt.Println("hit", string(msg.Data()))
					sub.Close()
					break
				}
			}
		}
	}
	return
}

type MessageHander func(ctx context.Context, msg iface.PubSubMessage) error

func SubscirbeMessage(ctx context.Context, node *rpc.HttpApi, network string, handlers map[MessageType]MessageHander) (err error) {
	var sub iface.PubSubSubscription
	if sub, err = node.PubSub().Subscribe(ctx, network); err != nil {
		return
	}
	defer sub.Close()
	for {
		if msg, err := sub.Next(ctx); err != nil {
			log.Println(err)
		} else {
			if typ := gjson.GetBytes(msg.Data(), "typ"); typ.Exists() && typ.Type == gjson.Number {
				if h, exists := handlers[MessageType(typ.Num)]; exists {
					if err := h(ctx, msg); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}
