package options

import (
	"fmt"
	"sync"

	"github.com/abmpio/abmp/pkg/log"
	"github.com/abmpio/configurationx"
	"github.com/abmpio/configurationx/options/mongodb"
	"github.com/mitchellh/mapstructure"
)

const (
	ConfigurationKey string = "plugins.xapikey"
)

var (
	_options XApiKeyOptions
	_once    sync.Once
)

type XApiKeyOptions struct {
	DatabaseName     string `json:"databaseName"`
	MongodbClientKey string `json:"mongodbClientKey"`

	//路由的前缀,如果加上此值,则所有路由都将加上这个前缀，加上时会自动以 /头,以/结尾
	RouterPrefixPath        string `json:"routerPrefixPath"`
	DisableControllerRegist bool   `json:"disableControllerRegist"`

	DisableXApiKey bool `json:"disableXApiKey"`
}

func (o *XApiKeyOptions) normalize() {
	if len(o.MongodbClientKey) <= 0 {
		o.MongodbClientKey = mongodb.AliasName_Default
	}
	if len(o.DatabaseName) <= 0 {
		o.DatabaseName = configurationx.GetInstance().Options.Mongodb.GetDefaultOptions().Database
	}
}

func GetOptions() *XApiKeyOptions {
	_once.Do(func() {
		if err := configurationx.GetInstance().UnmarshFromKey(ConfigurationKey, &_options, func(dc *mapstructure.DecoderConfig) {
			dc.TagName = "json"
		}); err != nil {
			err = fmt.Errorf("无效的配置文件,%s", err)
			log.Logger.Error(err.Error())
			panic(err)
		}
		_options.normalize()
	})
	return &_options
}
