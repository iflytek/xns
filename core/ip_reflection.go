package core

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

// ip 地址映射，根据ip 地址获取客户端的地域信息
type ipRange struct {
	start uint32
	end   uint32

	cityCodeString     string
	cityCode           int
	provinceCodeString string
	provinceCode       int
	//regionCodeString   string
	regionCode         int
	netCode            string
	countryCode        string
	//
}

type ipReflections struct {
	ips []*ipRange // ip 地址范围，已经根据start 排序
}

// 根据ip地址起始范围排序
func (ip *ipReflections) sort() {
	sort.Slice(ip.ips, func(i, j int) bool {
		return ip.ips[i].start < ip.ips[j].start
	})
}

//使用二分法查询ip地址所在范围
func (ir *ipReflections) getIpRange(ip uint32) *ipRange {
	lo := 0
	hi := len(ir.ips) - 1
	mid := 0
	var mv *ipRange
	for lo <= hi {
		mid = (lo + hi) / 2
		mv = ir.ips[mid]
		if mv.start <= ip && mv.end >= ip {
			return mv
		} else if ip < mv.start {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return nil
}

func (ir *ipReflections) LoadFrom(file string) error {
	rd, err := os.Open(file)
	if err != nil {
		return err
	}
	defer rd.Close()
	ir.ips = make([]*ipRange, 0, 300000)
	sc := bufio.NewScanner(rd)
	for sc.Scan() {
		startIp := 0
		endIp := 0
		countryCode := ""
		provinceCode := 0
		cityCode := 0
		netCode := ""
		n, err := fmt.Sscanf(sc.Text(), "%d %d %s %d %d %s ", &startIp, &endIp, &countryCode, &provinceCode, &cityCode, &netCode)
		if err != nil {
			return fmt.Errorf("scan line '%s'  error:%w", sc.Text(), err)
		}
		if n != 6{
			return  fmt.Errorf("scan line '%s'  error: scan num < 6 ;%d", sc.Text(),n)
		}
		ir.ips = append(ir.ips, &ipRange{
			start:              uint32(startIp),
			end:                uint32(endIp),
			cityCodeString:     strconv.Itoa(cityCode),
			cityCode:           cityCode,
			provinceCodeString: strconv.Itoa(provinceCode),
			provinceCode:       provinceCode,
			//regionCodeString:   "",
			netCode:     netCode,
			countryCode: countryCode,
		})
	}

	ir.sort()
	return nil

}

func (ir *ipReflections) Len() int {
	return len(ir.ips)
}
