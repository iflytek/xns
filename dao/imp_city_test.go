package dao

import (
	"database/sql"
	"testing"
)

var city City

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=127.0.0.1 port=55432 dbname=nameserverdev user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	city = NewCityDao(testdb)
}

func Test_impCity_Init(t *testing.T) {
	err := city.Init()
	if err != nil{
		panic(err)
	}
}
