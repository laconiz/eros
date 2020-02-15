package elastic

import "github.com/olivere/elastic"

type Option struct {
	Address  string
	User     string
	Password string
	Sniff    bool
}

func (o *Option) Parse() []elastic.ClientOptionFunc {

	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(o.Address),
		elastic.SetSniff(o.Sniff),
	}

	if o.User != "" && o.Password != "" {
		opts = append(opts, elastic.SetBasicAuth(o.User, o.Password))
	}

	return opts
}
