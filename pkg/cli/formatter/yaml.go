package formatter

import (
	"io"
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/seal-io/seal/utils/json"
)

// YamlFormatter use to convert response to yaml format.
type YamlFormatter struct{}

func (f *YamlFormatter) Format(resp *http.Response) ([]byte, error) {
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

	b, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return b, nil
}
