package pkg

import (
	"aproxy/pkg/logger"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/goinggo/mapstructure"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	configPath     string
	config         *Bootstrap
	configLoadFunc ConfigLoadFunc = DefaultConfigLoad
)

func GetBootstrap() *Bootstrap {
	return config
}

// Load config file and parse
func Load(path string) *Bootstrap {
	logger.Infof("[dubboproxy go] load path:%s", path)

	configPath, _ = filepath.Abs(path)
	if yamlFormat(path) {
		RegisterConfigLoadFunc(YAMLConfigLoad)
	}
	if cfg := configLoadFunc(path); cfg != nil {
		config = cfg
	}

	return config
}

// ConfigLoadFunc parse a input(usually file path) into a proxy config
type ConfigLoadFunc func(path string) *Bootstrap

// RegisterConfigLoadFunc can replace a new config load function instead of default
func RegisterConfigLoadFunc(f ConfigLoadFunc) {
	configLoadFunc = f
}

func yamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".yaml" || ext == ".yml" {
		return true
	}
	return false
}

func YAMLConfigLoad(path string) *Bootstrap {
	log.Println("load config in YAML format from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [yaml load] load config failed, ", err)
	}
	cfg := &Bootstrap{}

	bytes, err := yaml.YAMLToJSON(content)
	if err != nil {
		log.Fatalln("[config] [yaml load] convert YAML to JSON failed, ", err)
	}

	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] yaml unmarshal config failed, ", err)
	}

	// other adapter

	for i, l := range cfg.StaticResources.Listeners {
		if l.Address.SocketAddress.ProtocolStr == "" {
			l.Address.SocketAddress.ProtocolStr = "HTTP"
		}
		l.Address.SocketAddress.Protocol = ProtocolType(ProtocolTypeValue[l.Address.SocketAddress.ProtocolStr])

		hc := &HttpConfig{}
		if l.Config != nil {
			if v, ok := l.Config.(map[string]interface{}); ok {
				switch l.Name {
				case "net/http":
					if err := mapstructure.Decode(v, hc); err != nil {
						logger.Error(err)
					}

					cfg.StaticResources.Listeners[i].Config = hc
				}
			}
		}

		for _, fc := range l.FilterChains {
			if fc.Filters != nil {
				for i, fcf := range fc.Filters {
					hcm := &HttpConnectionManager{}
					if fcf.Config != nil {
						switch fcf.Name {
						case "dgp.filters.http_connect_manager":
							if v, ok := fcf.Config.(map[string]interface{}); ok {
								if err := mapstructure.Decode(v, hcm); err != nil {
									logger.Error(err)
								}

								fc.Filters[i].Config = hcm
							}
						}
					}
				}
			}
		}

	}

	for _, c := range cfg.StaticResources.Clusters {
		var discoverType int32
		if c.TypeStr != "" {
			if t, ok := DiscoveryTypeValue[c.TypeStr]; ok {
				discoverType = t
			} else {
				c.TypeStr = "EDS"
				discoverType = DiscoveryTypeValue[c.TypeStr]
			}
		} else {
			c.TypeStr = "EDS"
			discoverType = DiscoveryTypeValue[c.TypeStr]
		}
		c.Type = DiscoveryType(discoverType)

		var lbPolicy int32
		if c.LbStr != "" {
			if lb, ok := LbPolicyValue[c.LbStr]; ok {
				lbPolicy = lb
			} else {
				c.LbStr = "RoundRobin"
				lbPolicy = LbPolicyValue[c.LbStr]
			}
		} else {
			c.LbStr = "RoundRobin"
			lbPolicy = LbPolicyValue[c.LbStr]
		}
		c.Lb = LbPolicy(lbPolicy)
	}

	return cfg
}

func DefaultConfigLoad(path string) *Bootstrap {
	log.Println("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [default load] load config failed, ", err)
	}
	cfg := &Bootstrap{}
	// translate to lower case
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalln("[config] [default load] json unmarshal config failed, ", err)
	}
	return cfg

}
