package api

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/models"
)

type idcService struct {
	IdcDao      dao.Idc
	RegionDao   dao.Region
	ProvinceDao dao.Province
	CityDao     dao.City
}

func (s *idcService) Create(idc *Idc) (*models.Idc, int, error) {
	_, err := s.IdcDao.GetByName(idc.Name)
	if err != nil {
		if err != dao.NoElemError {
			return nil, CodeDbError, err
		}
	} else {
		return nil, CodeConflict, fmt.Errorf("idc %s has already exists", idc.Name)
	}
	idcm := &models.Idc{
		Name: idc.Name,
	}
	idcm.Description = idc.Desc
	return idcm, 0, s.IdcDao.Create(idcm)
}

func (s *idcService) Update(id string, idc Idc) (idcm *models.Idc, code int, err error) {
	idcm, err = s.IdcDao.GetByIdOrName(id)
	if err != nil {
		if err == dao.NoElemError {
			code = CodeNotFound
			err = fmt.Errorf("update idc error, idc not found,%w", err)
			return
		}
		code = CodeDbError
		return
	}
	idcm.Description = or(idc.Desc, idcm.Description)
	idcm.Name = or(idc.Name, idcm.Description)
	err = s.IdcDao.Update(idcm.Id, idcm)
	if err != nil {
		code = CodeDbError
	}

	return
}

type idcCheckFunc func(idcId string) (bool, error)

func checkIdcRefs(idcId string, funs ...idcCheckFunc) (ok bool, err error) {
	for _, fun := range funs {
		ok, err = fun(idcId)
		if err != nil || ok {
			return
		}
	}
	return false, nil
}

func (s *idcService) Delete(id string) error {
	idc, err := s.IdcDao.GetByIdOrName(id)
	if err != nil {
		if err == dao.NoElemError {
			return nil
		}
		return err
	}
	// 检查 idc机房亲和性  引用
	_, err = checkIdcRefs(idc.Id, s.RegionDao.IfReferenceIdc, s.CityDao.IfReferenceIdc, s.ProvinceDao.IfReferenceIdc)
	if err != nil {
		return err
	}
	//if ok {
	//	return fmt.Errorf("idc is still referenced by region , city or province .cannot delete")
	//}
	return s.IdcDao.Delete(idc.Id)
}

func (s *idcService) Get(idOrName string) (*models.Idc, error) {
	return s.IdcDao.GetByIdOrName(idOrName)
}

func (s *idcService) GetList() ([]*models.Idc, error) {
	return s.IdcDao.GetList()
}

//creater=123
