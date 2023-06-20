package metric

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/seal-io/seal/utils/log"
)

// Index returns a http.Handler to process the metrics exporting.
func Index(maxInFlight int, timeout time.Duration) http.Handler {
	opts := promhttp.HandlerOpts{
		ErrorLog:            log.WithName("metrics"),
		ErrorHandling:       promhttp.HTTPErrorOnError,
		Registry:            reg,
		DisableCompression:  false,
		MaxRequestsInFlight: maxInFlight,
		Timeout:             timeout,
		EnableOpenMetrics:   true,
	}

	return promhttp.HandlerFor(reg, opts)
}
