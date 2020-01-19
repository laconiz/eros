package command

import (
	"os"
	"strings"
)

func ParseStringArg(name string, value string) string {
	if str, ok := parseArg(name); ok {
		return str
	}
	return value
}

func parseArg(name string) (string, bool) {

	unix := "-" + name
	gnu := "--" + name

	for index, arg := range os.Args {

		if arg == unix && index+1 < len(os.Args) {
			return os.Args[index+1], true
		}

		if strings.HasPrefix(arg, gnu) {
			sp := strings.SplitN(arg, "=", 1)
			if len(sp) == 2 {
				return sp[1], true
			}
		}
	}

	return "", false
}
