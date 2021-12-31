package api

import (
	"github.com/xfyun/xns/fastserver"
	"github.com/xfyun/xns/models"
)

type Pool struct {
	Description    string                 `json:"description"`
	Name           string                 `json:"name" format:"name"`
	LbMode         int                    `json:"lb_mode" desc:"负载均衡模式,0:机房就近选择，1:就近选择基础上按照机房权重选择" enum:"0,1"`                         // 1，就近分配原则，并且根据权重来分配 2，某一地机房全部没了，更具权重分配到其他机房。
	LbConfig       map[string]interface{} `json:"lb_config" desc:"负载均衡配置"`                       // 负载均衡配置
	FailOverConfig map[string]interface{} `json:"fail_over_config" desc:"负载均衡选择失败时，兜底配置,object"` // fail配置 ，dx:[hu,gz],hu:[dx,gz],gz:[]
}

//新增地址池
var addPoolAPI = &fastserver.Api{
	Name:        "add pool",
	Method:      POST,
	Route:       "/pools",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "新增地址池",
	Handler: func(ctx *fastserver.Context, req *Pool) (code int, resp *Resp) {
		rsp, code, err := poolInstance.Create(req)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(rsp)
	},
}

type updatePoolReq struct {
	IdReq
	Pool
}

// 更新地址池
var updatePoolAPI = &fastserver.Api{
	Name:        "update pool",
	Method:      PATCH,
	Route:       "/pools/:id",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "更新地址池",
	RequestModel: func() interface{} {
		m := make(map[string]interface{})
		return &m
	},
	RequestExample: &updatePoolReq{},
	HandleFunc: func(ctx *fastserver.Context, model interface{}) (code int, resp interface{}) {
		req := model.(*map[string]interface{})
		id, _ := ctx.Params.Get("id")
		rsp, code, err := poolInstance.Update(id, *req)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return 200, &Resp{Data: rsp}
	},
	ResponseExample: &Resp{Data: &models.Pool{}},
}

var getPoolsAPI = &fastserver.Api{
	Name:        "get pools",
	Method:      GET,
	Route:       "/pools",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取所有的地址池",
	Handler: func(ctx *fastserver.Context) (code int, resp *Resp) {
		res, code, err := poolInstance.GetList()
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
}

var getPoolAPI = &fastserver.Api{
	Name:        "get pool",
	Method:      GET,
	Route:       "/pools/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取所有的地址池",
	Handler: func(ctx *fastserver.Context,req *IdReq) (code int, resp *Resp) {
		res, code, err := poolInstance.GetPool(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(res)
	},
}

var deletePoolAPI = &fastserver.Api{
	Name:        "delete pool",
	Method:      DELETE,
	Route:       "/pools/:id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除地址池",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		code, err := poolInstance.Delete(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(nil)
	},
}

type AddPoolGroupReq struct {
	PoolId  string `json:"pool_id" desc:"poolId or name" from:"path"`
	GroupId string `json:"group_id" desc:"groupId or name"`
	Weight  int    `json:"weight" minimum:"1" maximum:"100" desc:"权重"`
}

var poolAddGroupAPI = &fastserver.Api{
	Name:        "add group to pool",
	Method:      POST,
	Route:       "/pools/:pool_id/groups",
	ContentType: fastserver.ContentTypeApplicationJson,
	Desc:        "给地址池添加服务器组",
	Handler: func(ctx *fastserver.Context, req *AddPoolGroupReq) (code int, resp *Resp) {
		ref, code, err := poolInstance.AddPoolGroup(req.PoolId, req.GroupId, req.Weight)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return newSuccessResp(ref)
	},
}

type DeletePoolGroupReq struct {
	PoolId  string `json:"id" desc:"poolId " from:"path"`
	GroupId string `json:"group_id" desc:"groupId or name" from:"path"`
}

var poolDeleteGroupAPI = &Api{
	Name:        "delete pool group",
	Method:      DELETE,
	Route:       "/pools/:id/groups/:group_id",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "删除pool 上的group",
	Handler: func(ctx *fastserver.Context, req *DeletePoolGroupReq) (code int, resp *Resp) {
		code, err := poolInstance.DeletePoolGroup(req.PoolId, req.GroupId)
		if err != nil {
			return mapCodeToHttp(code), newErrorResp(code, err.Error())
		}
		return newSuccessResp(nil)
	},
	ResponseExample: nil,
}

var poolGetGroupAPI = &fastserver.Api{
	Name:        "get pool groups",
	Method:      GET,
	Route:       "/pools/:id/groups",
	ContentType: fastserver.ContentTypeNone,
	Desc:        "获取地址池的服务器组",
	Handler: func(ctx *fastserver.Context, req *IdReq) (code int, resp *Resp) {
		rsp, code, err := poolInstance.GetPoolGroups(req.Id)
		if err != nil {
			return newErrorHttpResp(code, err)
		}
		return newSuccessResp(rsp)
	},
}
