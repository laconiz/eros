package context

import (
	"fmt"
	"github.com/laconiz/eros/utils/json"
)

type Fields map[string]interface{}

func New(parent *Context) *Context {
	context := &Context{fields: Fields{}}
	if parent != nil {
		context.Fields(parent.fields)
	}
	return context
}

type Context struct {
	raw    []byte
	fields Fields
}

func (context *Context) Fields(fields Fields) *Context {
	for key, value := range fields {
		context.fields[key] = value
	}
	return context
}

func (context *Context) Json() []byte {
	if len(context.fields) == 0 {
		return nil
	}
	if context.raw != nil {
		return context.raw
	}
	raw, err := json.Marshal(context.fields)
	if err == nil {
		context.raw = raw
	} else {
		context.raw = []byte(fmt.Sprint(context.fields))
	}
	return context.raw
}

func (context *Context) Text() string {
	return string(context.Json())
}
