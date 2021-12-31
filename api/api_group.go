package api

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

type HealthyCheckConfig struct {
	Host         string `json:"host" desc:"mode=http时必定传，"`
	Method       string `json:"method" desc:"mode=http时必定传，" enum:"GET,POST,PATCH,PUT,HEAD,OPTION"`
	Path         string `json:"path" desc:"mode=http时必定传，"`
	Timeout      int    `json:"timeout" desc:"健康检查超时时间，单位s" minimum:"1" maximum:"60"`
	Body         string `json:"body" desc:"mode=http 时需要"`
	SuccessCodes []int  `json:"success_codes" desc:"mode=http 时需要，健康检查成功响应status code" maximum:"599" minimum:"100"`
}

type Group struct {
	Desc               string                 `json:"description"`
	Name               string                 `json:"name" desc:"group 名称" required:"true" format:"name"`
	HealthyCheckMode   string                 `json:"healthy_check_mode" desc:"健康检查模式" enum:"tcp,http"` // 是否开启healthy check  ,关闭
	HealthyCheckConfig *HealthyCheckConfig    `json:"healthy_check_config" desc:"健康检查配置,json 格式"`       // 健康检查模式
	HealthyNum         int                    `json:"healthy_num" desc:"对于检查不健康的节点，成功多少次认为健康，0表示不开启健康检查" minimum:"0"`
	UnHealthyNum       int                    `json:"un_healthy_num" desc:"对于健康的节点，失败多少次认为节点不健康" minimum:"0"`
	HealthyInterval    int                    `json:"healthy_interval" desc:"对于健康的节点，健康检查时间间隔，单位s " minimum:"0"`
	UnHealthyInterval  int                    `json:"unhealthy_interval" desc:"对于不健康的节点，健康检查时间间隔，单位s " minimum:"0"`
	LbMode             string                 `json:"lb_mode" desc:"负载均衡方式" enum:"loop"`
	LbConfig           map[string]interface{} `json:"lb_config" desc:"负载均衡配置,json 格式"`
	ServerTags         map[string]string      `json:"server_tags" desc:"给服务器打的标签,json "`                                 // 172.21.164.32:
	//Weight             int                    `json:"weight" desc:"group 的权重" minimum:"0" maximum:"100" required:"true"` // 机房的权重
	IpAllocNum         int                    `json:"ip_alloc_num" desc:"一次性下发的ip数量" minimum:"0"`                        // 一次下发多少数量的ip，0 为全部下发，最大不超过group 中包含的ip
	DefaultServers     string                 `json:"default_servers" format:"ipv4s" desc:"默认ip集合"`
	Port               int                    `json:"port" desc:"下发的端口号，同时用于健康检查" maximum:"65535" minimum:"1" required:"true"`
}

type AddGroupReq struct {
	Group
	IdcId string `json:"idc_id" desc:"group 所在机房" required:"true" ` // 机房名称
}

var addGroupAPI = &fastserver.Api{
	Name:        "add group",
	Method:      POST,
	Route:       "/groups",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "新增group",
	Handler: func(ctx *fastserver.Context, req *AddGroupReq) (statusCode int, resp *Resp) {
		gp, code, err := groupInstance.Create(req)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: gp}
	},
	ResponseExample: &Resp{Data: &models.Group{}},
}

type UpdateGroupReq struct {
	Id string `json:"id" desc:"group id or name" from:"path"`
	Group
}

var updateGroupAPI = &fastserver.Api{
	Name:        "update group",
	Method:      PATCH,
	Route:       "/groups/:id",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新group",
	RequestModel: func() interface{} {
		m := make(map[string]interface{})
		return &m
	},
	RequestExample: &UpdateGroupReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*map[string]interface{})
		id, _ := ctx.Params.Get("id")
		grp, code, err := groupInstance.Update(id, *req)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: grp}
	},
	ResponseExample: &Resp{Data: &models.Group{}},
}

var getListGroupAPI = &fastserver.Api{
	Name:        "get group list",
	Method:      GET,
	Route:       "/groups",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取服务器组列表",
	Handler: func(ctx *fastserver.Context) (code int, resp *Resp) {
		data, code, err := groupInstance.GetList()
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: data}
	},
}

var getGroupAPI = &fastserver.Api{
	Name:        "get group",
	Method:      GET,
	Route:       "/groups/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "根据id 或者name 获取 group",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		data, code, err := groupInstance.Get(req.Id)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: data}
	},
	ResponseExample: &Resp{Data: &models.Group{}},
}

var deleteGroupAPI = &fastserver.Api{
	Name:        "delete group",
	Method:      DELETE,
	Route:       "/groups/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "根据id 或者name 获取 删除group",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		var err error
		code, err = groupInstance.Delete(req.Id)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: nil}
	},
}

// 向服务器组添加服务器
type AddGroupServerReq struct {
	Id       string `json:"id" desc:"group_id or name" from:"path"`
	ServerId string `json:"server_ip" desc:"server_ip" required:"true" format:"ipv4"`
	Weight   int    `json:"weight" minimum:"1" maximum:"100" required:"true"`
}

var groupAddServerAPI = &fastserver.Api{
	Name:        "add server to group",
	Method:      POST,
	Route:       "/groups/:id/servers",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "向服务器组添加服务器",
	Handler: func(ctx *fastserver.Context, req *AddGroupServerReq) (code int, resp *Resp) {
		ref, code, err := groupInstance.AddServer(req.Id, req.ServerId, req.Weight)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: ref}
	},
	ResponseExample: &Resp{Data: &models.ServerGroupRef{}},
}

type DeleteGroupServerReq struct {
	Id       string `json:"id" from:"path"`
	ServerId string `json:"server_id" from:"path"`
}

var groupDeleteServerAPI = &fastserver.Api{
	Name:        "delete group server",
	Method:      DELETE,
	Route:       "/groups/:id/servers/:server_id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除服务器组的服务器",
	Handler: func(ctx *fastserver.Context, req *DeleteGroupServerReq) (code int, resp *Resp) {
		code, err := groupInstance.DeleteServer(req.Id, req.ServerId)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{}
	},
}

var groupGetServerAPI = &fastserver.Api{
	Name:        "get servers of group",
	Method:      GET,
	Route:       "/groups/:id/servers",
	ContentType: fastserver.ContentTypeNone,
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		res, code, err := groupInstance.GetGroupServers(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
	ResponseExample: &Resp{Data: []*GroupServer{{}}},
}

var groupStatusAPI = &fastserver.Api{
	Name:        "get group status",
	Method:      GET,
	Route:       "/group_status",
	ContentType: fastserver.ContentTypeNone,
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		res := core.GetGroupStatus()
		return newSuccessResp(res)
	},
}


//var getLBModes = &fastserver.Api{
//	Name:        "get group status",
//	Method:      GET,
//	Route:       "/pool_lb_modes",
//	ContentType: fastserver.ContentTypeNone,
//	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
//		return newSuccessResp()
//	},
//}
