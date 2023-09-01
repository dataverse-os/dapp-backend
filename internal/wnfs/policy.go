package wnfs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/dataverse-os/dapp-backend/ceramic"
	"github.com/tidwall/gjson"
)

var _ ceramic.CommitVerifier = (*IndexFileVerifier)(nil)

type IndexFileVerifier struct{}

func (*IndexFileVerifier) ProtectedFields() []string {
	return []string{
		"contentId", "contentType",
	}
}

func (v *IndexFileVerifier) ValidateData(data json.RawMessage) (err error) {
	obj := gjson.ParseBytes(data)
	if err = v.validateStreamId(obj.Get("contentId").String()); err != nil {
		return
	}
	// check acl
	if err = v.validateACL(obj.Get("accessControl").String()); err != nil {
		return
	}
	return
}

func (v *IndexFileVerifier) ValidatePatches(patch ceramic.Patch) (err error) {
	if patch.Path == "/accessControl" {
		if aclStr, ok := patch.Value.(string); ok {
			if err = v.validateACL(aclStr); err != nil {
				return
			}
		} else {
			err = fmt.Errorf("cannot parse accessControl field value: %s as string", patch.Value)
		}
	}
	return
}

func (*IndexFileVerifier) validateStreamId(streamIdStr string) (err error) {
	var streamId ceramic.StreamId
	if streamId, err = ceramic.ParseStreamID(streamIdStr); err != nil {
		return
	}
	_ = streamId
	// TODO check streamId not fs stream
	// TODO check streamId is Dapp stream
	// TODO check streamId can get from ceramic
	return
}

func (*IndexFileVerifier) validateACL(data string) (err error) {
	var acl []byte
	if acl, err = base64.RawStdEncoding.DecodeString(data); err != nil {
		return
	}
	obj := gjson.ParseBytes(acl)
	obj.Get("encryptionProvider.decryptionConditions").ForEach(func(_, value gjson.Result) bool {
		// TODO: check model is dapp
		return true
	})
	return
}

func CheckPolicy(streamContents []json.RawMessage) (err error) {
	return
}
