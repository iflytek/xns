package resource

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_generateCites(t *testing.T) {
	if err := generateCites();err != nil{
		panic(err)
	}
}

func Test_generateProvince(t *testing.T) {
	if err := generateProvince();err != nil{
		panic(err)
	}
}


type key struct {
	Ip string
	port int
	namespace string
}

func asmk(ip,nm string,port int)string{
	return  "ip_"+ip+"nm_"+nm+"port_"+strconv.Itoa(port)
}

func TestMap(t *testing.T) {
	m := map[key]string{}
	m[key{Ip: "dsd",port: 34}] = "234"

	fmt.Println()
}

func BenchmarkTest(b *testing.B) {
	m := map[key]string{}

	m[key{Ip: "dsd",port: 34,namespace: "gbk"}] = "234"
	var rs string
	for i := 0; i < b.N; i++ {
		rs = m[key{Ip: "dsd",port: 34,namespace: "gbk"}]
	}



	fmt.Println(rs)
}

func BenchmarkTest2(b *testing.B) {
	m := map[string]string{}
	m[asmk("dsd","gbk",34)] = "234"
	var rs string
	for i := 0; i < b.N; i++ {
		rs = m[asmk("dsd","gbk",34)]
	}
	fmt.Println(rs)

}
