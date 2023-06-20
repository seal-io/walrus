package formatter

import (
	"bytes"
	"io"
	"net/http"

	"github.com/seal-io/seal/utils/json"
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

	var data interface{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")

	err = enc.Encode(data)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
