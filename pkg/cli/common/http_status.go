package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/seal-io/walrus/utils/json"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}

func CheckResponseStatus(resp *http.Response, name string) error {
	switch {
	default:
		return nil
	case resp.StatusCode == http.StatusConflict:
		return NewRetryableError("conflict")
	case resp.StatusCode == http.StatusTooManyRequests:
		return NewRetryableError("too many request")
	case resp.StatusCode == http.StatusUnauthorized:
		return fmt.Errorf("unauthorized, please check the validity of the token")
	case resp.StatusCode == http.StatusNotFound:
		if name == "" {
			return errors.New("not found")
		}

		return fmt.Errorf("%s not found", name)
	case resp.StatusCode < 200 || resp.StatusCode >= 300:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unexpected status code %d, failed to read response body: %w", resp.StatusCode, err)
		}

		var (
			errMsg = fmt.Errorf("unexpected status code %d, %s", resp.StatusCode, string(body))
			errRes ErrorResponse
		)

		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return errMsg
		}

		if errRes.Message != "" {
			errMsg = fmt.Errorf("unexpected status code %d, %s", resp.StatusCode, errRes.Message)
		}

		return errMsg
	}
}
