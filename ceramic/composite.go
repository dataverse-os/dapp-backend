package ceramic

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func ExtarctStreamID(composite string) (streamID string, err error) {
	for k := range gjson.Get(composite, "models").Map() {
		streamID = k
		return
	}
	err = fmt.Errorf("streamID not found")
	return
}
