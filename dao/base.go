package dao

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

const (
	TableIdc             = "t_idc"
	TableServer          = "t_server"
	TableGroup           = "t_group"
	TableServerGroupRef  = "t_server_group_ref"
	TablePool            = "t_pool"
	TableGroupPoolRef    = "t_group_pool_ref"
	TableRoute           = "t_route"
	TableService         = "t_service"
	TableCity            = "t_city"
	TableProvince        = "t_province"
	TableRegion          = "t_region"
	TableCountry         = "t_country"
	TableParamValueEnums = "t_custom_param_enum"
	TableUser            = "t_user"
)

type baseDao struct {
	db          *sql.DB
	queryFields string
	format      FiledNameFormat
	model       interface{}
	channel     string
	table       string
}

//@model 数据库模型
//@channel channel
func newBaseDao(db *sql.DB, model interface{}, channel, table string) *baseDao {
	b := &baseDao{
		db:      db,
		model:   model,
		channel: channel,
		table:   table,
	}
	b.format = postgresFormat
	b.queryFields = strings.Join(getSqlQueryTags(reflect.TypeOf(model), b.format), ",")
	return b
}

func (base *baseDao) execTx(tx *sql.Tx, sql string) error {
	_, err := tx.Exec(sql)
	if err != nil {
		return fmt.Errorf("execTx sql error,sql=%s,err=%w", sql, err)
	}
	return nil
}

func (base *baseDao) execTxWithAffc(tx *sql.Tx, sql string) (sql.Result, error) {
	res, err := tx.Exec(sql)
	if err != nil {
		return nil, fmt.Errorf("execTx sql error,sql=%s,err=%w", sql, err)
	}
	return res, nil
}

func (base *baseDao) exec(sql string) error {
	_, err := base.db.Exec(sql)
	if err != nil {
		return fmt.Errorf("exec sql error,sql=%s,err=%w", sql, err)
	}
	return nil
}

func (base *baseDao) query(sqls string) (r *sql.Rows, err error) {
	r, err = base.db.Query(sqls)
	if err != nil {
		return nil, fmt.Errorf("query error,sql=%s,err=%w", sqls, err)
	}
	return r, nil
}

func (base *baseDao) queryResults(sqls string, res interface{}) (err error) {
	rows, err := base.query(sqls)
	if err != nil {
		return err
	}
	return unmarshalFromRows(rows, res)
}

func (base *baseDao) queryByCond(cond string, res interface{}) (err error) {
	sqlString := fmt.Sprintf("select %s from %s where %s", base.queryFields, base.table, cond)
	return base.queryResults(sqlString, res)
}

func (base *baseDao) insert(tx *sql.Tx, i interface{}) error {
	sqlString := generateInsertSql(base.table, base.format, i)
	return base.execTx(tx, sqlString)
}

func (base *baseDao) insertOnly(i interface{}) error {
	sqlString := generateInsertSql(base.table, base.format, i)
	return base.exec(sqlString)
}

func (base *baseDao) insertAndSendEvent(i interface{}, eventData string) error {
	tx, err := base.beginTx()
	if err != nil {
		return err
	}
	if err = base.insert(tx, i); err != nil {
		return err
	}

	err = addClusterEvent(tx, EventCreate, base.channel, eventData)
	if err != nil {
		base.rollbackTx(tx)
		return err
	}
	return tx.Commit()
}

func (base *baseDao) updateAndSendEvent(cond string, i interface{}, eventData string) error {
	tx, err := base.beginTx()
	if err != nil {
		return err
	}
	if err = base.update(tx, cond, i); err != nil {
		return err
	}
	err = addClusterEvent(tx, EventUpdate, base.channel, eventData)
	if err != nil {
		base.rollbackTx(tx)
		return err
	}
	return tx.Commit()
}

func (base *baseDao) getCount(cond string) (int, error) {
	sqlString := fmt.Sprintf("select count(*) from %s where %s", base.table, cond)
	res, err := base.query(sqlString)
	if err != nil {
		return 0, err
	}
	defer res.Close()
	for res.Next() {
		var a int
		err = res.Scan(&a)
		if err != nil {
			return 0, err
		}
		return a, nil
	}
	return 0, err
}

func (base *baseDao) deleteAndSendEvent(cond string, eventData string) error {
	count, err := base.getCount(cond)
	if err != nil {
		return err
	}
	if count == 0 {
		return NoElemError
	}

	tx, err := base.beginTx()
	if err != nil {
		return err
	}
	sqlString := fmt.Sprintf("delete from %s where %s", base.table, cond)
	if err = base.execTx(tx, sqlString); err != nil {
		return err
	}
	err = addClusterEvent(tx, EventDelete, base.channel, eventData)
	if err != nil {
		base.rollbackTx(tx)
		return err
	}
	return tx.Commit()
}

func (base *baseDao) deleteOnly(cond string) error {
	sqlString := fmt.Sprintf("delete from %s where %s", base.table, cond)
	return base.exec(sqlString)
}

//@cond 更新条件
func (base *baseDao) update(tx *sql.Tx, cond string, i interface{}) error {
	sql := generateUpdateSql(base.format, base.table, cond, i)
	return base.execTx(tx, sql)
}

func (base *baseDao) beginTx() (*sql.Tx, error) {
	tx, err := base.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx error:%w", err)
	}
	return tx, nil
}

func (base *baseDao) queryAll(res interface{}) error {
	sql := fmt.Sprintf("select %s from %s ;", base.queryFields, base.table)
	return base.queryResults(sql, res)
}

func (i *baseDao) rollbackTx(tx *sql.Tx) {
	tx.Rollback()
}

func (i *baseDao) queryCount(cond string) (int, error) {
	sql := fmt.Sprintf("select count(*) from %s where %s ;", i.table, cond)
	rows, err := i.query(sql)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			return c, err
		}
		return c, nil
	}
	return 0, err
}

func generateIdCond(id string) string {
	return fmt.Sprintf("id = '%s'", id)
}

type cond struct {
	cond strings.Builder
}

func newCond() *cond {
	return &cond{}
}

func (c *cond) eq(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" = ")
	c.cond.WriteString(valOfsqlSlice(val))
	return c
}

func (c *cond) contains(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" like ")
	c.cond.WriteString(fmt.Sprintf("'%%%v%%'", val))
	return c
}

func (c *cond) gt(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" > ")
	c.cond.WriteString(valOfsqlSlice(val))
	return c
}

func (c *cond) lt(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" < ")
	c.cond.WriteString(valOfsqlSlice(val))
	return c
}

func (c *cond) le(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" <= ")
	c.cond.WriteString(valOfsqlSlice(val))
	return c
}

func (c *cond) ge(name string, val interface{}) *cond {
	c.cond.WriteString(name)
	c.cond.WriteString(" >= ")
	c.cond.WriteString(valOfsqlSlice(val))
	return c
}

func (c *cond) and() *cond {
	c.cond.WriteString(" and ")
	return c
}

func (c *cond) or() *cond {
	c.cond.WriteString(" or ")
	return c
}

func (c *cond) truth() *cond {
	c.cond.WriteString(" 1 = 1")
	return c
}

func (c *cond) String() string {
	return c.cond.String()
}

func (c *cond) queryConds(kvs ...condKV) *cond {
	vs := make([]condKV, 0)
	for _, kv := range kvs {
		if kv.val == "" {
			continue
		}
		vs = append(vs, kv)
	}

	switch len(vs) {
	case 0:
		c.truth()
		return c
	case 1:
		vs[0].cond(vs[0].key, vs[0].val)
	default:
		vs[0].cond(vs[0].key, vs[0].val)
		for _, kv := range vs[1:] {
			c.and()
			kv.cond(kv.key, kv.val)
		}
	}
	return c
}

type condKV struct {
	key  string
	val  interface{}
	cond func(k string, v interface{}) *cond
}

func newIdCond(id string) string {
	return fmt.Sprintf("id = '%s'", id)
}

func newCodeCond(code int) string {
	return fmt.Sprintf("code = '%d'", code)
}

func newNameCond(id string) string {
	return fmt.Sprintf("name = '%s'", id)
}

func valOfsqlSlice(i interface{}) string {
	switch v := i.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case int, int32, int64, uint, uint64:
		return fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("'%v'", i)
}
