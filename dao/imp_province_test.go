package dao

import (
	"database/sql"
	"testing"
)
//
var prov Province

func init() {

	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserverdev user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	prov = NewProvinceDao(testdb)
}

func Test_impProvince_Init(t *testing.T) {
	must(prov.Init())
}
