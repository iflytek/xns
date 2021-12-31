package msc

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/tools"
	"net"
)

type JsonMarshal interface {
	Json(buf *bytes.Buffer) error
}

//--- PASS: Test_newRequestContext (1.09s)
type DnsResponse struct {
	Dns     []*HostIp `json:"dns"`
	Default []*HostIp `json:"default"`
	Mode    string    `json:"mode,omitempty"`
}

type HostIp struct {
	Host string `json:"host"`
	Sip  []Sip  `json:"sip"`
}

type Sip struct {
	Debug  map[string]interface{} `json:"debug,omitempty"`
	Svc    string                 `json:"svc,omitempty"`
	Ips    []Ip                   `json:"ips"`
	Ttl    int                    `json:"ttl"`
	Compel *Compel                `json:"compel,omitempty"`
}

type Compel struct {
	Value bool
}

var (
	bTrue       = []byte("1")
	bFalse      = []byte("0")
	compelTrue  = &Compel{Value: true}
	compelFalse = &Compel{Value: false}
)

func (c *Compel) MarshalJSON() ([]byte, error) {
	if c.Value {
		return bTrue, nil
	}
	return bFalse, nil
}

type Ip struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
	Idc  string `json:"idc,omitempty"`
}

type ErrorResp struct {
	Message string `json:"message"`
}

func newErrorMessage(msg string) *ErrorResp {
	return &ErrorResp{Message: msg}
}

var errinvalidIp = &ErrorResp{
	Message: "header X-Mock-Ip or client Ip should be valid ipv4 ",
}

func FetchOne(ctx *fastserver.Context) {
	c := &core.Context{}
	parameters := ctx.GetRequestHeader("X-Par")
	mockIp := ctx.GetRequestHeader("X-Mock-Ip")
	ip := ctx.ClientIp()
	ctx.SetResponseHeader("Content-Type", fastserver.ContentTypeApplicationJson)
	ctx.SetResponseHeader("Server", "ifly-nameserver")
	if mockIp != "" {
		ip, _ = parseIp(mockIp)
	}
	if ip == nil {
		ctx.AbortWithStatusJson(400, errinvalidIp)
		return
	}
	mode := string(ctx.FastCtx.QueryArgs().Peek("mode"))
	if mode == "ipv6" {
		c.GetV6 = true
	}

	core.InitRequestContext(c, ip, parameters)
	var res core.Address
	err := core.ResolveIpsByHost(c, "", &res)
	if err != nil {
		//ctx.AbortWithStatusJson(500, newErrorMessage(err.Error()))
		logger.Err().Error("resolve dns error", err, ",ctx", c)
		//return
	}
	rsp := newDnsResponse(&res, c, ctx, mode)
	ctx.AbortWithStatusJson(200, rsp)
	//ctx.AbortWithData(200,rsp.Bytes())
	//rsp.Release()
}

func newDnsResponse(address *core.Address, c *core.Context, ctx *fastserver.Context, mode string) *DnsResponse {
	var sip Sip
	newHostIp(&sip, address)
	initDebugInfo(ctx.FastCtx.QueryArgs().GetBool("return_debug_info"), &sip, c)
	dns := &HostIp{
		Host: address.Host,
		Sip: []Sip{
			sip,
		},
	}
	def := newDefaultHostIp(address)
	if len(def.Ips) == 0 {
		switch len(sip.Ips) {
		case 0:

		case 1:
			def = sip
		default:
			def.Ips = sip.Ips[1:]
		}

	}
	return &DnsResponse{
		Mode: mode,
		Dns:  []*HostIp{dns},
		Default: []*HostIp{{
			Host: address.Host,
			Sip: []Sip{
				def,
			},
		},
		},
	}
}

func parseIps(ips []string, idc string, p int, ver string) []Ip {
	switch ver {
	case "1.0":
		idc = ""
	}
	ipss := make([]Ip, len(ips))
	for i, ip := range ips {
		port := 0
		ip, port = tools.SplitIp(ip)
		if port < 0 {
			port = p
		}
		ipss[i] = Ip{
			Ip:   ip,
			Port: port,
			Idc:  idc,
		}
	}
	return ipss
}

func newHostIp(ip *Sip, address *core.Address) {
	var cpl *Compel = nil
	switch address.Ver {
	case "1.3":
		if !address.Mul {
			cpl = compelFalse
		}
	}
	ips := parseIps(address.Ips, address.IdcName, address.Port, address.Ver)
	*ip = Sip{
		Svc:    address.Svc,
		Ips:    ips,
		Ttl:    address.Ttl,
		Compel: cpl,
	}
}

func newDefaultHostIp(address *core.Address) Sip {

	ips := parseIps(address.Default, address.IdcName, address.Port, address.Ver)
	var cpl *Compel
	switch address.Ver {
	case "1.3":
		if !address.Mul {
			cpl = compelFalse
		}
	}

	return Sip{
		Svc:    address.Svc,
		Ips:    ips,
		Ttl:    address.Ttl,
		Compel: cpl,
	}
}

var invalidIpError = errors.New("invalid ip ")

func parseIp(s string) (net.IP, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, invalidIpError
	}
	ip = ip.To4()
	if ip == nil {
		return nil, invalidIpError
	}
	return ip, nil
}

func parseHosts(body []byte) (res []string) {
	sc := bufio.NewScanner(bytes.NewReader(body))
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	return
}
