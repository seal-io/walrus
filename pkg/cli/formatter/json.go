package formatter

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/seal-io/walrus/utils/json"
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
