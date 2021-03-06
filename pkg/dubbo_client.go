package pkg

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

import (
	"aproxy/pkg/logger"
)

import (
	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
)

// TODO java class name elem
const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
)

var (
	_DubboClient *DubboClient
	onceClient   = sync.Once{}
	dgCfg        dg.ConsumerConfig
)

type DubboClient struct {
	mLock              sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

func SingleDubboClient() *DubboClient {
	if _DubboClient == nil {
		onceClient.Do(func() {
			_DubboClient = NewDubboClient()
		})
	}

	return _DubboClient
}

func NewDubboClient() *DubboClient {
	return &DubboClient{
		mLock:              sync.RWMutex{},
		GenericServicePool: make(map[string]*dg.GenericService),
	}
}

func (dc *DubboClient) Init() error {
	dgCfg = dg.GetConsumerConfig()
	dg.SetConsumerConfig(dgCfg)
	dg.Load()
	dc.GenericServicePool = make(map[string]*dg.GenericService)
	return nil
}

func (dc *DubboClient) Close() error {
	return nil
}

func (dc *DubboClient) Call(r *Request) (resp Response, err error) {
	dm := r.Api.Metadata.(*DubboMetadata)
	gs := dc.Get(dm.Interface, dm.Version, dm.Group, dm)

	var reqData []interface{}

	l := len(dm.Types)
	switch {
	case l == 1:
		t := dm.Types[0]
		switch t {
		case JavaStringClassName:
			var s string
			if err := json.Unmarshal(r.Body, &s); err != nil {
				logger.Errorf("params parse error:%+v", err)
			} else {
				reqData = append(reqData, s)
			}
		case JavaLangClassName:
			var i int
			if err := json.Unmarshal(r.Body, &i); err != nil {
				logger.Errorf("params parse error:%+v", err)
			} else {
				reqData = append(reqData, i)
			}
		default:
			bodyMap := make(map[string]interface{})
			if err := json.Unmarshal(r.Body, &bodyMap); err != nil {
				return *EmptyResponse, err
			} else {
				reqData = append(reqData, bodyMap)
			}
		}
	case l > 1:
		if err = json.Unmarshal(r.Body, &reqData); err != nil {
			return *EmptyResponse, err
		}
	}

	logger.Debugf("[dubbogo proxy] invoke, method:%v, types:%v, reqData:%v", dm.Method, dm.Types, reqData)

	if resp, err := gs.Invoke(context.Background(), []interface{}{dm.Method, dm.Types, reqData}); err != nil {
		return *EmptyResponse, err
	} else {
		logger.Debugf("[dubbogo proxy] dubbo client resp:%v", resp)
		return *NewResponse(resp), nil
	}
}

func (dc *DubboClient) get(key string) *dg.GenericService {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	return dc.GenericServicePool[key]
}

func (dc *DubboClient) check(key string) bool {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	if _, ok := dc.GenericServicePool[key]; ok {
		return true
	} else {
		return false
	}
}

func (dc *DubboClient) create(interfaceName, version, group string, dm *DubboMetadata) *dg.GenericService {
	key := strings.Join([]string{interfaceName, version, group}, "_")
	referenceConfig := dg.NewReferenceConfig(interfaceName, context.TODO())
	referenceConfig.InterfaceName = interfaceName
	referenceConfig.Cluster = constant.DEFAULT_CLUSTER
	var registers []string
	for k := range dgCfg.Registries {
		registers = append(registers, k)
	}
	referenceConfig.Registry = strings.Join(registers, ",")

	if dm.ProtocolTypeStr == "" {
		referenceConfig.Protocol = dubbo.DUBBO
	} else {
		referenceConfig.Protocol = dm.ProtocolTypeStr
	}

	referenceConfig.Version = version
	referenceConfig.Group = group
	referenceConfig.Generic = true
	if dm.Retries == "" {
		referenceConfig.Retries = "3"
	} else {
		referenceConfig.Retries = dm.Retries
	}
	dc.mLock.Lock()
	defer dc.mLock.Unlock()
	referenceConfig.GenericLoad(interfaceName)
	time.Sleep(200 * time.Millisecond) //sleep to wait invoker create
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	dc.GenericServicePool[key] = clientService
	return clientService
}

func (dc *DubboClient) Get(interfaceName, version, group string, dm *DubboMetadata) *dg.GenericService {
	key := strings.Join([]string{interfaceName, version, group}, "_")
	if dc.check(key) {
		return dc.get(key)
	} else {
		return dc.create(interfaceName, version, group, dm)
	}
}
