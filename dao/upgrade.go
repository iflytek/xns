package dao

import "database/sql"

type upgradeFunc func(db *sql.DB)error


type UpgradeCtrl struct {
	From string
	To string
	Func upgradeFunc
}

// 从0.0.0 版本升级
func upgradeFrom_0_0_0 (db *sql.DB)error{
	_ ,err := db.Exec("create table if not exists  t_user  (username text unique, password text, type text, id uuid primary key, create_at int, update_at int, description text);")
	if err != nil{
		return err
	}
	return nil
}


var upgradeFrom = []UpgradeCtrl{
	{
		From: "0.0.0",
		To:   "1.0.0",
		Func: upgradeFrom_0_0_0,
	},
}
//
//func Upgrade(db *sql.DB)error{
//	ver ,err := currentDBVersion()
//	if err != nil{
//		return err
//	}
//
//
//
//}


func currentDBVersion()(string,error){
	panic("implement me")
}
