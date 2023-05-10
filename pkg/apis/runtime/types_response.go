package runtime

import (
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetResponseCollection(c *gin.Context, data any, totalSize int) ResponseCollection {
	req := struct {
		RequestPagination `query:",inline"`
		RequestGrouping   `query:",inline"`
	}{
		RequestPagination: RequestPagination{
			Page: -1,
		},
	}
	_ = binding.MapFormWithTag(&req, c.Request.URL.Query(), "query")

	var currentSize int
	dataRef := reflect.ValueOf(data)
	if dataRef.Kind() == reflect.Slice {
		currentSize = dataRef.Len()
	}

	var totalPage int
	if req.Page > 0 {
		if req.PerPage <= 0 {
			req.PerPage = 100
		}
		totalPage = int(math.Ceil(float64(totalSize) / float64(req.PerPage)))
	} else {
		req.Page = 1
		req.PerPage = totalSize
		totalPage = 1
	}

	partial := currentSize < totalSize
	group := len(req.Groups) != 0

	nextPage := req.Page + 1
	if !partial || nextPage > totalPage {
		nextPage = 0
	}

	return ResponseCollection{
		Items: data,
		Pagination: ResponsePagination{
			Page:      req.Page,
			PerPage:   req.PerPage,
			Total:     totalSize,
			TotalPage: totalPage,
			Partial:   partial,
			Group:     group,
			NextPage:  nextPage,
		},
	}
}

type ResponseCollection struct {
	Items      any                `json:"items"`
	Pagination ResponsePagination `json:"pagination"`
}

type ResponsePagination struct {
	Page      int  `json:"page"`
	PerPage   int  `json:"perPage"`
	Total     int  `json:"total"`
	TotalPage int  `json:"totalPage"`
	Partial   bool `json:"partial"`
	Group     bool `json:"group,omitempty"`
	NextPage  int  `json:"nextPage,omitempty"`
}

// ResponseStream is similar to render.Reader,
// but be able to close the reader out of the handler processing.
type ResponseStream struct {
	ContentType   string
	ContentLength int64
	Headers       map[string]string
	Reader        io.ReadCloser
}

func (r ResponseStream) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if r.ContentLength > 0 {
		if r.Headers == nil {
			r.Headers = map[string]string{}
		}
		r.Headers["Content-Length"] = strconv.FormatInt(r.ContentLength, 10)
	}
	header := w.Header()
	for k, v := range r.Headers {
		if header.Get(k) == "" {
			header.Set(k, v)
		}
	}
	_, err = io.Copy(w, r.Reader)
	return
}

func (r ResponseStream) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if vs := header["Content-Type"]; len(vs) == 0 {
		contentType := "application/octet-stream"
		if r.ContentType != "" {
			contentType = r.ContentType
		}
		header["Content-Type"] = []string{contentType}
	}
}

func (r ResponseStream) Close() error {
	if r.Reader == nil {
		return nil
	}
	return r.Reader.Close()
}
