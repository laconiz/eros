package userMgr

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/laconiz/eros/database/redis"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/model/db/postgresMgr"
	"github.com/laconiz/eros/oceanus/example/model/db/redisMgr"
	"github.com/laconiz/eros/utils/mathe"
	"math/rand"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

type UserID = model.UserID

// ---------------------------------------------------------------------------------------------------------------------
// 生成一个新用户

func New(info *model.User) (*model.User, error) {

	value, err := rds.Key().Incr(key, random.Int63n(step)+1)
	if err != nil {
		return nil, fmt.Errorf("new user id error: %w", err)
	}

	user := &model.User{
		UserID:   UserID(value),
		Name:     info.Name,
		Avatar:   info.Avatar,
		Gender:   info.Gender,
		Phone:    info.Phone,
		Password: info.Password,
	}

	err = pgl.Create(user).Error
	return user, err
}

// ---------------------------------------------------------------------------------------------------------------------
// 根据指定的用户ID获取用户

func Get(userID UserID) (*model.User, error) {

	user := &model.User{}
	err := pgl.First(user, "WHERE user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, model.RecordNotFound
	}

	return user, err
}

// ---------------------------------------------------------------------------------------------------------------------
// 根据指定的用户ID列表获取用户组

func Gets(list []UserID) ([]*model.User, error) {
	var users []*model.User
	err := pgl.Find(&users, "WHERE user_id IN (?)", list).Error
	return users, err
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	pgl    *gorm.DB
	rds    *redis.Redis
	random *rand.Rand
)

const (
	name = "user"
	key  = "user.MaxID"
	base = 10000000 // 基础用户ID
	step = 10       // 用户ID随机步数
)

func init() {

	var err error

	// 连接redis
	rds, err = redisMgr.Connect(name)
	if err != nil {
		panic(err)
	}

	// 连接sql
	pgl, err = postgresMgr.Connect(name)
	if err != nil {
		panic(err)
	}

	// 同步表结构
	pgl.AutoMigrate(&model.User{})

	// 获取redis记录
	exist, err := rds.Key().Exist(key)
	if err != nil {
		logger.Err(err).Error("check redis user id record error")
		return
	}

	// 已存在记录
	if !exist {
		logger.Info("user id record existed")
		return
	}

	max := 0
	pgl.Model(&model.User{}).Select("MAX(user_id)").Row().Scan(&max)
	max = mathe.MaxInt(max, base)

	ok, err := rds.Key().SetNX(key, max)
	if err != nil {
		logger.Err(err).Data(max).Error("set user id record error")
		return
	}

	if ok {
		logger.Data(max).Info("set user id record success")
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// ---------------------------------------------------------------------------------------------------------------------
// 日志接口

var logger = model.Logger(name)
