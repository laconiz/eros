package tokenMgr

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/laconiz/eros/database/redis"
	"github.com/laconiz/eros/oceanus/example/model"
	"github.com/laconiz/eros/oceanus/example/model/db/redisMgr"
	"strconv"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------------------------------------------------

const sep = "-"

var (
	ErrInvalid = errors.New("invalid token")
	ErrExpired = errors.New("expired token")
)

func New(id model.UserID) (string, error) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	src := fmt.Sprintf("%d%s%d", id, sep, now)
	dst := cipher.Encrypt([]byte(src))

	token := hex.EncodeToString(dst)
	if err := hash.Set(id, token); err != nil {
		return "", err
	}

	return token, nil
}

// 校验TOKEN合法性
func Verify(token string) (model.UserID, error) {

	dst, err := hex.DecodeString(token)
	if err != nil {
		return 0, err
	}

	src := cipher.Decrypt(dst)
	arr := strings.Split(string(src), sep)
	if len(arr) != 2 {
		return 0, ErrInvalid
	}

	// 获取TOKEN颁发时间
	ms, err := strconv.ParseInt(arr[1], 10, 64)
	if err != nil {
		return 0, ErrInvalid
	}

	// 校验TOKEN颁发时间
	tm := time.Unix(ms/1000+option.Expired, ms%1000)
	if time.Now().After(tm) {
		return 0, ErrExpired
	}

	// 获取DB中保存的TOKEN
	var stored string
	ok, err := hash.Get(arr[0], &stored)
	if err != nil {
		return 0, err
	}

	// 无法查询到TOKEN
	if !ok {
		return 0, ErrInvalid
	}

	// TOKEN已被更新
	if (stored) != token {
		return 0, ErrExpired
	}

	id, err := strconv.ParseInt(arr[0], 10, 64)
	return model.UserID(id), err
}

// ---------------------------------------------------------------------------------------------------------------------

var hash *redis.Hash

func init() {

	if conn, err := redisMgr.Connect(model.ModuleToken); err != nil {
		panic(fmt.Errorf("connect to redis error: %w", err))
	} else {
		hash = conn.Hash(model.RdsHashToken)
	}

	logger.Info("redis connected")
}

var logger = model.Logger(model.ModuleToken)
