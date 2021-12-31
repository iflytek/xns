package env

import (
	"fmt"
	"os"
	"testing"
)

/*

nameserver_listen :4567   #nameserver 监听端口
admin_listen :8806        # 管理api 监听端口
metrics_listen :5678      # prometheus 端口
dns_listen_port 53
login{
    enabled false
    session_timeout 30
}
database { #数据库配置
    driver postgres
    connect_url "host=10.1.87.70 port=55432 dbname=nameserverdev3 user=kong sslmode=disable"
}

resources{
    ip_src /Users/sjliu/temp/ip   # ip 地址资源
}

# 日志配置
log{
    access{
        file ./log/access.log
        level info
    }

    admin{
        file ./log/admin.log
        level info
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
        file ./log/error.log
         level info
    }

    cluster {
        file ./log/cluster.log
         level info
    }
}

*/
func TestNewEnv(t *testing.T) {
	e := NewEnv(`
nameserver_listen {{.nameserver_listen}}   #nameserver 监听端口
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
	os.Setenv("NS_PG_HOST","10.1.87.70")
	os.Setenv("NS_LOG_ERROR_LEVEL","debug")
	e.SetValue("nameserver_listen", WithDef("nameserver_listen", ":4567"))
	e.SetValue("admin_listen", WithDef("admin_listen", ":8806"))
	e.SetValue("metrics_listen", WithDef("metrics_listen", ":5678"))
	e.SetValue("login_enabled", WithDef("login_enabled", "true"))
	e.SetValue("login_session_timeout", WithDef("login_session_timeout", "720000"))
	e.SetValue("db_connection_url", Render(func() string {
		host := GetENVValue("pg_host")
		port := GetENVValue("pg_port")
		dbName := GetENVValue("pg_db_name")
		user := GetENVValue("pg_user")
		password := GetENVValue("pg_password")
		sslMode := GetENVValueWithDef("pg_sslmode","false")
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s sslmode=%s password=%s",host,port,dbName,user,sslMode,password)
	}))
	e.SetValue("ip_src",WithDef("ip_src","/root/ip.src"))
	e.SetValue("log_admin",WithDef("log_admin","./log/admin.log"))
	e.SetValue("log_admin_level",WithDef("log_admin_level","info"))

	e.SetValue("log_access",WithDef("log_access","./log/access.log"))
	e.SetValue("log_access_level",WithDef("log_access_level","info"))

	e.SetValue("log_error",WithDef("log_error","./log/error.log"))
	e.SetValue("log_error_level",WithDef("log_error_level","error"))

	err := e.Parse(os.Stdout)
	if err != nil{
		panic(err)
	}
}
