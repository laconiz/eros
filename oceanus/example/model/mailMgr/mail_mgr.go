package mailMgr

import (
	"github.com/jinzhu/gorm"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/model/db/postgresMgr"
)

func New() {

}

func NewGlobal() {

}

// ---------------------------------------------------------------------------------------------------------------------

var (
	pgl *gorm.DB
)

func init() {

	var err error

	pgl, err = postgresMgr.Connect(model.ModuleMail)
	if err != nil {
		panic(err)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 日志接口

var logger = model.Logger(model.ModuleMail)
