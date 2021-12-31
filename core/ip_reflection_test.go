package core

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

var ipf *ipReflections

func init() {
	ipf = &ipReflections{
		ips: []*ipRange{
			{start: 0, end: 1},
			{start: 2, end: 5},
			{start: 6, end: 10},
			{start: 12, end: 14},
			{start: 15, end: 20},
			{start: 24, end: 30},
			{start: 35, end: 60},
		},
	}
	err := ipf.LoadFrom("/Users/sjliu/temp/ip")
	if err != nil {
		panic(err)
	}
	fmt.Println("load length:",len(ipf.ips))
}

func Test_ipReflections_getIpRange(t *testing.T) {
	fmt.Println(ipf.getIpRange(996874242))
}

func BenchmarkIp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ipf.getIpRange(996874242)
	}
}

func TestIp(t *testing.T) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, 3707298363)
	fmt.Println(b)
	fmt.Println(fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3]))
}

func TestName(t *testing.T) {
	ip := net.ParseIP("42.62.42.12")
	v := binary.BigEndian.Uint32(ip.To4())
	fmt.Println(v, len(ip))
}

// 404 serviceId
//


var nams struct{
	user struct{}
}
