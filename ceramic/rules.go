package ceramic

import "encoding/json"

type CommitVerifier interface {
	ProtectedFields() []string
	ValidateData(data json.RawMessage) (err error)
	ValidatePatches(patche Patch) (err error)
}

func ValidateData(data json.RawMessage, rules []CommitVerifier) (err error) {
	for _, rule := range rules {
		if err = rule.ValidateData(data); err != nil {
			return
		}
	}
	return
}

func ValidatePatches(patches []Patch, rules []CommitVerifier) (err error) {
	for _, rule := range rules {
		for _, patch := range patches {
			if err = rule.ValidatePatches(patch); err != nil {
				return
			}
		}
	}
	return
}
