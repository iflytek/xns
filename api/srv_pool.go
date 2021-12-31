package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/uid"
	"sort"
)

type poolService struct {
	PoolDao         dao.Pool
	PoolGroupRefDao dao.GroupPoolRef
	GroupDao        dao.Group
	IdcDao          dao.Idc
}

func newPoolService(poolDao dao.Pool, poolGroupRefDao dao.GroupPoolRef, groupDao dao.Group) *poolService {
	return &poolService{
		PoolDao:         poolDao,
		PoolGroupRefDao: poolGroupRefDao,
		GroupDao:        groupDao,
	}
}

func (ps *poolService) Create(req *Pool) (pool *models.Pool, code int, err error) {
	// 1 检查name是否存在
	pool, err = ps.PoolDao.GetByIdOrName(req.Name)
	if err == nil {
		code = CodeConflict
		err = fmt.Errorf("pool '%s' already exists", req.Name)
		return
	} else {
		if err != dao.NoElemError {
			code = CodeDbError
			return
		}
	}

	pool = &models.Pool{
		Base: models.Base{
			Description: req.Description,
		},
		Name:           req.Name,
		LbMode:         req.LbMode,
		LbConfig:       string(configJsonToString(req.LbConfig)),
		FailOverConfig: string(configJsonToString(req.FailOverConfig)),
	}

	err = ps.PoolDao.Create(pool)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (ps *poolService) Update(id string, req map[string]interface{}) (pool *models.Pool, code int, err error) {
	pool, err = ps.PoolDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "get pool '%s' error :%w", id, err)
		return
	}
	// 修改
	err = ps.PoolDao.Patch(pool.Id, req)
	if err != nil {
		code = CodeDbError
		return
	}

	pool, err = ps.PoolDao.GetById(pool.Id)
	if err != nil {
		code, err = convertErrorf(err, "get pool '%s' error :%w", id, err)
		return
	}
	return
}

func (ps *poolService) Delete(id string) (code int, err error) {
	pool, err := ps.PoolDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
		return
	}

	err = ps.PoolDao.Delete(pool.Id)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
	}
	return
}

func (ps *poolService) GetList() (pools []*models.Pool, code int, err error) {
	pools, err = ps.PoolDao.GetList()
	if err != nil {
		code = CodeDbError
		return
	}

	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Name < pools[j].Name
	})
	return
}

func (ps *poolService) GetPool(id string) (res *models.Pool, code int, err error) {
	res, err = ps.PoolDao.GetByIdOrName(id)
	if err != nil {
		code, err = convertErrorf(err, "get pool '%s' error:%w", id, err)
		return
	}
	return
}

func (ps *poolService) getPoolAndGroup(poolId string, groupId string) (pool *models.Pool, group *models.Group, code int, err error) {
	pool, err = ps.PoolDao.GetByIdOrName(poolId)
	if err != nil {
		code, err = convertErrorf(err, "get pool '%s' error:%w", poolId, err)
		return
	}
	group, err = ps.GroupDao.GetByIdOrName(groupId)
	if err != nil {
		code, err = convertErrorf(err, "get group '%s' error:%w", groupId, err)
		return
	}
	return
}

func (ps *poolService) AddPoolGroup(poolId string, groupId string, weight int) (ref *models.GroupPoolRef, code int, err error) {
	var group *models.Group
	var pool *models.Pool
	pool, group, code, err = ps.getPoolAndGroup(poolId, groupId)
	if err != nil {
		return
	}
	// 检查地址池中是否存在同样的机房
	// 一个地址池中不允许存在
	var groups []*PoolGroups
	groups, code, err = ps.GetPoolGroups(poolId)
	if err != nil {
		return
	}
	for _, poolGroups := range groups {
		if poolGroups.Group.Id == group.Id { // 存在相等的groupId ,更新
			ref = &models.GroupPoolRef{
				Base: models.Base{
					Id: poolGroups.Id,
				},
				GroupId: group.Id,
				PoolId:  pool.Id,
				Weight:  weight,
			}
			err = ps.PoolGroupRefDao.Update(poolGroups.Id, ref)
			if err != nil {
				code = CodeDbError
			}
			return
		}
		if poolGroups.Group.IdcId == group.IdcId {
			code = CodeRequestError
			err = fmt.Errorf("地址池中已经存在属于机房 '%s' 的服务器组了", ps.nameOfIdc(group.IdcId))
			return
		}
	}

	ref = &models.GroupPoolRef{
		Base:    models.Base{},
		GroupId: group.Id,
		PoolId:  pool.Id,
		Weight:  weight,
	}
	err = ps.PoolGroupRefDao.Create(ref)
	if err != nil {
		code = CodeDbError
	}
	return
}

func (ps *poolService) nameOfIdc(id string) string {
	idc ,err := ps.IdcDao.GetByIdOrName(id)
	if err != nil{
		return id
	}
	return idc.Name
}

func (ps *poolService) DeletePoolGroup(poolId string, groupId string) (code int, err error) {
	var group *models.Group
	var pool *models.Pool
	pool, group, code, err = ps.getPoolAndGroup(poolId, groupId)
	if err != nil {
		return
	}
	ref, err := ps.PoolGroupRefDao.GetPoolGroupRef(pool.Id, group.Id)
	if err != nil {
		code, err = convertErrorf(err, "%w", err)
		return
	}
	code, err = nopNotFoundCode(ps.PoolGroupRefDao.Delete(ref.Id))
	return
}

func (ps *poolService) UpdateRef(poolId string, groupId string, weight int) (ref *models.GroupPoolRef, code int, err error) {
	var group *models.Group
	var pool *models.Pool
	pool, group, code, err = ps.getPoolAndGroup(poolId, groupId)
	if err != nil {
		return
	}
	ref, err = ps.PoolGroupRefDao.GetPoolGroupRef(pool.Id, group.Id)
	if err != nil {
		code, err = convertErrorf(err, "get ref error:%w", err)
		return
	}
	ref.Weight = weight
	err = ps.PoolGroupRefDao.Update(ref.Id, ref)
	if err != nil {
		code = CodeDbError
	}
	return
}

type groupWrap struct {
	*models.Group
	Idc *models.Idc `json:"idc"`
}

type PoolGroups struct {
	Id       string        `json:"id"`
	PoolId   string        `json:"pool_id"`
	Group    *models.Group `json:"group"`
	Weight   int           `json:"weight"`
	CreatAt  int           `json:"creat_at"`
	UpdateAt int           `json:"update_at"`
}

func (ps *poolService) GetPoolGroups(pooId string) (res []*PoolGroups, code int, err error) {
	var realPoolIp string
	if uid.IsUUID(pooId) {
		realPoolIp = pooId
	} else {
		var pool *models.Pool
		pool, err = ps.PoolDao.GetByIdOrName(pooId)
		if err != nil {
			code, err = convertErrorf(err, "get pool '%s' error:%w", pooId, err)
			return
		}
		realPoolIp = pool.Id
	}

	var refs []*models.GroupPoolRef
	refs, err = ps.PoolGroupRefDao.GetPoolGroups(realPoolIp)
	if err != nil {
		code = CodeDbError
		return
	}

	res = make([]*PoolGroups, 0, len(refs))
	for _, ref := range refs {
		var group *models.Group
		group, err = ps.GroupDao.GetById(ref.GroupId)
		if err != nil {
			code, err = convertErrorf(err, "get group '%s' error:%w", ref.GroupId, err)
			return
		}

		res = append(res, &PoolGroups{
			Id:       ref.Id,
			PoolId:   ref.PoolId,
			Group:    group,
			Weight:   ref.Weight,
			CreatAt:  ref.CreateAt,
			UpdateAt: ref.UpdateAt,
		})
	}
	return
}

func nopNotFoundCode(e error) (code int, err error) {
	if e == dao.NoElemError {
		return 0, nil
	}
	return CodeDbError, e

}
