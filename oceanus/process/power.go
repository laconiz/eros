package process

import (
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"
)

// 获取一个IPV4地址的权值
func addrPower(addr string) (uint64, error) {
	// 分离IP和端口
	ap := strings.Split(addr, ":")
	if len(ap) != 2 {
		return 0, errors.New("invalid addr format")
	}
	// 解析IP地址
	ip := net.ParseIP(ap[0])
	if ip == nil {
		return 0, errors.New("invalid ip address")
	}
	// 反序列化端口号
	port, err := strconv.ParseUint(ap[1], 10, 64)
	if err != nil || port > 65535 {
		return 0, errors.New("invalid port address")
	}
	// 计算权值
	power := big.NewInt(0).SetBytes(ip.To4()).Uint64()
	return port<<32 | power, nil
}
