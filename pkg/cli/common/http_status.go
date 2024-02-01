package common

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/seal-io/walrus/utils/json"
)

var NotFoundReg = regexp.MustCompile(`model: (.+?) not found`)

type ErrorResponse struct {
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}

func CheckResponseStatus(resp *http.Response) error {
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode == http.StatusConflict:
		return NewRetryableError("conflict")
	case resp.StatusCode == http.StatusTooManyRequests:
		return NewRetryableError("too many request")
	case resp.StatusCode == http.StatusUnauthorized:
		return fmt.Errorf("unauthorized, please check the validity of the token")
	case resp.StatusCode == http.StatusNotFound:
		commonErrMsg := fmt.Sprintf("got not found response from %s", resp.Request.URL.String())

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%s, failed to read response body: %w", commonErrMsg, err)
		}

		var errRes ErrorResponse

		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return fmt.Errorf("%s, failed to decode response body: %w", commonErrMsg, err)
		}

		matches := NotFoundReg.FindStringSubmatch(errRes.Message)

		switch {
		case len(matches) >= 2:
			return fmt.Errorf("%s not found", matches[1])
		case errRes.Message != "":
			return fmt.Errorf("%s, %s", commonErrMsg, errRes.Message)
		default:
			return errors.New(commonErrMsg)
		}
	default:
		commonErrMsg := fmt.Sprintf("unexpected status code %d from %s",
			resp.StatusCode, resp.Request.URL.String())

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%s, failed to read response body: %w", commonErrMsg, err)
		}

		var (
			errMsg = fmt.Errorf("%s, %s", commonErrMsg, string(body))
			errRes ErrorResponse
		)

		err = json.Unmarshal(body, &errRes)
		if err != nil {
			return errMsg
		}

		if errRes.Message != "" {
			errMsg = fmt.Errorf("%s, %s", commonErrMsg, errRes.Message)
		}

		return errMsg
	}
}
