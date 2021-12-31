package dao

import (
	"database/sql"
	"fmt"
)

func Init(connectionUrl string) ([]interface{}, *sql.DB, error) {
	db, err := sql.Open("postgres", connectionUrl)
	if err != nil {
		return nil, nil, fmt.Errorf("init database error:%w", err)
	}

	daos := []interface{}{
		NewIdcImp(db),
		NewGroupImp(db),
		NewGroupServerRef(db),
		NewPool(db),
		NewGroupPoolRef(db), //dependencies
		NewServiceDao(db),
		NewRouteDao(db),
		NewClusterEventDao(db),
		NewCityDao(db),
		NewProvinceDao(db),
		NewRegionDao(db),
		NewCountryDao(db),
		NewParamDao(db),
		NewUserDao(db),
		NewDomainDao(db),
	}
	return daos, db, nil
}
