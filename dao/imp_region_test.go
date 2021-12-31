package dao

import (
	"database/sql"
	"testing"
)
var region Region
func init() {

	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserverdev user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	region = NewRegionDao(testdb)
}

func Test_impRegion_Init(t *testing.T) {
	must(region.Init())
}
