package dao

import (
	"database/sql"
	"fmt"
	"testing"
)

var(
	route Route
)

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserver user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	route = NewRouteDao(testdb)
}



func Test_impRoute_QueryRoutes(t *testing.T) {
	fmt.Println(route.QueryRoutes("","appid=500"))
	fmt.Println(route.QueryRoutes("evo",""))
}
