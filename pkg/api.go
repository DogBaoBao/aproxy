package pkg

import (
	"aproxy/pkg/service"
	"encoding/json"
	"errors"
	"github.com/goinggo/mapstructure"
)

func init() {
	AddApiDiscoveryService(ApiDiscoveryService_Dubbo, NewApiDiscoveryService())
}

// Api
type Api struct {
	Name     string      `json:"name"`
	ITypeStr string      `json:"itype"`
	IType    ApiType     `json:"-"`
	OTypeStr string      `json:"otype"`
	OType    ApiType     `json:"-"`
	Status   Status      `json:"status"`
	Metadata interface{} `json:"metadata"`
	Method   string      `json:"method"`
	RequestMethod
	Client Client
}

type DubboMetadata struct {
	ApplicationName      string   `yaml:"application_name" json:"application_name" mapstructure:"application_name"`
	Group                string   `yaml:"group" json:"group" mapstructure:"group"`
	Version              string   `yaml:"version" json:"version" mapstructure:"version"`
	Interface            string   `yaml:"interface" json:"interface" mapstructure:"interface"`
	Method               string   `yaml:"method" json:"method" mapstructure:"method"`
	Types                []string `yaml:"type" json:"types" mapstructure:"types"`
	Retries              string   `yaml:"retries"  json:"retries,omitempty" property:"retries"`
	ProtocolTypeStr      string   `yaml:"protocol_type"  json:"protocol_type,omitempty" property:"protocol_type"`
	SerializationTypeStr string   `yaml:"serialization_type"  json:"serialization_type,omitempty" property:"serialization_type"`
}

type ApiDiscoveryService struct {
}

func NewApiDiscoveryService() *ApiDiscoveryService {
	return &ApiDiscoveryService{}
}

func (ads *ApiDiscoveryService) AddApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	aj := &Api{}
	if err := json.Unmarshal(request.Body, aj); err != nil {
		return *service.EmptyDiscoveryResponse, err
	}

	apiCache[aj.Name] = aj

	if aj.Metadata == nil {

	} else {
		if v, ok := aj.Metadata.(map[string]interface{}); ok {
			if d, ok := v["dubbo"]; ok {
				dm := &DubboMetadata{}
				if err := mapstructure.Decode(d, dm); err != nil {
					return *service.EmptyDiscoveryResponse, err
				}
				aj.Metadata = dm
			}
		}

		aj.RequestMethod = RequestMethod(RequestMethodValue[aj.Method])
	}

	return *service.NewSuccessDiscoveryResponse(), nil
}

func (ads *ApiDiscoveryService) GetApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	n := string(request.Body)

	if a, ok := apiCache[n]; ok {
		return *service.NewDiscoveryResponse(a), nil
	}

	return *service.EmptyDiscoveryResponse, errors.New("not found")
}

func (a *Api) FindApi(name string) (*Api, bool) {
	ads := GetMustApiDiscoveryService(ApiDiscoveryService_Dubbo)
	if api, err := ads.GetApi(*service.NewDiscoveryRequest([]byte(name))); err != nil {
		return nil, false
	} else {
		return api.Data.(*Api), true
	}
}

func (a *Api) MatchMethod(method string) bool {
	i := RequestMethodValue[method]
	if a.RequestMethod == RequestMethod(i) {
		return true
	}

	return false
}

func (a *Api) IsOk() bool {
	return a.Status == Up
}

// RegisterSelf register api self
func (a *Api) RegisterSelf() {
	apiCache[a.Name] = a
}

// Register register a new api
func (a *Api) Register(api *Api) {
	apiCache[api.Name] = api
}

// Offline api offline
func (a *Api) Offline(name string) {
	if v, ok := apiCache[name]; ok {
		v.Status = Down
	}
}

// Online api online
func (a *Api) Online(name string) {
	if v, ok := apiCache[name]; ok {
		v.Status = Up
	}
}
