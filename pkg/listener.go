package pkg

import (
	"aproxy/pkg/logger"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Listener is a server, listener a port
type Listener struct {
	Name    string  `yaml:"name" json:"name" mapstructure:"name"`
	Address Address `yaml:"address" json:"address" mapstructure:"address"`

	FilterChains []FilterChain `yaml:"filter_chains" json:"filter_chains" mapstructure:"filter_chains"`
	Config       interface{}   `yaml:"config" json:"config" mapstructure:"config"`
}

func (l *Listener) Start() {
	switch l.Address.SocketAddress.Protocol {
	case HTTP:
		l.httpListener()
	default:
		panic("un support protocol start: " + l.Address.SocketAddress.ProtocolStr)
	}
}

func (l *Listener) httpListener() {
	hl := &DefaultHttpListener{
		pool: sync.Pool{},
	}
	hl.pool.New = func() interface{} {
		return l.allocateContext()
	}

	var hc HttpConfig
	if l.Config == "" {

	} else {
		if c, ok := l.Config.(HttpConfig); ok {
			hc = c
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	srv := http.Server{
		Addr:           resolveAddress(l.Address.SocketAddress.Address + ":" + strconv.Itoa(l.Address.SocketAddress.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubboproxy go] httpListener start by config : %+v", l)

	log.Println(srv.ListenAndServe())
}

func (l *Listener) allocateContext() *HttpContext {
	return &HttpContext{
		l:                     l,
		FilterChains:          l.FilterChains,
		httpConnectionManager: l.findHttpManager(),
		BaseContext:           &BaseContext{},
	}
}

func (l *Listener) findHttpManager() HttpConnectionManager {
	for _, fc := range l.FilterChains {
		for _, f := range fc.Filters {
			if f.Name == HttpConnectManagerFilter {
				return *f.Config.(*HttpConnectionManager)
			}
		}
	}

	return *DefaultHttpConnectionManager()
}

type HttpConfig struct {
	IdleTimeoutStr  string `yaml:"idle_timeout" json:"idle_timeout" mapstructure:"idle_timeout"`
	ReadTimeoutStr  string `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty" mapstructure:"read_timeout"`
	WriteTimeoutStr string `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty" mapstructure:"write_timeout"`
	MaxHeaderBytes  int    `json:"max_header_bytes,omitempty" yaml:"max_header_bytes,omitempty" mapstructure:"max_header_bytes"`
}

type DefaultHttpListener struct {
	pool sync.Pool
}

func (s *DefaultHttpListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hc := s.pool.Get().(*HttpContext)
	hc.r = r
	hc.writermem.reset(w)
	hc.reset()

	hc.AppendFilterFunc(Logger(), Recover())

	hc.BuildFilters()

	s.handleHTTPRequest(hc)

	s.pool.Put(hc)
}

func (s *DefaultHttpListener) handleHTTPRequest(c *HttpContext) {
	if len(c.BaseContext.filters) > 0 {
		c.Next()
		c.writermem.WriteHeaderNow()
		return
	}

	// TODO redirect
}

func resolveInt2IntProp(currentV, defaultV int) int {
	if currentV == 0 {
		return defaultV
	}

	return currentV
}

func resolveStr2Time(currentV string, defaultV time.Duration) time.Duration {
	if currentV == "" {
		return defaultV
	} else {
		if duration, err := time.ParseDuration(currentV); err != nil {
			return 20 * time.Second
		} else {
			return duration
		}
	}
}

func resolveAddress(addr string) string {
	if addr == "" {
		logger.Debug("Addr is undefined. Using port :8080 by default")
		return ":8080"
	}

	return addr
}
