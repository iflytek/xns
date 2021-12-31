package dao

import (
	"database/sql"
	"fmt"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools"
)

const (
	expireSeconds = 1 * 24 * 3600 // 事件在数据库中的保存时间
)
const (
	EventCreate = "create"
	EventUpdate = "update"
	EventDelete = "delete"
)

const (
	ChannelIdc         = "idc"
	ChannelServer      = "server"
	ChannelGroup       = "group"
	ChannelGroupServer = "group_server_ref"
	ChannelPool        = "pool"
	ChannelPoolGroup   = "pool_group_ref"
	ChannelRoute       = "route"
	ChannelService     = "service"
	ChannelRegion      = "region"
	ChannelProvince    = "province"
	ChannelCity        = "city"
	ChannelCountry     = "country"
	ChannelParams      = "parameters"
)

var (
	clusterEventTableName = "t_cluster_event"
)

// 向cluster event 表中新增一条更新事件
func addClusterEvent(tx *sql.Tx, event, channel string, data string) error {
	now := tools.CurrentTimestamp()
	sqlString := fmt.Sprintf("insert into %s (event,channel,data,at,expire_at)values('%s','%s','%s',%d,%d)",
		clusterEventTableName, event, channel, data, now, now+expireSeconds)
	_, err := tx.Exec(sqlString)
	if err != nil {
		return fmt.Errorf("insert into cluster event error:sql=%s,err=%w", sqlString, err)
	}
	return nil
}

func getEventCount(db *sql.DB) (int, error) {
	countSql := fmt.Sprintf("select count(*) from %s", clusterEventTableName)
	rows, err := db.Query(countSql)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		c := 0
		err = rows.Scan(&c)
		if err != nil {
			return 0, err
		}
		return c, nil
	}

	return 0, fmt.Errorf("no count found")
}

// 启动时，先获取最大的id，让后加载数据库，最后从最大id为起始点，开始拉取事件。来避免事件被丢失，和重复加载，
// 事件重复加载应该是允许的。
func GetMaxEventId(db *sql.DB) (int, error) {
	n, err := getEventCount(db)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return -1, nil
	}

	sqlString := fmt.Sprintf("select max(id) from %s;", clusterEventTableName)
	rows, err := db.Query(sqlString)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return 0, err
		}
		return id, nil
	}
	return -1, nil
}

type ClusterEventDao struct {
	*baseDao
	table string
}

func NewClusterEventDao(db *sql.DB) *ClusterEventDao {
	return &ClusterEventDao{
		baseDao: newBaseDao(db, &models.ClusterEvent{}, "", clusterEventTableName),
		table:   clusterEventTableName,
	}
}

func (c *ClusterEventDao) PullNewEvent(fromId int) (res []*models.ClusterEvent, err error) {
	sqlString := fmt.Sprintf("select %s from %s where id > %d", c.queryFields, c.table, fromId)
	err = c.queryResults(sqlString, &res)
	return
}

func (c *ClusterEventDao) GetMaxEventId() (int, error) {
	return GetMaxEventId(c.db)
}

func (c *ClusterEventDao) ClearEvents() error {
	sqlString := fmt.Sprintf("delete from %s where expire_at < %d", c.table, tools.CurrentTimestamp())
	return c.exec(sqlString)
}
