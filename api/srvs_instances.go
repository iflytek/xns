package api

import (
	"github.com/xfyun/xns/tools/inject"
)

var (
	idcServiceInstance *idcService     = &idcService{}
	//serverInstance     *serverService  = &serverService{}
	groupInstance      *groupService   = &groupService{}
	poolInstance       *poolService    = &poolService{}
	serviceInstance    *serviceService = &serviceService{}
	regionInstance     *regionService = &regionService{}
	paramInstance  = &paramService{}
	domainInstance = &domainService{}
)

// 自动依赖注入
func Init(daoDeps []interface{}) { // 将dao 注入 services
	inject.Inject([]interface{}{ //services
		idcServiceInstance,
		//serverInstance,
		groupInstance,
		poolInstance,
		serviceInstance,
		regionInstance,
		paramInstance,
		domainInstance,
	}, daoDeps)

}

