package core

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
)

type Daoes struct {
	ServiceDao dao.Service
	//ServerDao         dao.Server
	GroupServerRefDao dao.GroupServerRef
	GroupDao          dao.Group
	PoolDao           dao.Pool
	GroupPoolRefDao   dao.GroupPoolRef
	RouteDao          dao.Route
	IdcDao            dao.Idc
	City              dao.City
	Province          dao.Province
	Region            dao.Region
	Param             dao.ParamsEnums
	Country           dao.Country
	User              dao.User
}

func (d *Daoes) InitCmds() map[string]func() error {
	return map[string]func() error{
		"city":     d.City.Init,
		"province": d.Province.Init,
		"region":   d.Region.Init,
		"params":   d.Param.Init,
		"country":  d.Country.Init,
	}
}

type Args struct {
	Dao            *Daoes
	IpResourcePath string
}

func Init(arg Args) error {
	daoes := arg.Dao
	// 初始化ip地址池
	err := gIpReflection.LoadFrom(arg.IpResourcePath)
	if err != nil {
		return err
	}

	if gIpReflection.Len() < 1 {
		return fmt.Errorf("load ip resorce error ,may load wrong ip resource file, a valid length of ip resoure should be larger than 200000,but now is %d", gIpReflection.Len())
	}

	idcs, err := daoes.IdcDao.GetList()
	if err != nil {
		return err
	}

	//servers, err := daoes.ServerDao.GetList()
	//if err != nil {
	//	return err
	//}
	refs, err := daoes.GroupServerRefDao.GetList()
	if err != nil {
		return err
	}
	// add group
	groups, err := daoes.GroupDao.GetList()
	if err != nil {
		return err
	}

	// add Pool
	pools, err := daoes.PoolDao.GetList()
	if err != nil {
		return err
	}

	prefs, err := daoes.GroupPoolRefDao.GetList()

	services, err := daoes.ServiceDao.GetList()
	if err != nil {
		return err
	}
	routes, err := daoes.RouteDao.GetList()
	if err != nil {
		return err
	}

	regions, err := daoes.Region.GetList()
	if err != nil {
		return err
	}

	provinces, err := daoes.Province.GetList()
	if err != nil {
		return err
	}

	cities, err := daoes.City.GetList()
	if err != nil {
		return err
	}
	return initData(refs, groups, prefs, pools, routes, services, idcs, regions, provinces, cities)
}

func initData(sgref []*models.ServerGroupRef,
	groups []*models.Group, pgref []*models.GroupPoolRef,
	pools []*models.Pool, routes []*models.Route, services []*models.Service, idcs []*models.Idc,
	regions []*models.Region, provinces []*models.Province, cities []*models.City) error {

	for _, idc := range idcs {
		AddIdc(idc)
	}

	var err error
	for _, ref := range sgref {
		if err = AddGroupServerRef(ref); err != nil {
			return err
		}
	}

	for _, group := range groups {
		if err = AddGroup(group); err != nil {
			return err
		}
	}

	for _, ref := range pgref {
		if err = AddPoolGroupRef(ref); err != nil {
			return err
		}
	}

	for _, pool := range pools {
		err = AddPool(pool)
		if err != nil {
			return err
		}
	}

	for _, service := range services {
		if err = AddService(service); err != nil {
			return err
		}
	}

	for _, route := range routes {
		err = AddRoute(route)
		if err != nil {
			return err
		}
	}

	for _, region := range regions {
		if err = AddRegion(region); err != nil {
			return err
		}
	}

	for _, prov := range provinces {
		if err = AddProvince(prov); err != nil {
			return err
		}
	}

	for _, city := range cities {
		if err = AddCity(city); err != nil {
			return err
		}
	}

	return err
}
