package pkg

import "net/http"

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
	default406Body = []byte("406 api not up")
)

func init() {
	AddFilterFunc(HttpApiFilter, ApiFilter())
}

func ApiFilter() FilterFunc {
	return func(c Context) {
		url := c.GetUrl()
		method := c.GetMethod()

		isHit := false
		a := &Api{}
		if api, b := a.FindApi(url); b {
			if !api.MatchMethod(method) {
				c.WriteWithStatus(http.StatusMethodNotAllowed, default405Body)
			} else {
				if !api.IsOk() {
					c.WriteWithStatus(http.StatusNotAcceptable, default406Body)
				} else {
					isHit = true
					c.Api(api)
					c.Next()
				}
			}
		} else {
			// status must set first
			c.WriteWithStatus(http.StatusNotFound, default404Body)
		}

		if !isHit {
			c.AddHeader(HeaderKeyContextType, HeaderValueTextPlain)
			c.Abort()
		}

	}
}
