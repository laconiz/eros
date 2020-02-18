//

package logis

import (
	"github.com/laconiz/eros/utils/json"
)

type ContextRaw = json.RawMessage

func NewContext(parent *Context) *Context {
	return &Context{parent: parent, fields: map[string]interface{}{}}
}

type Context struct {
	parent *Context
	fields map[string]interface{}
}

func (context *Context) Field(key string, value interface{}) *Context {
	context.fields[key] = value
	return context
}

func (context *Context) Fields(fields map[string]interface{}) *Context {
	for key, value := range fields {
		context.fields[key] = value
	}
	return context
}

func (context *Context) Context() map[string]interface{} {
	data := map[string]interface{}{}
	for c := context; c != nil; c = c.parent {
		for key, value := range c.fields {
			if _, ok := data[key]; !ok {
				data[key] = value
			}
		}
	}
	return data
}

func (context *Context) Json() []byte {
	data := context.Context()
	raw, err := json.Marshal(data)
	if err != nil {
		return []byte(err.Error())
	}
	return raw
}

func (context *Context) Text() string {
	return string(context.Json())
}
