package cluster_events

import (
	"fmt"
	"github.com/xfyun/xns/dao"
	"github.com/xfyun/xns/logger"
	"github.com/xfyun/xns/models"
	"github.com/xfyun/xns/tools/inject"
	"time"
)

const (
	clearExpireIntervalSeconds      = 1 * 24 * 3600
	defaultEventPullIntervalSeconds = 5
)

var events = map[string]EventExecutor{
	dao.ChannelRoute:       new(routeEventExecutor),
	dao.ChannelGroup:       new(groupEventExecutor),
	dao.ChannelGroupServer: new(serverGroupRefExecutor),
	dao.ChannelPool:        new(poolEventExecutor),
	dao.ChannelPoolGroup:   new(groupPoolRefExecutor),
	dao.ChannelService:     new(serviceExecutor),
	dao.ChannelIdc:         new(idcEventExecutor),
	dao.ChannelRegion:      new(regionEventExecutor),
	dao.ChannelCity:        new(cityEventExecutor),
	dao.ChannelProvince:    new(provinceEventExecutor),
}

// 集群更新事件相关
type EventExecutor interface {
	Create(channel string, data string) error
	Delete(channel string, data string) error
	Update(channel string, data string) error
}

type ClusterEventManager struct {
	executors    map[string]EventExecutor // key : channel
	Dao          *dao.ClusterEventDao
	currentId    int
	pullInterval int // seconds
}

//@daoDeps: dao 实例依赖
func (c *ClusterEventManager) Init(daoDeps []interface{}) error {
	id, err := c.Dao.GetMaxEventId()
	if err != nil {
		return fmt.Errorf("get max id error:%w",err)
	}

	c.currentId = id
	logger.Event().Info("start init event, current eventId is:", id)
	for channel, executor := range events {
		inject.InjectOne(executor, daoDeps) // 将dao 依赖注入 event 实例
		c.RegisterEventExecutor(channel, executor)
	}

	return nil
}

func (c *ClusterEventManager) StartPullEvents() {
	fmt.Println("about to start event pull task and event clear task,pull interval is ", c.pullInterval)
	go c.startPullTask()
	go c.startClearExpireEventsTask()
}

func (c *ClusterEventManager) RegisterEventExecutor(channel string, exec EventExecutor) {
	if c.executors[channel] != nil {
		panic("event executor channel of '" + channel + "' has already registered")
	}
	c.executors[channel] = exec
}

func (c *ClusterEventManager) doEvent(event *models.ClusterEvent) (err error) {
	exec := c.executors[event.Channel]
	if exec == nil {
		return fmt.Errorf("executor %s not found", event.Channel)
	}
	switch event.Event {
	case dao.EventDelete:
		return exec.Delete(event.Channel, event.Data)
	case dao.EventUpdate:
		return exec.Update(event.Channel, event.Data)
	case dao.EventCreate:
		return exec.Create(event.Channel, event.Data)
	default:
		return fmt.Errorf("unknow event type:%s", event.Event)
	}
}

func (c *ClusterEventManager) startPullTask() {

	tick := time.NewTicker(time.Duration(c.pullInterval) * time.Second)
	for range tick.C {
		c.pullAndDoEvents()
	}
}

//定时从数据中获取数据更新事件
func (c *ClusterEventManager) pullAndDoEvents() {
	events, err := c.Dao.PullNewEvent(c.currentId)
	if err != nil {
		logger.Event().Error("pull event error:", err)
		return
	}
	c.currentId = c.maxIdOfEvents(events)
	for _, event := range events {
		err = c.doEvent(event)
		if err != nil {
			logger.Event().Error("do event error:", err, " event:", event)
		} else {
			logger.Event().Info("success handle event,", "event", event)
		}
	}
}

// 定期清理数据库中过期无用的cluster_event
func (c *ClusterEventManager) startClearExpireEventsTask() {
	ticker := time.NewTicker(time.Second * clearExpireIntervalSeconds)
	for range ticker.C {
		err := c.Dao.ClearEvents()
		if err != nil {
			logger.Event().Error("clear expire event error:", err)
		} else {
			logger.Event().Info("success clear expire events")
		}
	}
}

func (c *ClusterEventManager) maxIdOfEvents(events []*models.ClusterEvent) int {
	if len(events) == 0 {
		return c.currentId
	}
	id := -1
	for _, event := range events {
		if event.Id > id {
			id = event.Id
		}
	}
	return id
}
