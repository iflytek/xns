package cluster_events

import "github.com/xfyun/xns/tools/inject"

func InitClusterEvents(daoDeps []interface{}, pullInterval int) error {
	if pullInterval <= 0 {
		pullInterval = defaultEventPullIntervalSeconds
	}

	cem := &ClusterEventManager{executors: map[string]EventExecutor{}, pullInterval: pullInterval}
	// 将dao 注入
	inject.InjectOne(cem,daoDeps)

	err := cem.Init(daoDeps)
	if err != nil {
		return err
	}

	cem.StartPullEvents()
	return nil
}
