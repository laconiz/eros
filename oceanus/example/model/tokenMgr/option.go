package tokenMgr

import (
	"fmt"
	"github.com/laconiz/eros/database/consul/consulor"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/utils/crypto/des"
)

var option = &struct {
	DesKey  string `json:"desKey"`  // 加密KEY
	Expired int64  `json:"expired"` // 过期时间(秒)
}{}

var cipher *des.DES

func init() {

	err := consulor.KV().Load(model.OptionPrefix+model.ModuleToken, option)
	if err != nil {
		panic(fmt.Errorf("get option error: %w", err))
	}
	logger.Data(option).Info("option")

	if option.Expired <= 0 {
		panic(fmt.Errorf("invalid expired seconds: %v", option.Expired))
	}

	if cipher, err = des.New([]byte(option.DesKey)); err != nil {
		panic(fmt.Errorf("invalid des key: %w", err))
	}
}
