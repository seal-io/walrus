package httpx

import (
	"net/http"

	"github.com/seal-io/utils/json"
	"github.com/seal-io/utils/pools/bytespool"
)

func PureJSON(w http.ResponseWriter, code int, v any) error {
	buf := bytespool.GetBuffer()
	defer bytespool.Put(buf)

	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(buf.Bytes())
	return err
}

func JSON(w http.ResponseWriter, code int, v any) error {
	buf := bytespool.GetBuffer()
	defer bytespool.Put(buf)

	err := json.NewEncoder(buf).Encode(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = w.Write(buf.Bytes())
	return err
}
