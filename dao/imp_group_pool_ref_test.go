package dao

import (
	"database/sql"
	"github.com/xfyun/xns/models"
	"testing"
)
var (
	gpf GroupPoolRef
)

func init() {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserver user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	gpf = NewGroupPoolRef(testdb)
}

func Test_impGroupPoolRef_Create(t *testing.T) {
	ref := &models.GroupPoolRef{
		Base:    models.Base{},
		GroupId: "b212c32a-416b-452b-aac4-0c1ac43e9aa4",
		PoolId:  "",
	}
	must(gpf.Create(ref))
}
