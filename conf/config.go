package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/seeadoog/errors"
	"github.com/seeadoog/ngcfg"
	"github.com/xfyun/xns/tools/env"
	"io/ioutil"
	"sync/atomic"
)

type Config struct {
	AdminListen      string `json:"admin_listen"`                      // admin api 监听端口
	NameserverListen string `json:"nameserver_listen" required:"true"` // nameserver api 监听端口
	//DnsListenPort    int    `json:"dns_listen_port"`
	MetricsListen string `json:"metrics_listen" default:":5678"` // metrics 监听端口
	Login         struct {
		Enabled        bool `json:"enabled"`
		SessionTimeout int  `json:"session_timeout" default:"3600"`
	} `json:"login"`
	Database struct {
		Driver     string `json:"driver"`
		ConnectUrl string `json:"connect_url"`
	} `json:"database" required:"true"`

	ClusterEvents struct {
		PullInterval int `json:"pull_interval"` // 	集群事件拉去时间间隔，默认5s
	} `json:"cluster_events"`

	Resources struct {
		IpSrc string `json:"ip_src" required:"true"`
	} `json:"resources" required:"true"`
	Log map[string]*LogConf `json:"log" required:"true"` // 日志相关配置
}

var confInst atomic.Value

type LogConf struct {
	Level         string `json:"level" default:"error"`
	File          string `json:"file" default:"./log/server.log"`
	MaxSize       int    `json:"max_size" default:"200"`
	MaxBackup     int    `json:"max_backup" default:"30"`
	MaxAge        int    `json:"max_age"`
	Async         bool   `json:"async" default:"true"`
	CacheMaxCount int    `json:"cache_max_count"`
	BatchSize     int    `json:"batch_size"`
	Wash          int    `json:"wash"`
	Caller        bool   `json:"caller" default:"true"`
	CallerSkip    int    `json:"caller_skip" default:"2"`
}

// 解析配置文件
func parseCfgFromFile(file string) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.WithStack(err)
	}
	cfg := &Config{}
	if err := ngcfg.UnmarshalFromBytesCtx(bytes, cfg); err != nil {
		return err
	}
	confInst.Store(cfg)
	saveFileAsJson(cfg, file)
	return nil
}

func saveFileAsJson(f interface{}, name string) {
	b, _ := json.Marshal(f)
	bf := bytes.NewBuffer(nil)
	json.Indent(bf, b, "", "    ")
	err := ioutil.WriteFile(name+".json", bf.Bytes(), 0666)
	if err != nil {
		fmt.Println("save confing json error:", err)
	}
}

func Load() *Config {
	return confInst.Load().(*Config)
}

type Args struct {
	ConfigName string
}

func Init(a *Args) error {
	if a.ConfigName == "" || a.ConfigName == "env" {
		a.ConfigName = "./xns.cfg"
		return parseEnvConfig(a)
	}
	return parseFile(a)
}

func parseFile(a *Args) error {
	return parseCfgFromFile(a.ConfigName)
}

func parseEnvConfig(a *Args) error {
	e := env.NewEnv(`
nameserver_listen {{.server_listen}}   #nameserver 监听端口
admin_listen {{.admin_listen}}        # 管理api 监听端口
metrics_listen {{.metrics_listen}}      # prometheus 端口

login{
    enabled {{.login_enabled}}
    session_timeout {{.login_session_timeout}}
}
database { #数据库配置
    driver postgres
    connect_url "{{.db_connection_url}}"
}

resources{
    ip_src "{{.ip_src}}" # ip 地址资源
}

log{
    access{
        file {{.log_access}}
        level {{.log_access_level}}
    }

    admin{
       file {{.log_admin}}
        level {{.log_admin_level}}
    }

    runtime{
        file ./log/runtime.log
         level info
    }

    debug {
        file ./log/debug.log
         level info
    }

    error {
        file {{.log_error}}
        level {{.log_error_level}}
    }

    cluster {
        file ./log/cluster.log
         level info
    }
}


`)
	type Render = env.Render
	var (
		GetENVValue        = env.GetENVValue
		GetENVValueWithDef = env.GetENVValueWithDef
		WithDef            = env.WithDef
	)

	e.SetValue("server_listen", WithDef("server_listen", ":4567"))
	e.SetValue("admin_listen", WithDef("admin_listen", ":8806"))
	e.SetValue("metrics_listen", WithDef("metrics_listen", ":5678"))
	e.SetValue("login_enabled", WithDef("login_enabled", "false"))
	e.SetValue("login_session_timeout", WithDef("login_session_timeout", "720000"))
	e.SetValue("db_connection_url", Render(func() string {
		host := GetENVValue("pg_host")
		port := GetENVValue("pg_port")
		dbName := GetENVValue("pg_db_name")
		user := GetENVValue("pg_user")
		password := GetENVValue("pg_password")
		sslMode := GetENVValueWithDef("pg_sslmode", "disable")
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=%s password=%s", host, port, dbName, user, sslMode, password)
	}))
	e.SetValue("ip_src", WithDef("ip_src", "/root/ip.src"))
	e.SetValue("log_admin", WithDef("log_admin", "./log/admin.log"))
	e.SetValue("log_admin_level", WithDef("log_admin_level", "info"))

	e.SetValue("log_access", WithDef("log_access", "./log/access.log"))
	e.SetValue("log_access_level", WithDef("log_access_level", "info"))

	e.SetValue("log_error", WithDef("log_error", "./log/error.log"))
	e.SetValue("log_error_level", WithDef("log_error_level", "error"))

	file := bytes.NewBuffer(nil)
	err := e.Parse(file)
	if err != nil {
		return fmt.Errorf("parse config from env error:%w", err)
	}
	err = ioutil.WriteFile(a.ConfigName, file.Bytes(), 0666)
	if err != nil {
		return fmt.Errorf("create config file error:%w", err)
	}

	return parseFile(a)
}


