package dao

import (
	_ "github.com/bmizerany/pq"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools"
	"github.com/xfyun/xns/tools/uid"
)

type Idc interface {
	GetById(id string) (*models.Idc, error)                     // 通过id 获取
	GetByName(name string) (*models.Idc, error)                 // 通过名称获取
	GetByIdOrName(idOrName string) (idc *models.Idc, err error) // 通过名称或者id获取
	GetList() (idcs []*models.Idc, err error)                   //获取列表
	Create(idc *models.Idc) error                               // 新增
	Update(id string, idc *models.Idc) error                    // 更新
	Delete(id string) error                                     // 删除
}

//type Server interface {
//	GetById(id string) (*models.Server, error)                     // 通过id 获取
//	GetByName(name string) (*models.Server, error)                 // 通过名称获取
//	GetByIdOrName(idOrName string) (idc *models.Server, err error) // 通过名称或者id获取
//	GetList() (srvs []*models.Server, err error)                   //获取列表
//	Create(srv *models.Server) error                               // 新增
//	Update(id string, srv *models.Server) error                    // 更新
//	Delete(id string) error                                        // 删除
//	QueryIdcReference(idc_id string) (bool, error)
//}

type Group interface {
	GetById(id string) (*models.Group, error)                     // 通过id 获取
	GetByName(name string) (*models.Group, error)                 // 通过名称获取
	GetByIdOrName(idOrName string) (idc *models.Group, err error) // 通过名称或者id获取
	GetList() (srvs []*models.Group, err error)                   //获取列表
	Create(srv *models.Group) error                               // 新增
	Update(id string, srv *models.Group) error                    // 更新
	Delete(id string) error                                       // 删除
	Patch(id string, group map[string]interface{}) error
}

type GroupServerRef interface {
	GetById(id string) (*models.ServerGroupRef, error)   // 通过id 获取
	GetList() (srvs []*models.ServerGroupRef, err error) //获取列表
	Create(srv *models.ServerGroupRef) error             // 新增
	Update(id string, srv *models.ServerGroupRef) error  // 更新
	Delete(id string) error
	GetByRef(serverId, groupId string) (res *models.ServerGroupRef, err error)
	GetGroupServers(groupId string) (res []*models.ServerGroupRef, err error)
}

type Pool interface {
	GetById(id string) (*models.Pool, error) // 通过id 获取
	GetByName(name string) (pool *models.Pool, err error)
	GetByIdOrName(idOrName string) (idc *models.Pool, err error) // 通过名称或者id获取
	GetList() (pools []*models.Pool, err error)                  //获取列表
	Create(pool *models.Pool) error                              // 新增
	Update(id string, srv *models.Pool) error                    // 更新
	Delete(id string) error
	Patch(id string, pool map[string]interface{}) error
}

type GroupPoolRef interface {
	GetById(id string) (res *models.GroupPoolRef, err error) // 通过id 获取
	GetList() (refs []*models.GroupPoolRef, err error)       //获取列表
	Create(ref *models.GroupPoolRef) error                   // 新增
	Update(id string, srv *models.GroupPoolRef) error        // 更新
	Delete(id string) error
	Patch(id string, pool map[string]interface{}) error
	GetPoolGroupRef(poolId, groupId string) (ref *models.GroupPoolRef, err error)
	GetPoolGroups(poolId string) (refs []*models.GroupPoolRef, err error)
}

type Service interface {
	GetById(id string) (*models.Service, error) // 通过id 获取
	GetByName(name string) (pool *models.Service, err error)
	GetByIdOrName(idOrName string) (pool *models.Service, err error) // 通过名称或者id获取
	GetList() (pools []*models.Service, err error)                   //获取列表
	Create(pool *models.Service) error                               // 新增
	Update(id string, srv *models.Service) error                     // 更新
	Delete(id string) error
	Patch(id string, pool map[string]interface{}) error
}

type Route interface {
	GetById(id string) (res *models.Route, err error) // 通过id 获取
	GetList() (refs []*models.Route, err error)       //获取列表
	Create(ref *models.Route) error                   // 新增
	Update(id string, srv *models.Route) error        // 更新
	Delete(id string) error
	Patch(id string, pool map[string]interface{}) error
	GetServiceRoutes(serviceId string) (res []*models.Route, err error)
	QueryRoutes(host,rule string)(res []*models.Route,err error) // 模糊查询
	QueryRoutesByRuleCond(conds string)(res []*models.Route,err error)
}

type City interface {
	Create(c *models.City)(err error)
	GetByCode(code int) (city *models.City, err error)
	GetById(id string) (city *models.City, err error)
	GetProvinceCities(provCode int) (city []*models.City, err error)
	GetList() (res []*models.City, err error)
	Delete(code int) error
	Update(code int,c *models.City)(err error)
	IfReferenceIdc(idcId string) (bool, error)
	Init() error
}

type Province interface {
	GetByCode(code int) (res *models.Province, err error)
	GetList() (res []*models.Province, err error)
	GetById(id string) (res *models.Province, err error)
	GeProvinceByRegionCode(regionCode int) (res []*models.Province, err error)
	Delete(code int) error
	Create(p *models.Province)(err error)
	IfReferenceIdc(idcId string) (bool, error)
	Update(code int ,p *models.Province)(err error)
	Init() error
}

type Region interface {
	GetByCode(code int) (res *models.Region, err error)
	GetById(id string) (res *models.Region, err error)
	GetList() (res []*models.Region, err error)
	Create(region *models.Region) error
	Update(id string, region *models.Region) error
	Delete(code int) error
	IfReferenceIdc(idcId string) (bool, error)
	Init() error
}

type Country interface {
	GetByCode(code int) (res *models.Country, err error)
	GetList() (res []*models.Country, err error)
	Create(c *models.Country) error
	Delete(code int) error
	Init() error
}

type ParamsEnums interface {
	Create(enum *models.CustomParamEnum) (err error)
	Delete(paraName, paramValue string) (err error)
	GetValues(paramName string) (res []*models.CustomParamEnum, err error)
	GetParamList() ([]*models.CustomParamEnum, error)
	Init()error
}

type User interface {
	Get(username string)(*models.User,error)
	GetCount()(int,error)
	Create(m *models.User)error
}



type Domain interface {
	Create(host,group string)error
	GetAll()(res []*models.Domain, err error)
}

func postgresFormat(s string) string {
	return s
}

func createBase(b *models.Base) {
	now := tools.CurrentTimestamp()
	b.CreateAt = now
	b.UpdateAt = now
	b.Id = uid.UUid()
}

func updateBase(b *models.Base) {
	now := tools.CurrentTimestamp()
	b.UpdateAt = now
}

func patchBase(m map[string]interface{}) {
	now := tools.CurrentTimestamp()
	m["update_at"] = now
}

func deleteBase(b *models.Base) {
	b.UpdateAt = tools.CurrentTimestamp()
}


