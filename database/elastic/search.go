package elastic

import (
	"github.com/olivere/elastic"
	"reflect"
)

type Query interface {
	Query() elastic.Query
}

func makeQuery(value reflect.Value) elastic.Query {

	var queries []elastic.Query

	for i := 0; i < value.NumField(); i++ {
		if query, ok := value.Field(i).Interface().(Query); ok {
			if query := query.Query(); query != nil {
				queries = append(queries, query)
			}
		}
	}

	switch len(queries) {
	case 0:
		return elastic.NewMatchAllQuery()
	case 1:
		return queries[0]
	default:
		return elastic.NewBoolQuery().Must(queries...)
	}
}

type Sorter struct {
	Sorters []elastic.Sorter
	Values  []interface{}
	Element reflect.Type
}
