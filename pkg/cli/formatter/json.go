package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JsonFormatter use to convert response to json format.
type JsonFormatter struct{}

func (f *JsonFormatter) Format(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)

	defer func() { _ = resp.Body.Close() }()

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, nil
	}

	var data any

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer

	// TODO(michelia): Normally we use github.com/seal-io/walrus/utils/json as our json util,
	// since the upstream json-iterator has a bug with MarshalIndent, we use the standard library here,
	// upgrade to github.com/seal-io/walrus/utils/json when the bug is fixed.
	// https://github.com/json-iterator/go/issues/645
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	err = enc.Encode(data)
	if err != nil {
		return nil, err
	}

	// Response status is not 200.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, b.String())
	}

	return b.Bytes(), nil
}
