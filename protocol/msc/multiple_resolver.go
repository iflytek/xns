package msc

import (
	"encoding/json"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/logger"
)

type MultipleRequest struct {
	Hosts []HostSrv `json:"hosts"`
}

type HostSrv struct {
	Host string   `json:"host"`
	Svcs []string `json:"svcs"`
}

func parseRequest(req *MultipleRequest, b []byte) error {
	err := json.Unmarshal(b, req)
	if err != nil {
		return err
	}
	return nil
}

var (
	requestIsNotJsonErr = &ErrorResp{
		Message: "request is not json",
	}
)

func MultipleResolver(ctx *fastserver.Context) {
	c := &core.Context{}
	parameters := ctx.GetRequestHeader("X-Par")
	mockIp := ctx.GetRequestHeader("X-Mock-Ip")
	ip := ctx.ClientIp()
	ctx.SetResponseHeader("Content-Type",fastserver.ContentTypeApplicationJson)
	ctx.SetResponseHeader("Server","ifly-nameserver")
	if mockIp != "" {
		ip, _ = parseIp(mockIp)
	}
	if ip == nil {
		ctx.AbortWithStatusJson(400, errinvalidIp)
		return
	}
	core.InitRequestContext(c, ip, parameters)

	req := &MultipleRequest{}
	if err := parseRequest(req, ctx.FastCtx.PostBody()); err != nil {
		ctx.AbortWithStatusJson(200, requestIsNotJsonErr)
		return
	}
	mode := string(ctx.FastCtx.QueryArgs().Peek("mode"))
	if mode == "ipv6"{
		c.GetV6 = true
	}
	resp := &MultipleResponse{
		Mode: mode,
	}
	for _, host := range req.Hosts {
		hip := HostIp{
			Host: host.Host,
		}

		def := HostIp{
			Host: host.Host,
		}

		for _, svc := range host.Svcs {
			c.Set("svc", svc)
			addr := &core.Address{
				Mul: true,
			}
			err := core.ResolveIpsByHost(c, host.Host, addr)
			if err != nil {
				logger.Err().Error("resolve multiple error:", err, ",host:", host.Host, ",svc", svc)
				//continue
			}
			sip := Sip{}
			newHostIp(&sip, addr)
			initDebugInfo(ctx.FastCtx.QueryArgs().GetBool("return_debug_info"),&sip,c)
			hip.Sip = append(hip.Sip, sip)
			defip := newDefaultHostIp(addr)
			if len(defip.Ips) == 0 {
				defip = sip
			}
			def.Sip = append(def.Sip, defip)
		}

		resp.Dns = append(resp.Dns, hip)
		resp.Default = append(resp.Default, hip)
	}
	ctx.AbortWithStatusJson(200, resp)

}

type MultipleResponse struct {
	Dns     []HostIp `json:"dns"`
	Default []HostIp `json:"default"`
	Mode string `json:"mode"`
}

func initDebugInfo(ok bool, hip *Sip, c *core.Context) {
	if !ok {
		return
	}
	hip.Debug = map[string]interface{}{}
	hip.Debug["pool"] = c.Pool
	hip.Debug["group"] = c.Group
	hip.Debug["service"] = c.Service
	hip.Debug["route"] = c.Route
	hip.Debug["idcAffinity"] = c.IdcAffinity()
	hip.Debug["params"] = c.Params()
	hip.Debug["lbmode"] = c.LbMode
	hip.Debug["clientIp"] = c.ClientIp.String()
}
