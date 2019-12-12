package elastic

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
	"reflect"
)

var (
	errNilDoc = errors.New("nil document")
)

type Config struct {
	Address  string // 地址
	User     string // 用户名 无则留空
	Password string // 密码
	Sniff    bool   // 是否开启集群嗅探
	Log      bool   // 记录日志
}

type Elastic struct {
	conf Config
	*elastic.Client
}

// 插入一个文档
func (e *Elastic) Insert(v interface{}) error {

	t := reflect.TypeOf(v)
	if t == nil {
		return errNilDoc
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("value must be struct or struct point, got %v", t)
	}

	_, err := e.Index().
		Index(gorm.ToColumnName(t.Name())).
		Type(Document).
		BodyJson(v).
		Do(context.Background())

	return err
}

const Document = "_doc"

// 新建一个elastic连接
func New(conf Config) (*Elastic, error) {

	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(conf.Address),
		elastic.SetSniff(conf.Sniff),
		// elastic.SetErrorLog(errorLogger),
	}

	if conf.User != "" {
		opts = append(opts, elastic.SetBasicAuth(conf.User, conf.Password))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Elastic{Client: client}, nil
}
