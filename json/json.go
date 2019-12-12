package json

import (
	jsonIter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var json jsonIter.API

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func init() {

	extra.RegisterFuzzyDecoders()

	json = jsonIter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()
}
