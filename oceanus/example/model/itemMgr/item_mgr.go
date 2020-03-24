package itemMgr

import (
	"fmt"
	"github.com/laconiz/eros/database/elastic"
	"github.com/laconiz/eros/database/redis"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/model/db/elasticMgr"
	"github.com/laconiz/eros/oceanus/example/model/db/redisMgr"
	"github.com/laconiz/eros/oceanus/example/proto"
	"github.com/robfig/cron"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type UserID = proto.UserID
type ItemID = proto.ItemID
type Reason = proto.ItemChangeReason

// ---------------------------------------------------------------------------------------------------------------------
// 获取用户的物品列表

func Items(userID UserID) (proto.ItemMap, error) {
	m := map[ItemID]int64{}
	err := rds.Hash(key(userID)).GetAll(&m)
	return m, err
}

// ---------------------------------------------------------------------------------------------------------------------
// 获取指定物品数量

func Num(userID UserID, itemID ItemID) (int64, error) {
	var num int64
	_, err := rds.Hash(key(userID)).Get(itemID, &num)
	return num, err
}

// ---------------------------------------------------------------------------------------------------------------------
// 修改物品数量

func Change(userID UserID, itemID ItemID, value int64, reason Reason) (int64, bool, error) {

	// 修改物品数量
	latest, success, err := rds.Hash(key(userID)).Consume(itemID, value)
	if err != nil {
		return latest, success, err
	}

	// 物品修改日志
	log := &model.ItemChangeLog{
		Item:    itemID,
		User:    userID,
		Value:   value,
		Latest:  latest,
		Reason:  reason,
		Success: success,
		Time:    time.Now(),
	}

	// 记录日志
	if err := elt.Insert(log); err != nil {
		logger.Err(err).Data(log).Error("insert change log error")
	}

	return latest, success, err
}

// ---------------------------------------------------------------------------------------------------------------------
// 生成User的物品表ID

func key(userID UserID) string {
	return fmt.Sprintf("%s.%d", model.ModuleItem, userID)
}

// ---------------------------------------------------------------------------------------------------------------------
// 初始化

var (
	rds *redis.Redis
	elt *elastic.Elastic
)

func init() {

	var err error

	// 连接redis
	rds, err = redisMgr.Connect(model.ModuleItem)
	if err != nil {
		panic(err)
	}

	// 连接elastic
	elt, err = elasticMgr.Connect(model.ModuleItem)
	if err != nil {
		panic(err)
	}

	// TODO 更新物品改变日志索引别名

	// TODO 创建基于日期的新物品改变日志索引
	schedule := cron.New()
	schedule.AddFunc("59 23 * * *", func() {
		// rds.Singleton().Exec()
	})
	schedule.Run()
}

// ---------------------------------------------------------------------------------------------------------------------
// 日志接口

var logger = model.Logger(model.ModuleItem)
