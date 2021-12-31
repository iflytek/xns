package tools

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)
//split ip to ip and port
//if port is not valid ,port will return -2
func SplitIp(ipp string) (ip string, port int) {
	indexStr := ":"
	offset := 1
	if strings.Contains(ipp,"]"){ // ip v6 地址
		indexStr = "]:"
		offset = 2
	}

	idx := strings.LastIndex(ipp, indexStr)
	if idx < 0 {
		return ipp, -1
	}
	var err error
	port, err = strconv.Atoi(ipp[idx+offset:])
	if err != nil || port < 0 || port > 65535{
		port = -2
	}
	return ipp[:idx+offset-1], port
}

func IsValidIpAddress(ipp string) error {
	//ips, port := SplitIp(ipp)
	//if port == -2 {
	//	return fmt.Errorf("ip '%s' has invalid port", ipp)
	//}

	//ip := net.ParseIP(ips)
	//if ip == nil || !strings.Contains(ips, ".") {
	//	return fmt.Errorf("invalid ip '%s'", ipp)
	//}

	_ ,ok := IsValidIp(ipp)
	if !ok{
		return fmt.Errorf("invalid ip:%s",ipp)
	}
	return nil
}



func IsValidIp(ip string) (int, bool) {
	if strings.Contains(ip, "]:") { // [aa:aa]:90 带端口的ipv6
		kvs := strings.SplitN(ip, "]:", 2)
		ip = strings.TrimLeft(kvs[0], "[") // 去掉端口号和[]
		if !isValidPort(kvs[1]){
			return -1,false
		}

		if isv6(ip) {
			return 1, true
		}
		return -1, false
	}
	// ipv6 不带端口
	if strings.HasPrefix(ip, "[") && strings.HasSuffix(ip, "]") {
		if isv6(ip[1 : len(ip)-1]) {
			return 1, true
		}
		return -1, false
	}
	// ipv4 去掉端口号
	if strings.Contains(ip, ":") {
		kvs := strings.SplitN(ip, ":", 2)
		ip = kvs[0]
		if !isValidPort(kvs[1]){
			return -1,false
		}
	}
	if isv4(ip) {
		return 0, true
	}
	return -1, false
}

func isv6(ip string) bool {
	ipp := net.ParseIP(ip)
	if ipp == nil {
		return false
	}
	if strings.Contains(ip, ":") {
		return true
	}
	return false
}

func isv4(ip string) bool {
	ipp := net.ParseIP(ip)
	if ipp == nil {
		return false
	}
	if strings.Contains(ip, ".") {
		return true
	}
	return false
}


func isValidPort(port string)bool{
	p ,err := strconv.Atoi(port)
	if err != nil{
		return false
	}
	return  p <= 65535 && p > 0
}
