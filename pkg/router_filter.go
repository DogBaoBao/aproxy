package pkg

import (
	"net/http"
	"strings"
)

var (
	default403Body = []byte("403 for bidden")
)

func init() {
	AddFilterFunc(HttpRouterFilter, HttpRouting())
}

// HttpRouting http router filter
func HttpRouting() FilterFunc {
	return func(c Context) {
		routingFilter(c.(*HttpContext))
	}
}

// routingFilter
func routingFilter(c *HttpContext) {
	result := true
	for _, v := range c.httpConnectionManager.RouteConfig.Routes {
		result = routeMatch(c, v)
		if result {
			httpHeaderCorsHandler(c, v)
			break
		}
	}

	if !result {
		c.WriteWithStatus(http.StatusForbidden, default403Body)
		c.Abort()
	}
}

// routeMatch will match router with request, only true or false way
func routeMatch(c *HttpContext, r Router) bool {
	result := true
	if len(r.Match.Headers) > 0 {
		for _, v := range r.Match.Headers {
			result = HttpHeaderMatch(c, v)
			if !result {
				break
			}
		}
	}

	if !result {
		return result
	}

	result = HttpRouteMatch(c, r.Match)

	if !result {
		return result
	}

	return HttpRouteActionMatch(c, r.Route)
}

// httpHeaderCorsHandler will set cors, handler mean can do c.Next()
func httpHeaderCorsHandler(c *HttpContext, r Router) {
	var acao string
	if r.Route.Cors.Enabled {
		acao = strings.Join(r.Route.Cors.AllowOrigin, "|")
	}

	c.Next()

	if acao != "" {
		c.AddHeader(HeaderKeyAccessControlAllowOrigin, acao)
	}
}
