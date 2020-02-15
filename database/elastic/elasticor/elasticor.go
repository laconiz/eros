package elasticor

import "github.com/laconiz/eros/database/elastic"

var client *elastic.Elastic

func init() {

	var option elastic.Option

	var err error
	if client, err = elastic.New(option); err != nil {

	}
}
