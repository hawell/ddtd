package server

import (
	"bytes"
	"context"
	"ddtd/types"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tinylib/msgp/msgp"
)

type Server struct {
	config     *Config
	httpServer *http.Server
}

var allTraces []types.Trace

func isApiCall(r *http.Request) bool {
	if strings.HasPrefix(r.URL.Path, "/info") || strings.HasPrefix(r.URL.Path, "/ter.URL.Pathlemetry") ||
		strings.HasPrefix(r.URL.Path, "/v0.7") || strings.HasPrefix(r.URL.Path, "/v0.4") ||
		strings.HasPrefix(r.URL.Path, "/api/get_traces") || strings.HasPrefix(r.URL.Path, "/api/clear") {
		return true
	}
	return false
}

func NewServer(config Config) (*Server, error) {
	apiRouter := gin.New()
	handleRecovery := func(c *gin.Context, err interface{}) {
		ErrorResponse(c, http.StatusInternalServerError, err.(string))
		c.Abort()
	}
	bodySizeMiddleware := func(c *gin.Context) {
		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, config.MaxBodyBytes)

		c.Next()
	}
	apiRouter.Use(gin.CustomRecovery(handleRecovery))
	apiRouter.Use(bodySizeMiddleware)

	apiRouter.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiRouter.Any("info", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	apiRouter.Any("telemetry/proxy/api/v2/apmtelemetry", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	apiRouter.Any("v0.7/config", func(ctx *gin.Context) { ctx.Status(http.StatusOK) })
	apiRouter.Any("v0.4/traces", func(ctx *gin.Context) {
		fmt.Println("traces")
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if len(body) != 0 {
			b := bytes.Buffer{}
			_, err = msgp.UnmarshalAsJSON(&b, body)
			if err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			fmt.Println(b.String())
			err = AddTrace(ctx, b.Bytes())
			if err != nil {
				fmt.Println(err.Error())
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		ctx.JSON(http.StatusOK, nil)

	})
	apiRouter.GET("api/get_traces", func(ctx *gin.Context) {
		traces, err := GetTraces(ctx)
		if err != nil {
			fmt.Println(err.Error())
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, traces)
	})
	apiRouter.POST("api/clear", func(ctx *gin.Context) {
		Clear(ctx)
	})

	panelRouter := gin.New()
	panelRouter.Static("/", config.PanelRoot)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isApiCall(r) {
			apiRouter.Handler().ServeHTTP(w, r)
		} else {
			panelRouter.Handler().ServeHTTP(w, r)
		}
	})

	s := &http.Server{
		Addr:           config.BindAddress,
		Handler:        handler,
		ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &Server{
		config:     &config,
		httpServer: s,
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func AddTrace(ctx context.Context, b []byte) error {
	fmt.Println(string(b))
	var tracess [][]types.Trace
	err := json.Unmarshal(b, &tracess)
	if err != nil {
		return err
	}
	for _, traces := range tracess {
		allTraces = append(allTraces, traces...)
	}
	return nil
}

func GetTraces(ctx context.Context) ([]types.Trace, error) {
	m := map[int64]int{}
	for i := range allTraces {
		m[allTraces[i].SpanID] = i
	}
	for i := range allTraces {
		setLevel(m, allTraces, i)
	}
	return allTraces, nil
}

func setLevel(m map[int64]int, l []types.Trace, i int) {
	if l[i].ParentID == 0 {
		l[i].Level = 1
		return
	}
	pindex := m[l[i].ParentID]
	if l[pindex].Level == 0 {
		setLevel(m, l, pindex)
	}
	l[i].Level = l[pindex].Level + 1
}

func Clear(ctx context.Context) {
	allTraces = []types.Trace{}
}
