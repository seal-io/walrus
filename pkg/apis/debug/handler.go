package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/seal-io/walrus/pkg/apis/runtime"
	"github.com/seal-io/walrus/utils/errorx"
	"github.com/seal-io/walrus/utils/log"
	"github.com/seal-io/walrus/utils/version"
)

func Version() runtime.Handle {
	info := gin.H{
		"version": version.Version,
		"commit":  version.GitCommit,
	}

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, info)
	}
}

func PProf() runtime.HTTPHandler {
	// NB(thxCode): init from net/http/pprof.
	m := http.NewServeMux()
	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	return m
}

func SetFlags() runtime.ErrorHandle {
	return func(ctx *gin.Context) error {
		// Support set flags log-debug and log-verbosity.
		var input struct {
			LogDebug     *bool   `query:"log-debug"`
			LogVerbosity *uint64 `query:"log-verbosity"`
		}

		if err := binding.MapFormWithTag(&input, ctx.Request.URL.Query(), "query"); err != nil {
			return errorx.WrapHttpError(http.StatusBadRequest, err, "invalid query params")
		}

		resp := map[string]any{}

		if input.LogDebug != nil {
			level := log.InfoLevel
			if *input.LogDebug {
				level = log.DebugLevel
			}

			log.SetLevel(level)
			resp["log-debug"] = *input.LogDebug
		}

		if input.LogVerbosity != nil {
			log.SetVerbosity(*input.LogVerbosity)
			resp["log-verbosity"] = *input.LogVerbosity
		}

		ctx.JSON(http.StatusOK, resp)

		return nil
	}
}

func GetFlags() runtime.ErrorHandle {
	return func(ctx *gin.Context) error {
		resp := map[string]any{
			"log-debug":     log.GetLevel() == log.DebugLevel,
			"log-verbosity": log.GetVerbosity(),
		}

		ctx.JSON(http.StatusOK, resp)

		return nil
	}
}
