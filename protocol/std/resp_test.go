
package std

import (
	"encoding/json"
	"fmt"
	"github.com/xfyun/xns/tools/fastjson"
	"testing"
)

func TestResponses_Json(t *testing.T) {
	resp := &Responses{
		Dns:      []HostIps{
			{
				Host: "baidu.com",
				Ips:  []Ip{
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
				},
				Ttl:  499,
			},
			{
				Host: "baidu.com",
				Ips:  []Ip{
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
					{Ip: "127.0.0.1",Port: 455},
				},
				Ttl:  499,
			},

		},
		ClientIp: "127.0.0.1",
	}

	jw := fastjson.AcquireJsonWriter()
	resp.Json(jw)
	fmt.Println(string(jw.String()))
	fastjson.ReleaseJsonWriter(jw)
}

func BenchmarkJson(b *testing.B) {

	for i := 0; i < b.N; i++ {
		resp := &Responses{
			Dns:      []HostIps{
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},

			},
			ClientIp: "127.0.0.1",
		}
		jw := fastjson.AcquireJsonWriter()
		resp.Json(jw)
		fastjson.ReleaseJsonWriter(jw)
	}
}

func BenchmarkJson2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		resp := &Responses{
			Dns:      []HostIps{
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},

			},
			ClientIp: "127.0.0.1",
		}
		json.Marshal(resp)
	}
}

func BenchmarkEasyjson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp := &Responses{
			Dns:      []HostIps{
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},
				{
					Host: "baidu.com",
					Ips:  []Ip{
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
						{Ip: "127.0.0.1",Port: 455},
					},
					Ttl:  499,
				},

			},
			ClientIp: "127.0.0.1",
		}
		resp.MarshalJSON()

	}
}

