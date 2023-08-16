package runtime

import (
	"io"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/walrus/pkg/apis/runtime/bind"
)

// NoPageResponse returns the given data without pagination.
func NoPageResponse(data any) *ResponseCollection {
	return &ResponseCollection{
		Items: data,
	}
}

// PageResponse returns the given data in a pagination,
// which calculates the pagination info with the given request page and perPage.
func PageResponse(page, perPage int, data any, dataTotalSize int) *ResponseCollection {
	var currentSize int

	dataRef := reflect.ValueOf(data)
	if dataRef.Kind() == reflect.Slice {
		currentSize = dataRef.Len()
	}

	var totalPage int

	if page > 0 {
		if perPage <= 0 {
			perPage = 100
		}
		totalPage = int(math.Ceil(float64(dataTotalSize) / float64(perPage)))
	} else {
		page = 1
		perPage = dataTotalSize
		totalPage = 1
	}

	partial := currentSize < dataTotalSize

	nextPage := page + 1
	if !partial || nextPage > totalPage {
		nextPage = 0
	}

	return &ResponseCollection{
		Items: data,
		Pagination: &ResponsePagination{
			Page:      page,
			PerPage:   perPage,
			Total:     dataTotalSize,
			TotalPage: totalPage,
			Partial:   partial,
			NextPage:  nextPage,
		},
	}
}

// FullPageResponse returns the given data in a pagination,
// which treats the given data as a full page.
func FullPageResponse(data any, dataTotalSize int) *ResponseCollection {
	return PageResponse(-1, 0, data, dataTotalSize)
}

// TypedResponse returns the given data in typed.
func TypedResponse(typ string, data any) *ResponseCollection {
	return &ResponseCollection{
		Type:  typ,
		Items: data,
	}
}

// getPageResponse gains the request pagination from request,
// and returns the given data in a pagination.
func getPageResponse(c *gin.Context, data any, dataTotalSize int) *ResponseCollection {
	reqPagination := RequestPagination{
		Page: -1,
	}
	bind.WithQuery(c, &reqPagination)

	return PageResponse(reqPagination.Page, reqPagination.PerPage, data, dataTotalSize)
}

// ResponseCollection holds the response data of collection with a pagination.
type ResponseCollection struct {
	Type       string              `json:"type,omitempty"`
	Items      any                 `json:"items"`
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

// ResponsePagination holds the pagination data.
type ResponsePagination struct {
	Page      int  `json:"page"`
	PerPage   int  `json:"perPage"`
	Total     int  `json:"total"`
	TotalPage int  `json:"totalPage"`
	Partial   bool `json:"partial"`
	NextPage  int  `json:"nextPage,omitempty"`
}

// ResponseFile is similar to render.Reader,
// but be able to close the file reader out of the handler processing.
type ResponseFile struct {
	ContentType   string
	ContentLength int64
	Headers       map[string]string
	Reader        io.ReadCloser
}

func (r ResponseFile) Render(w http.ResponseWriter) (err error) {
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

func (r ResponseFile) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if vs := header["Content-Type"]; len(vs) == 0 {
		contentType := "application/octet-stream"
		if r.ContentType != "" {
			contentType = r.ContentType
		}
		header["Content-Type"] = []string{contentType}
	}
}

func (r ResponseFile) Close() error {
	if r.Reader == nil {
		return nil
	}

	return r.Reader.Close()
}
