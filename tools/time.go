package tools

import "time"
// 返回当前的时间戳，精确到ms
func CurrentTimestamp()int{
	return int(time.Now().Unix())
}
