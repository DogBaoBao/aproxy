package pkg

type Request struct {
	Body   []byte
	Header map[string]string
	Api    *Api
}

func NewRequest(b []byte, api *Api) *Request {
	return &Request{
		Body: b,
		Api:  api,
	}
}

type Response struct {
}

var EmptyResponse = &Response{}

type Client interface {
	Init() error
	Close() error

	Call(req *Request) (resp Response, err error)
}

type Endpoint struct {
	Address Address `yaml:"address"`
}
