package pkg

import (
	"aproxy/pkg/logger"
	"io/ioutil"
)

func init() {
	AddFilterFunc(HttpTransferDubboFilter, HttpDubbo())
}

func HttpDubbo() FilterFunc {
	return func(c Context) {
		doDubbo(c.(*HttpContext))
	}
}

func doDubbo(c *HttpContext) {
	api := c.GetApi()

	if bytes, err := ioutil.ReadAll(c.r.Body); err != nil {
		logger.Errorf("[dubboproxy go] read body err:%v!", err)
		c.WriteFail()
		c.Abort()
	} else {
		if api.Client == nil {
			api.Client = SingleDubboClient()
		}

		if _, err := api.Client.Call(NewRequest(bytes, api)); err != nil {
			logger.Errorf("[dubboproxy go] client do err:%v!", err)
			c.WriteFail()
			c.Abort()
		} else {
			c.Next()
		}
	}
}
