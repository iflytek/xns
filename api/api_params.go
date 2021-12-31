package api

import (
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

type AddParamReq struct {
	Name        string `json:"name" format:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

var addParamsAPI = &fastserver.Api{
	Name:            "add parameters",
	Method:          POST,
	Route:           "/parameters",
	ContentType:     fastserver.ContentTypeApplicationJson,
	Desc:            "添加参数对，用于用户自定义参数",
	Handler:         func(ctx *fastserver.Context,req *AddParamReq ) (int, *Resp) {
		res,code ,err := paramInstance.AddParam(req.Name,req.Value,req.Description)
		if err != nil{
			return newErrorHttpResp(code,err)

		}
		return newSuccessResp(res)
	},
	ResponseExample: Resp{Data: &models.CustomParamEnum{}},
}

type DeleteParamReq struct {
	Name        string `json:"name" format:"name" from:"path"`
	Value       string `json:"value" from:"path"`
}

var deleteParamsAPI = &fastserver.Api{
	Name:            "delete parameters",
	Method:          DELETE,
	Route:           "/parameters/:name/:value",
	ContentType:     fastserver.ContentTypeNone,
	Desc:            "删除用户自定义参数",
	Handler:         func(ctx *fastserver.Context,req *DeleteParamReq ) (int, *Resp) {
		code ,err := paramInstance.Delete(req.Name,req.Value)
		if err != nil{
			return newErrorHttpResp(code,err)

		}
		return newSuccessResp(nil)
	},
	ResponseExample: Resp{Data:nil},
}


var getAllParamsAPI = &fastserver.Api{
	Name:            "get All Params",
	Method:          GET,
	Route:           "/parameters",
	ContentType:     fastserver.ContentTypeNone,
	Desc:            "获取所有的参数对",
	Handler:         func(ctx *fastserver.Context ) (int, *Resp) {
		res,code ,err := paramInstance.GetAllParam()
		if err != nil{
			return newErrorHttpResp(code,err)

		}
		return newSuccessResp(res)
	},
}

type NameReq struct {
	Name string `json:"name" from:"path"`
}


var getParamsByNameAPI = &fastserver.Api{
	Name:            "get params value by name",
	Method:          GET,
	Route:           "/parameters/:name",
	ContentType:     fastserver.ContentTypeNone,
	Desc:            "获取所有的参数对",
	Handler:         func(ctx *fastserver.Context,req *NameReq ) (int, *Resp) {
		res,code ,err := paramInstance.GetParamValues(req.Name)
		if err != nil{
			return newErrorHttpResp(code,err)

		}
		return newSuccessResp(res)
	},
}


