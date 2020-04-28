// 命令行参数解析
package command

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// 从参数中获取合法的IP地址及端口
func ParseAddress(name string, value string) (string, error) {

	addr, ok := parse(name)
	if !ok {
		addr = name
	}

	items := strings.Split(addr, ":")
	if len(items) != 2 {
		return "", fmt.Errorf("invalid address[%s]", addr)
	}

	if ip := net.ParseIP(items[0]); ip == nil {
		return "", fmt.Errorf("invalid ip[%s]", items[0])
	}

	port, err := strconv.Atoi(items[1])
	if err != nil || port < 0 || port > 65535 {
		return "", fmt.Errorf("invalid port[%s]", items[1])
	}

	return addr, nil
}

func ParseStringArg(name string, value string) string {
	if str, ok := parse(name); ok {
		return str
	}
	return value
}

func parse(name string) (string, bool) {

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
