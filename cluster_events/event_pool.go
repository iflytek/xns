package cluster_events

import (
	"github.com/xfyun/xns/core"
	"github.com/xfyun/xns/dao"
)

type poolEventExecutor struct {
	PoolDao dao.Pool
}

func (p *poolEventExecutor) Create(channel string, data string) error {
	pool, err := p.PoolDao.GetById(data)
	if err != nil {
		return err
	}
	err = core.AddPool(pool)
	return err
}

func (p *poolEventExecutor) Delete(channel string, data string) error {
	return core.DeletePool(data)
}

func (p *poolEventExecutor) Update(channel string, data string) error {
	pool, err := p.PoolDao.GetById(data)
	if err != nil {
		return err
	}
	err = core.AddPool(pool)
	return err
}
