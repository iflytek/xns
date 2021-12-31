package api

import (
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

type Idc struct {
	Name string `json:"name" desc:"机房名称" format:"name"`
	Desc string `json:"description"`
}

// 机房api

func newIdcsResp(idcs []*models.Idc) *Resp {
	return &Resp{
		Data: idcs,
	}
}

func newIdcResp(idc *models.Idc) *Resp {
	return &Resp{
		Data: idc,
	}
}

// 添加机房
var addIdcApi = &fastserver.Api{
	Name:   "add idc",
	Method: POST,
	Route:  "/idcs",
	Desc:   "添加机房",
	RequestModel: func() interface{} {
		return &Idc{}
	},
	RequestExample:     nil,
	ResponseExample:    newIdcResp(&models.Idc{}),
	NotValidateRequest: false,
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*Idc)
		idc, code, err := idcServiceInstance.Create(req)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, newIdcResp(idc)
	},
}

// 获取机房列表
type IdcReq struct {
}

var getIdcsApi = &fastserver.Api{
	Name:   "get idcs",
	Method: GET,
	Route:  "/idcs",
	Desc:   "获取机房列表",
	RequestModel: func() interface{} {
		return &IdcReq{}
	},
	RequestExample: nil,
	ContentType: fastserver.ContentTypeNone,
	ResponseExample: newIdcsResp([]*models.Idc{
		{
			Name: "dx",
		},
	}),
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		//todo
		idcs, err := idcServiceInstance.GetList()
		if err != nil {
			return 500, newErrorResp(CodeDbError, err.Error())
		}
		return 0, newIdcsResp(idcs)
	},
}

type UpdateIdcReq struct {
	Id string `json:"id" from:"path" desc:"idc name or id"`
	Idc
}

var updateIdcAPI = &fastserver.Api{
	Name:        "update idc",
	Method:      PATCH,
	Route:       "/idcs/:id",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新机房",
	Handler: func(ctx *fastserver.Context, req *UpdateIdcReq) (code int, resp *Resp) {
		md, code, err := idcServiceInstance.Update(req.Id, req.Idc)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, newIdcResp(md)
	},
}

type IdReq struct {
	Id string `json:"id" from:"path" desc:"id"`
}

var getIdcAPI = &fastserver.Api{
	Name:        "get idc",
	Method:      GET,
	Route:       "/idcs/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取机房",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		md, err := idcServiceInstance.Get(req.Id)
		if err != nil {
			if err == dao.NoElemError {
				return 404,notFoundError
			}
			return 500, newErrorResp(CodeDbError, err.Error())
		}
		return 200, newIdcResp(md)
	},
}

var deleteIdcAPI = &fastserver.Api{
	Name:               "delete idc",
	Method:             DELETE,
	Route:              "/idcs/:id",
	ContentType:        fastserver.ContentTypeNone,
	Desc:               "删除机房",
	Handler:       func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp){
		err := idcServiceInstance.Delete(req.Id)
		if err != nil{
			code ,err = convertErrorf(err,"%w",err)
			return mapCodeToHttp(code),newErrorResp(CodeDbError,err.Error())
		}
		return 200,deleteSuccessResp
	},
}
