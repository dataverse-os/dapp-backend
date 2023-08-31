package ceramic

type Collection struct {
	Edges []struct {
		Cursor string      `json:"cursor"`
		Node   StreamState `json:"node"`
	} `json:"edges"`
}
