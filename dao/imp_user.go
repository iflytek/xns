package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
)

type userDao struct {
	*baseDao
}

func NewUserDao(db *sql.DB)User{
	return &userDao{
		baseDao:newBaseDao(db,&models.User{},"",TableUser),
	}
}

func (u *userDao) Get(username string) (res *models.User,err  error) {
	res = &models.User{}
	err = u.baseDao.queryByCond(newCond().eq("username",username).String(),res)
	return
}

func (u *userDao)GetCount()(int,error)  {
	count ,err := u.getCount("1=1")
	if err != nil{
		return 0, err
	}
	return count,nil
}

func (u *userDao)Create(m *models.User)error{
	createBase(&m.Base)
	return u.insertOnly(m)
}


