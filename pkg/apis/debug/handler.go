package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"

	"github.com/seal-io/seal/pkg/apis/runtime"
	"github.com/seal-io/seal/utils/version"
)

func Version() runtime.Handle {
	var info = gin.H{
		"version": version.Version,
		"commit":  version.GitCommit,
	}
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, info)
	}
}

func PProf() runtime.HTTPHandler {
	// NB(thxCode): init from net/http/pprof
	var m = http.NewServeMux()
	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	return m
}
