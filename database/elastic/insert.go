package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
)

func (e *Elastic) Insert(v interface{}) error {

	index, err := structIndex(v)
	if err != nil {
		return fmt.Errorf("insert document error: %w", err)
	}

	if _, err = e.Index().Index(index).Type(doc).BodyJson(v).Do(context.Background()); err != nil {
		return fmt.Errorf("insert document error: %w", err)
	}

	return nil
}

func (e *Elastic) InsertRaw(index string, raw []byte) error {

	if len(raw) == 0 {
		return errors.New("insert document raw error: raw is empty")
	}

	if _, err := e.Index().Index(index).Type(doc).BodyString(string(raw)).Do(context.Background()); err != nil {
		return fmt.Errorf("insert document raw error: %w", err)
	}
	return nil
}

func (e *Elastic) Inserts(v interface{}) error {

	name, err := sliceIndex(v)
	if err != nil {
		return fmt.Errorf("insert slice error: %w", err)
	}

	value := reflect.ValueOf(v)
	bulk := e.Bulk()
	for i := 0; i < value.Len(); i++ {
		index := value.Index(i)
		if !index.CanInterface() {
			return fmt.Errorf("insert slice error: invalid %dth value", i)
		}
		if index.Type().Kind() == reflect.Ptr && index.IsNil() {
			return fmt.Errorf("insert slice error: value %dth is nil", i)
		}
		bulk.Add(elastic.NewBulkIndexRequest().Index(name).Type(doc).Doc(index.Interface()))
	}

	if _, err = bulk.Do(context.Background()); err != nil {
		return fmt.Errorf("insert slice error: %w", err)
	}
	return nil
}

func (e *Elastic) InsertRaws(index string, raws [][]byte) error {

	bulk := e.Bulk()
	for idx, raw := range raws {
		if len(raw) == 0 {
			return fmt.Errorf("insert slice raw error: raw %dth is empty", idx)
		}
		bulk.Add(elastic.NewBulkIndexRequest().Index(index).Type(doc).Doc((json.RawMessage)(raw)))
	}

	if _, err := bulk.Do(context.Background()); err != nil {
		return fmt.Errorf("insert slice raw error: %w", err)
	}
	return nil
}
