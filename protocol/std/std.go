package std

import (
	"errors"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/protocol/consts"
	"github.com/xfyun/xns/tools"
	"net"
	"strings"
)

type ErrorResp struct {
	Message string `json:"message"`
}

var NoHostFoundError = &ErrorResp{
	Message: "hosts param is required",
}

var invaldIpResp = &ErrorResp{
	Message: "header X-Mock-Ip should be valid ipv4 ",
}

const (
	contextKey = consts.RequestContext
)

var invalidIpError = errors.New("invalid ip ")
func newErrorMessage(msg string) *ErrorResp {
	return &ErrorResp{Message: msg}
}

var errinvalidIp = &ErrorResp{
	Message: "header X-Mock-Ip or client Ip should be valid ipv4 ",
}

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

var (
	modeIpv6 = []byte("ipv6")
)

//resolve?host=www.baidu.com&appid=100ime
func Handler(ctx *fastserver.Context) {
	fast := ctx.FastCtx
	ctx.SetResponseHeader("Content-Type",fastserver.ContentTypeApplicationJson)
	ctx.SetResponseHeader("Server","ifly-nameserver")
	c := &core.Context{
		Getter: func(key string) string {
			return string(fast.QueryArgs().Peek(key))
		},
	}
	mockIp := ctx.GetRequestHeader("X-Mock-Ip")
	ip := ctx.ClientIp()
	if mockIp != "" {
		ip, _ = parseIp(mockIp)
	}
	if ip == nil{
		ctx.AbortWithStatusJson(400, errinvalidIp)
		return
	}

	hosts := string(fast.QueryArgs().Peek("hosts"))
	if hosts == "" {
		ctx.AbortWithStatusJson(400, NoHostFoundError)
		return
	}
	core.InitRequestContext(c, ip, "")

	res := make([]HostIps, 0, len(hosts))
	var err error
	queryArgs := fast.QueryArgs()
	debug :=queryArgs.GetBool("return_debug_info")
	mode := string(queryArgs.Peek("mode"))
	if mode == "ipv6"{
		c.GetV6 = true
	}else{
		mode = "ipv4"
	}

	for _, host := range strings.Split(hosts, ",") {
		var addr core.Address

		err = core.ResolveIpsByHost(c, host, &addr)

		if err != nil {

			logger.Err().Error("resolve host error:","err:",err," host:",host, " params:",c.Params())
		}
		hip := HostIps{
			Host: host,
			Ips:  convert2IP(addr.Ips, addr.Port),
			Ttl:  addr.Ttl,
		}
		if debug{
			hip.Debug = map[string]interface{}{}
			hip.Debug["pool"] = c.Pool
			hip.Debug["group"] = c.Group
			hip.Debug["service"] = c.Service
			hip.Debug["route"] = c.Route
			hip.Debug["idcAffinity"] = c.IdcAffinity()
			hip.Debug["params"] = c.Params()
			hip.Debug["lbmode"] = c.LbMode

		}
		res = append(res,hip )

	}

	rsp := &Responses{
		Dns:      res,
		ClientIp: ip.To4().String(),
		Mode:mode,
	}
	ctx.SetUserValue(contextKey,c)
	ctx.AbortWithStatusJson(200,rsp)
}

var bt = []byte("true")

func convert2IP(ips []string, p int) (res []Ip) {
	res = make([]Ip, len(ips))
	for i, ip := range ips {
		ip, port := tools.SplitIp(ip)
		if port < 0 {
			port = p
		}
		res[i] = Ip{
			Ip:   ip,
			Port: port,
		}
	}
	return res
}
