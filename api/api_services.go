package api

import (
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

var addServiceAPI = &fastserver.Api{
	Name:        "create service",
	Method:      POST,
	Route:       "/services",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "创建服务",
	Handler: func(ctx *fastserver.Context, req *Service) (code int, resp *Resp) {
		srv, code, err := serviceInstance.Create(req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(srv)
	},
}

type UpdateServiceReq struct {
	Id string `json:"id" desc:"service id or service name" format:"name" from:"path"`
	Service
}

var updateServiceAPI = &fastserver.Api{
	Name:        "update service",
	Method:      PATCH,
	Route:       "/services/:id",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新服务",
	RequestModel: func() interface{} {
		m := make(map[string]interface{})
		return &m
	},
	RequestExample: &UpdateServiceReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*map[string]interface{})
		id, _ := ctx.Params.Get("id")
		srv, code, err := serviceInstance.Update(id, *req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(srv)
	},
	ResponseExample: &Resp{Data: &models.Service{}},
}

var getServiceAPI = &fastserver.Api{
	Name:        "get service ",
	Method:      GET,
	Route:       "/services/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "根据id或者name获取服务",
	Handler: func(ctx *fastserver.Context, id *IdReq) (code int, resp *Resp) {
		srv, code, err := serviceInstance.Get(id.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(srv)
	},
	ResponseExample: &Resp{Data: &models.Service{}},
}

var getServicesAPI = &fastserver.Api{
	Name:               "get services list",
	Method:             GET,
	Route:              "/services",
	ContentType:        fastserver.ContentTypeNone,
	Desc:               "获取服务列表",
	HandleFunc:         nil,
	NotValidateRequest: false,
	Handler: func(ctx *fastserver.Context) (code int, resp *Resp) {
		srvs, code, err := serviceInstance.GetList()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(srvs)
	},
}

var deleteServiceAPI = &fastserver.Api{
	Name:        "delete service",
	Method:      DELETE,
	Route:       "/services/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除服务",
	Handler: func(ctx *fastserver.Context, id *IdReq) (code int, resp *Resp) {
		code, err := serviceInstance.Delete(id.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(nil)
	},
}

type AddRouteReq struct {
	ServiceId string `json:"service_id" desc:"service id" from:"path"`
	Route
}

var serviceAddRouteAPI = &fastserver.Api{
	Name:        "service add route",
	Method:      POST,
	Route:       "/services/:service_id/routes",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "给服务添加路由",
	Handler: func(ctx *fastserver.Context, req *AddRouteReq) (code int, resp *Resp) {
		r, code, err := serviceInstance.AddRoutes(req.ServiceId, &req.Route)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(r)
	},
}

var serviceGetRoutesAPI = &fastserver.Api{
	Name:        "service get routes",
	Method:      GET,
	Route:       "/services/:id/routes",
	ContentType: fastserver.ContentTypeNone,
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		r, code, err := serviceInstance.GetRoutes(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(r)
	},
	ResponseExample: nil,
}

var serviceDeleteRouteAPI = &fastserver.Api{
	Name:        "delete route",
	Method:      DELETE,
	Route:       "/routes/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除路由",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		code, err := serviceInstance.DeleteRoutes(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(nil)
	},
	ResponseExample: nil,
}

type queryRoutesReq struct {
	Host string `json:"host" desc:"host 模糊匹配" from:"query"`
	Rule string `json:"rule" desc:"rule 模糊匹配" from:"query"`
}

var getRoutesAPI = &fastserver.Api{
	Name:        "get all routes",
	Method:      GET,
	Route:       "/routes",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取全部的路由",
	Handler: func(ctx *fastserver.Context,req *queryRoutesReq) (code int, resp *Resp) {
		rs ,code, err := serviceInstance.GetAllRoutes(req.Host,req.Rule)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(rs)
	},
	ResponseExample: &Resp{Data: []*models.Route{{}}},
}

var getRouteAPI = &fastserver.Api{
	Name:        "get  route",
	Method:      GET,
	Route:       "/routes/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取路由",
	Handler: func(ctx *fastserver.Context,req *IdReq) (code int, resp *Resp) {
		rs ,code, err := serviceInstance.GetRoute(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(rs)
	},
	ResponseExample: Resp{Data: &models.Route{}},
}

type UpdateRouteReq struct {
	RouteId string `json:"id" from:"path" format:"uuid"`
	Route
}

var updateRoutesAPI = &fastserver.Api{
	Name:               "update routes",
	Method:             PATCH,
	Route:              "/routes/:id",
	ContentType:        fastserver.ContentTypeApplicationJson,
	Desc:               "更新路由",
	RequestModel: func() interface{} {
		m:=make(map[string]interface{})
		return &m
	},
	RequestExample:     &UpdateRouteReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		id,_:=ctx.Params.Get("id")
		req := model.(*map[string]interface{})
		rsp,code,err := serviceInstance.UpdateRoute(id,*req)
		if err != nil{
			return newErrorHttpResp(code,err)
		}
		return newSuccessResp(rsp)
	},
	NotValidateRequest: false,
	ResponseExample:    &Resp{Data: &models.Route{}},
}


