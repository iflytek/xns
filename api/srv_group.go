package api

import (
	"fmt"
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"sort"
)

type groupService struct {
	GroupDao          dao.Group
	IdcDao            dao.Idc
	ServerGroupRefDao dao.GroupServerRef
}

func (g *groupService) Create(gp *AddGroupReq) (group *models.Group, code int, err error) {
	//1 检查名称是否已经存在
	_, err = g.GroupDao.GetByName(gp.Name)
	if err == nil {
		code = CodeNotFound
		err = fmt.Errorf("group '%s' already exists", gp.Name)
		return
	} else {
		if err != dao.NoElemError {
			code = CodeDbError
			return
		}
	}
	//
	var idc *models.Idc
	idc, err = g.IdcDao.GetByIdOrName(gp.IdcId)
	if err != nil {
		err = fmt.Errorf("get idc '%s' error", gp.IdcId)
		code = CodeDbError
		return
	}

	group = &models.Group{
		Base: models.Base{
			Description: gp.Desc,
		},
		Name:               gp.Name,
		IdcId:              idc.Id,
		HealthyCheckMode:   gp.HealthyCheckMode,
		HealthyCheckConfig: configJsonToString(gp.HealthyCheckConfig),
		HealthyNum:         gp.HealthyNum,
		UnHealthyNum:       gp.UnHealthyNum,
		HealthyInterval:    gp.HealthyInterval,
		UnHealthyInterval:  gp.UnHealthyInterval,
		LbMode:             gp.LbMode,
		//LbConfig:           configJsonToString(gp.LbConfig),
		ServerTags:     configJsonToString(gp.ServerTags),
		//Weight:         gp.Weight,
		IpAllocNum:     gp.IpAllocNum,
		DefaultServers: gp.DefaultServers,
		Port:           gp.Port,
	}

	code, err = g.checkHealthConfig(group)
	if err != nil {
		return
	}

	err = g.GroupDao.Create(group)
	if err != nil {
		code = CodeDbError
	}
	return
}

func reqObject2String(req map[string]interface{}, keys ...string) {
	for _, key := range keys {
		val, ok := req[key]
		if ok {
			req[key] = configJsonToString(val)
		}
	}
}

func (g *groupService) Update(id string, req map[string]interface{}) (group *models.Group, code int, err error) {
	group, err = g.GroupDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "patch group error:%w", err)
		return
	}
	// object 转换为json string 类型放入数据库
	reqObject2String(req, "healthy_check_config")
	reqObject2String(req, "lb_config")
	reqObject2String(req, "server_tags")
	delete(req, "idc_id") // 机房不允许更改
	// 将req 写到 group 中
	err = patch(group, req)
	if err != nil {
		code = CodeRequestError
		return
	}

	code, err = g.checkHealthConfig(group)
	if err != nil {
		return
	}
	err = g.GroupDao.Update(group.Id, group)
	if err != nil {
		code, err = convertErrorf(err, "patch group p error:%w", err)
		return
	}
	group, err = g.GroupDao.GetByIdOrName(group.Id)
	if err != nil {
		code, err = convertErrorf(err, "patch group  get error:%w", err)
		return
	}
	return
}

func (g *groupService) checkHealthConfig(group *models.Group) (code int, err error) {
	_, err = core.NewHealthyCheck(group.HealthyCheckMode, string(group.HealthyCheckConfig))
	if err != nil {
		code = CodeRequestError
		err = fmt.Errorf("parse healthy check config error")
		return
	}
	return
}

func (g *groupService) Get(id string) (group *models.Group, code int, err error) {
	group, err = g.GroupDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "get group error:%w", err)
		return
	}
	return
}

func (g *groupService) GetList() (group []*models.Group, code int, err error) {
	group, err = g.GroupDao.GetList()
	if err != nil {
		code, err = convertErrorf(err, "get group error:%w", err)
		return
	}
	sort.Slice(group, func(i, j int) bool {
		return group[i].Name < group[j].Name
	})
	return
}

func (g *groupService) Delete(id string) (code int, err error) {
	var group *models.Group
	group, err = g.GroupDao.GetByIdOrName(id)
	if err != nil {
		code,err = convertErrorf(err,"%w",err)
		return
	}
	err = g.GroupDao.Delete(group.Id)
	if err != nil {
		code ,err = convertErrorf(err,"%w",err)
		return
	}
	return
}

func (g *groupService) AddServer(groupId, serverIp string, weight int) (ref *models.ServerGroupRef, code int, err error) {
	if weight <= 0 {
		code = CodeRequestError
		err = fmt.Errorf("weight must be lager than 0")
		return
	}

	var group *models.Group
	group, code, err = g.getServerAndGroup(groupId)
	if err != nil {
		return
	}

	ref ,err = g.ServerGroupRefDao.GetByRef(serverIp,group.Id)
	if err == nil{
		ref.Weight = weight
		err = g.ServerGroupRefDao.Update(ref.Id,ref)
		if err != nil{
			code,err = convertErrorf(err,"request db error:%w",err)
		}
		return
	}else{
		ref = &models.ServerGroupRef{
			Base:     models.Base{},
			ServerIp: serverIp,
			GroupId:  group.Id,
			Weight:   weight,
		}
	}
	err = g.ServerGroupRefDao.Create(ref)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (g *groupService) getServerAndGroup(groupId string) (group *models.Group, code int, err error) {

	group, err = g.GroupDao.GetByIdOrName(groupId)
	if err != nil {
		code, err = convertErrorf(err, "get group '%s' error:%w", groupId, err)
		return
	}
	return
}

func (g *groupService) UpdateServer(groupId, serverIp string, weight int) (code int, err error) {
	if weight <= 0 {
		code = CodeRequestError
		err = fmt.Errorf("weight must be larger than 0")
		return
	}
	var group *models.Group
	group, code, err = g.getServerAndGroup(groupId)
	if err != nil {
		return
	}

	var ref *models.ServerGroupRef
	ref, err = g.ServerGroupRefDao.GetByRef(serverIp, group.Id)
	if err != nil {
		if err == dao.NoElemError {
			return 0, nil
		}
		code = CodeDbError
		return
	}

	ref.Weight = weight
	err = g.ServerGroupRefDao.Update(ref.Id, ref)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (g *groupService) DeleteServer(groupId, serverIp string) (code int, err error) {
	var group *models.Group
	group, code, err = g.getServerAndGroup(groupId)
	if err != nil {
		return
	}
	var ref *models.ServerGroupRef
	ref, err = g.ServerGroupRefDao.GetByRef(serverIp, group.Id)
	if err != nil {
		code,err = convertErrorf(err,"%w",err)
		return
	}

	err = g.ServerGroupRefDao.Delete(ref.Id)
	if err != nil {
		code ,err = convertErrorf(err,"%w",err)
	}
	return
}

type GroupServer struct {
	GroupId  string `json:"group_id"`
	ServerIp string `json:"server"`
	Weight   int    `json:"weight"`
	CreateAt int `json:"create_at"`
	UpdateAt int `json:"update_at"`
	IsHealthy bool `json:"is_healthy"`
}

func (g *groupService) GetGroupServers(groupId string) (res []*GroupServer, code int, err error) {
	var group *models.Group
	group, err = g.GroupDao.GetByIdOrName(groupId)
	if err != nil {
		code, err = convertErrorf(err, "get group '%s' error:%w", groupId, err)
		return
	}
	var refs []*models.ServerGroupRef
	refs, err = g.ServerGroupRefDao.GetGroupServers(group.Id)
	if err != nil {
		code = CodeDbError
		return
	}
	status := core.GetGroupStatus()
	statMap := make(map[string]bool)
	for _, groupStatus := range status {
		if groupStatus.Name == group.Name{
			for _, node := range groupStatus.Nodes {
				statMap[node.Addr] = node.Healthy
			}
		}
	}
	res = make([]*GroupServer, len(refs))
	for i, ref := range refs {
		gs := &GroupServer{
			GroupId:  ref.GroupId,
			ServerIp: ref.ServerIp,
			Weight:   ref.Weight,
			CreateAt: ref.CreateAt,
			UpdateAt: ref.UpdateAt,
		}
		res[i] = gs
		isHealthy := true
		if group.HealthyNum > 0 {
			isHealthy = statMap[ref.ServerIp]
		}
		gs.IsHealthy = isHealthy

	}
	return

}
