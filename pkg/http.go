package pkg

type HttpConnectionManager struct {
	RouteConfig       RouteConfiguration `yaml:"route_config" json:"route_config" mapstructure:"route_config"`
	HttpFilters       []HttpFilter       `yaml:"http_filters" json:"http_filters" mapstructure:"http_filters"`
	ServerName        string             `yaml:"server_name" json:"server_name" mapstructure:"server_name"`
	IdleTimeoutStr    string             `yaml:"idle_timeout" json:"idle_timeout" mapstructure:"idle_timeout"`
	AccessLog         AccessLog          `yaml:"access_log" json:"access_log" mapstructure:"access_log"`
	GenerateRequestId bool               `yaml:"generate_request_id" json:"generate_request_id" mapstructure:"generate_request_id"`
}

func DefaultHttpConnectionManager() *HttpConnectionManager {
	return &HttpConnectionManager{
		RouteConfig: RouteConfiguration{
			Routes: []Router{
				{
					Match: RouterMatch{
						Prefix: "/api/v1",
					},
					Route: RouteAction{
						Cluster: "*",
					},
				},
			},
		},
		HttpFilters: []HttpFilter{
			{
				Name: "dgp.filters.http.router",
			},
		},
	}
}

type CorsPolicy struct {
	AllowOrigin      []string `yaml:"allow_origin" json:"allow_origin" mapstructure:"allow_origin"`
	AllowMethods     string   // access-control-allow-methods
	AllowHeaders     string   // access-control-allow-headers
	ExposeHeaders    string   // access-control-expose-headers
	MaxAge           string   // access-control-max-age
	AllowCredentials bool
	Enabled          bool `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
}

type HttpFilter struct {
	Name   string      `yaml:"name" json:"name" mapstructure:"name"`
	Config interface{} `yaml:"config" json:"config" mapstructure:"config"`
}

type RequestMethod int32

const (
	METHOD_UNSPECIFIED = 0 + iota // (DEFAULT)
	GET
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
)

var RequestMethodName = map[int32]string{
	0: "METHOD_UNSPECIFIED",
	1: "GET",
	2: "HEAD",
	3: "POST",
	4: "PUT",
	5: "DELETE",
	6: "CONNECT",
	7: "OPTIONS",
	8: "TRACE",
}

var RequestMethodValue = map[string]int32{
	"METHOD_UNSPECIFIED": 0,
	"GET":                1,
	"HEAD":               2,
	"POST":               3,
	"PUT":                4,
	"DELETE":             5,
	"CONNECT":            6,
	"OPTIONS":            7,
	"TRACE":              8,
}
