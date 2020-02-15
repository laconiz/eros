package json

import (
	"fmt"
	jsonIter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

type RawMessage = jsonIter.RawMessage

var json jsonIter.API

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func String(v interface{}) string {
	raw, err := Marshal(v)
	if err == nil {
		return string(raw)
	}
	return fmt.Sprintf("marshal %+v error: %v", v, err)
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
