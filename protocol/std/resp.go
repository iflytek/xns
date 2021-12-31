package std

import "github.com/xfyun/xns/tools/fastjson"

type HostIps struct {
	Debug map[string]interface{} `json:"debug,omitempty"`
	Host string `json:"host"`
	Ips  []Ip   `json:"ips"`
	Ttl  int    `json:"ttl"`

}

func (h HostIps) Json(jw *fastjson.JsonWriter) {
	jw.WriteObjectStart()
	jw.WriteString("host", h.Host)
	jw.WriteSep()
	jw.WriteArrayLeft("ips")
	ips := h.Ips
	switch len(ips) {
	case 0:
	case 1:
		ips[0].Json(jw)
	default:
		ips[0].Json(jw)
		for i := 1; i < len(ips); i++ {
			jw.WriteSep()
			ips[i].Json(jw)
		}
	}
	jw.WriteArrayRight()
	jw.WriteSep()
	jw.WriteInt("ttl", h.Ttl)
	jw.WriteObjectRight()
}

type Ip struct {
	Ip string `json:"ip"`
	//Tags map[string]string `json:"tags,omitempty"`
	Port int `json:"port"`
}

func (i Ip) Json(jw *fastjson.JsonWriter) {
	jw.WriteObjectStart()
	jw.WriteString("ip", i.Ip)
	jw.WriteSep()
	jw.WriteInt("port", i.Port)
	jw.WriteObjectRight()
}
//easyjson:json
type Responses struct {
	Dns      []HostIps `json:"dns"`
	ClientIp string    `json:"client_ip"`
	Mode string `json:"mode"`
}

func (r Responses) Json(jw *fastjson.JsonWriter) {
	jw.WriteObjectStart()
	jw.WriteArrayLeft("dns")
	dns := r.Dns
	switch len(dns) {
	case 0:

	case 1:
		dns[0].Json(jw)
	default:
		dns[0].Json(jw)
		for i := 1; i < len(dns); i++ {
			jw.WriteSep()
			dns[i].Json(jw)
		}
	}
	jw.WriteArrayRight()
	jw.WriteSep()
	jw.WriteString("client_ip",r.ClientIp)
	jw.WriteObjectRight()
}

type U struct {
	A int
}
