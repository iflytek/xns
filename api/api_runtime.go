package api

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"net"
)

type IpReq struct {
	Ip string `json:"ip" from:"path"`
}

var ipCitiesAPI = &Api{
	Name:               "get ip info",
	Method:             GET,
	Route:              "/ipinfo/:ip",
	ContentType:        fastserver.ContentTypeNone,
	Desc:               "",
	RequestModel:       nil,
	RequestExample:     nil,
	HandleFunc:         nil,
	NotValidateRequest: false,
	Handler: func(ctx *fastserver.Context,req *IpReq)(int ,*Resp) {
		ip := net.ParseIP(req.Ip)
		if ip != nil && ip.To4() != nil{
			rs := core.GetIpInfo(ip.To4())
			return 200,&Resp{Data: rs}
		}
		return 400,&Resp{
			Message: "ip is not valid",
		}

	},
	ResponseExample:    nil,
}

