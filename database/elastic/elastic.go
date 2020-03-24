package elastic

import (
	"github.com/olivere/elastic"
)

func New(option *Option) (*Elastic, error) {

	client, err := elastic.NewClient(option.Parse()...)
	if err != nil {
		return nil, err
	}

	return &Elastic{Client: client}, nil
}

type Elastic struct {
	*elastic.Client
}

const doc = "_doc"
