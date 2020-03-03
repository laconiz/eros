package formatter

import (
	"github.com/laconiz/eros/logis"
)

type Formatter interface {
	Format(*logis.Log) ([]byte, error)
}

const DefaultTimeLayout = "2006-01-02 15:04:05.000"
