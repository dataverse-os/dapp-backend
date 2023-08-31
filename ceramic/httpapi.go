package ceramic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type HttpApi struct{}

func (HttpApi) GetStreamId(ctx context.Context, streamId StreamId) (stream Stream, err error) {
	url := fmt.Sprintf("%s/api/v0/streams/%s", os.Getenv("CERAMIC_URL"), streamId)
	var req *http.Request
	var resp *http.Response
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	if err = json.NewDecoder(resp.Body).Decode(&stream); err != nil {
		return
	}
	return
}
