package pkg

// Router 路由
type Router struct {
	Match    RouterMatch `yaml:"match" json:"match"`
	Route    RouteAction `yaml:"route" json:"route"`
	Redirect RouteAction `yaml:"redirect" json:"redirect"`
	//"metadata": "{...}",
	//"decorator": "{...}"
}

type RouterMatch struct {
	Prefix        string          `yaml:"prefix" json:"prefix"`
	Path          string          `yaml:"path" json:"path"`
	Regex         string          `yaml:"regex" json:"regex"`
	CaseSensitive bool            // CaseSensitive default true
	Headers       []HeaderMatcher `yaml:"headers" json:"headers"`
}

type RouteAction struct {
	Cluster                     string            `yaml:"cluster" json:"cluster"` // Cluster cluster name
	ClusterNotFoundResponseCode int               `yaml:"cluster_not_found_response_code" json:"cluster_not_found_response_code"`
	PrefixRewrite               string            `yaml:"prefix_rewrite" json:"prefix_rewrite"`
	HostRewrite                 string            `yaml:"host_rewrite" json:"host_rewrite"`
	Timeout                     string            `yaml:"timeout" json:"timeout"`
	Priority                    int8              `yaml:"priority" json:"priority"`
	ResponseHeadersToAdd        HeaderValueOption `yaml:"response_headers_to_add" json:"response_headers_to_add"`       // ResponseHeadersToAdd add response head
	ResponseHeadersToRemove     []string          `yaml:"response_headers_to_remove" json:"response_headers_to_remove"` // ResponseHeadersToRemove remove response head
	RequestHeadersToAdd         HeaderValueOption `yaml:"request_headers_to_add" json:"request_headers_to_add"`         // RequestHeadersToAdd add request head
	Cors                        CorsPolicy        `yaml:"cors" json:"cors"`
}

// RouteConfiguration
type RouteConfiguration struct {
	InternalOnlyHeaders     []string          `yaml:"internal_only_headers" json:"internal_only_headers"`           // InternalOnlyHeaders used internal, clear http request head
	ResponseHeadersToAdd    HeaderValueOption `yaml:"response_headers_to_add" json:"response_headers_to_add"`       // ResponseHeadersToAdd add response head
	ResponseHeadersToRemove []string          `yaml:"response_headers_to_remove" json:"response_headers_to_remove"` // ResponseHeadersToRemove remove response head
	RequestHeadersToAdd     HeaderValueOption `yaml:"request_headers_to_add" json:"request_headers_to_add"`         // RequestHeadersToAdd add request head
	Routes                  []Router          `yaml:"routes" json:"routes"`
}
