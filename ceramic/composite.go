package ceramic

import (
	"fmt"
	"time"

	"github.com/tidwall/gjson"
	"gorm.io/datatypes"
)

type Composite struct {
	StreamID        string         `gorm:"column:stream_id"`
	ControllerDID   string         `gorm:"column:controller_did"`
	StreamContent   datatypes.JSON `gorm:"column:stream_content"`
	Tip             string         `gorm:"column:tip"`
	LastAnchoredAt  time.Time      `gorm:"column:last_anchored_at"`
	FirstAnchoredAt time.Time      `gorm:"column:first_anchored_at"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
}

func ExtarctStreamID(composite string) (streamID string, err error) {
	for k := range gjson.Get(composite, "models").Map() {
		streamID = k
		return
	}
	err = fmt.Errorf("streamID not found")
	return
}
