package generate_env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/template"
)

type group struct {
	id   string
	name string
}

type server struct {
	Id   string
	Ip   string
	Port string
	Idc  string
}

func readGroup() map[string]*group {
	file, err := os.Open("./group.txt")
	must(err)
	sc := bufio.NewScanner(file)
	res := make(map[string]*group)
	sc.Scan()
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		ksvs := strings.Fields(sc.Text())
		if len(ksvs) != 2 {
			panic("invalid group line:" + sc.Text())
		}
		res[ksvs[0]] = &group{
			id:   ksvs[0],
			name: ksvs[1],
		}
	}
	return res
}

func readServers() (res map[string]*server) {
	file, err := os.Open("./server.txt")
	must(err)
	sc := bufio.NewScanner(file)
	res = map[string]*server{}
	sc.Scan()
	for sc.Scan() {
		if sc.Text() == "" {
			continue
		}
		ksvs := strings.Fields(sc.Text())
		if len(ksvs) != 4 {
			panic("invalid server line:" + sc.Text())
		}
		res[ksvs[0]] = &server{
			Id:   ksvs[0],
			Ip:   ksvs[1],
			Port: ksvs[2],
			Idc:  ksvs[3],
		}
	}
	return res
}

type groupServers struct {
	groupId   string
	serversId []string
	servers   []string
}

func readGroupServers() (res map[string]*groupServers) {
	file, err := os.Open("./group_server.txt")
	must(err)
	sc := bufio.NewScanner(file)
	res = map[string]*groupServers{}
	sc.Scan()
	for sc.Scan() {
		ksvs := strings.Fields(sc.Text())
		if sc.Text() == "" {
			continue
		}
		if len(ksvs) != 2 {
			panic("invalid server line:" + sc.Text())
		}
		res[ksvs[0]] = &groupServers{
			groupId:   ksvs[0],
			serversId: strings.Split(ksvs[1], ","),
			servers:   nil,
		}
	}
	return res
}

type param struct {
	Key string
	Val string
}

type rule struct {
	Host   string
	Params []param
	Rank int
}

type rules struct {
	groupId string
	rank    string
	rule    string
	ruleRaw string
	rules   rule
}

func parseRule(r string) (rl rule) {
	for _, s := range strings.Split(r, "&") {
		kvs := strings.SplitN(s, "=", 2)
		if len(kvs) != 2 {
			panic("invalie rule:" + r)
		}
		if kvs[0] == "host" {
			rl.Host = kvs[1]
			continue

		}
		rl.Params = append(rl.Params, param{
			Key: kvs[0],
			Val: kvs[1],
		})
	}
	return rl
}

func readRules() map[string][]*rules {
	file, err := os.Open("./rule.txt")
	must(err)
	sc := bufio.NewScanner(file)
	res := map[string][]*rules{}
	sc.Scan()
	for sc.Scan() {
		ksvs := strings.Fields(sc.Text())
		if sc.Text() == "" {
			continue
		}
		if len(ksvs) != 3 {
			panic("invalid rule line:" + sc.Text())
		}
		rr := &rules{
			groupId: ksvs[0],
			rank:    ksvs[1],
			rule:    ksvs[2]+" ,rank= "+ksvs[1],
			ruleRaw: ksvs[2],
			rules:   parseRule(ksvs[2]),
		}
		rk ,err  := strconv.Atoi( ksvs[1])
		if err != nil{
			panic(err)
		}
		rr.rules.Rank =rk
		res[ksvs[0]] = append(res[ksvs[0]], rr)
	}
	return res
}

type elem struct {
	Group   string
	Servers string
	Rules   string
}

const (
	tlp = `
### 规则和group 整理

group | servers | rules (is_used=1)
-----| ------ | ------
{{range $e:=. -}}
{{$e.Group}}|{{$e.Servers}}|{{$e.Rules}}
{{ end }}
`
// b appod=5
	tlp_host = `

<table>
	<tr>
		<th>host</th>
		<th>group</th>
		<th>rule</th>
		<th>ips</th>
	</tr>
{{range $e:=. -}}
	<tr>
		<td rowspan="{{$e.Len}}"> {{$e.Host}}</td>
		<td>{{$e.G.Name}}</td>
		<td>{{$e.G.Rules}}</td>
		<td>{{$e.G.Ips}}</td>
	</tr>
	{{range $ee:=$e.Groups -}}
	<tr>
		<td>{{$ee.Name}} </td>
		<td>{{$ee.Rules}} </td>
		<td>{{$ee.Ips}} </td>
	</tr>

	{{- end}}
{{- end }}
</table>
`
)

func TestMain2(t *testing.T) {
	srvs := readGroupServers()
	serverMap := readServers()
	rule := readRules()
	elems := []elem{}
	gps := readGroup()
	gIps := map[string][]string{}
	gsrvs := map[string][]*server{}
	for _, g := range gps {
		ss := make([]string, 0)
		srvvs := []*server{}
		for _, sid := range srvs[g.id].serversId {
			srv := serverMap[sid]

			ss = append(ss, fmt.Sprintf("%s:%s %s", srv.Ip, srv.Port, srv.Idc))
			srvvs = append(srvvs,srv)
		}
		gsrvs[g.name] = srvvs
		gIps[g.id] = ss
		//fmt.Println(g.id,g.name,ss,rule[g.id])
		//fmt.Println("group:",g.id,g.name)
		//fmt.Printf("servers:%v\n",ss)
		rs := make([]string, 0)
		for _, r := range rule[g.id] {
			//fmt.Println(r.rule)
			rs = append(rs, r.rule)
		}
		elems = append(elems, elem{
			Group:   g.name,
			Servers: strings.Join(ss, "<br/>"),
			Rules:   strings.Join(rs, "<br/>"),
		})
		//fmt.Println("-----------------------------------------------------------------------------------------------")
	}
	tlp, err := template.New("").Parse(tlp)
	if err != nil {
		panic(err)
	}
	output, _ := os.Create("output.md")
	err = tlp.Execute(output, elems)
	must(err)
	hostsRules := map[string]map[string]*groupElem{}
	for _, r := range rule {
		for _, r2 := range r {
			hm := hostsRules[r2.rules.Host]
			if hm == nil{
				hm = map[string]*groupElem{}
				hostsRules[r2.rules.Host] = hm
			}
			gname := gps[r2.groupId].name
			gname = strings.TrimSuffix(gname,"\\t")
			g := hm[gname]
			if g== nil{
				g = &groupElem{
					Name:  gname,
					Rule:  nil,
					Ips:   strings.Join(gIps[r2.groupId],"</br>"),
				}
				hm[gname] =g
			}
			g.Rule = append(g.Rule,r2.rule)
			g.RuleRaw = append(g.RuleRaw,&r2.rules)
			//grs[r2.groupId].Rule = append(grs[r2.groupId].Rule ,r2.rule)
			//hostsRules[r2.rules.host] = append(hostsRules[r2.rules.host], fmt.Sprintf("%s %s", r2.rule, gps[r2.groupId].name))
		}

	}

	hoses := []hostelem{}
	for s, i := range hostsRules {
		gs := []*groupElem{}
		for _, g := range i {
			gs = append(gs,g)
			g.Rules = strings.Join(g.Rule,"<br/>")
		}
		if len(gs) > 0{
			hoses = append(hoses,hostelem{
				Host:   s,
				Groups:gs[1:],
				G: gs[0],
				Len: len(gs),
			})
		}

		fmt.Println("host:", s)
		for _, s2 := range i {
			fmt.Println(s2)

		}
		fmt.Println("--------------------------------------")
	}
	hostOut,_ := os.Create("hostout.md")
	tp ,err := template.New("2").Parse(tlp_host)
	must(err)
	must(tp.Execute(hostOut,hoses))


	// appids
	kvs := make(map[string]map[string]bool)
	for _, i := range rule {
		for _, r := range i {
			for _, p := range r.rules.Params {
				m := kvs[p.Key]
				if m== nil{
					m = map[string]bool{}
					kvs[p.Key] = m
				}
				m[p.Val] = true
			}
		}
	}

	for key, m := range kvs {
		for param, _ := range m {
			if in(key,"net","country","city","province"){
				continue
			}
			fmt.Println(key,param)
		}
	}
	ggsrbs := make(map[string][]*server)
	for _, hose := range hoses {
		for _, g := range hose.Groups {
			if strings.HasSuffix(g.Name,"\\t"){
				continue
			}
			ggsrbs [g.Name] = gsrvs[g.Name]

		}
		if hose.G != nil{
			if strings.HasSuffix(hose.G.Name,"\\t"){
				continue
			}
			ggsrbs[hose.G.Name] = gsrvs[hose.G.Name]
		}
	}
	gsrvb ,_:= json.Marshal(ggsrbs)
	fmt.Println(string(gsrvb))


	rsb,_ :=json.Marshal(hostsRules)
	fmt.Println(string(rsb))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func in( s string,ss ...string)bool{
	for _, val := range ss {
		if val == s{
			return true
		}
	}
	return false
}
// nginx_worker_processes = auto
// KONG_NGINX_WORKER_PROCESSES=16
// pg_max_concurrent_queries
type groupElem struct {
	Name string
	Rule []string
	RuleRaw []*rule
	Rules string
	Ips string

}

type hostelem struct {
	Host string
	Groups []*groupElem
	G *groupElem
	Len int
	Ips string
}
