package yiyun

import (
	"os"

	"strings"

	"github.com/jingweno/conf"
)

var (
	APPPath      = ""
	Server       = "yiyun pro server"
	StaticPath   = "/assets"
	StaticSource = "/assets"
	avg          *conf.Conf
)

//SetStaticPathPre 设置静态文件path前缀
func SetStaticPathPre(pre string) {
	if strings.HasPrefix(pre, "/") {
		StaticPath = pre
	} else {
		StaticPath = "/" + pre
	}
}

//SetStaticSourcePre 设置静态文件路由路径
func SetStaticSourcePre(pre string) {
	if strings.HasPrefix(pre, "/") {
		StaticSource = pre
	} else {
		StaticSource = "/" + pre
	}
}

//SetServerName 设置server name
func SetServerName(server string) {
	Server = server
}
func init() {
	if APPPath == "" {
		var err error
		APPPath, err = os.Getwd()
		if err != nil {
			Panic("httpserver get current app path error", err)
		}
		//APPPath = "/Users/qingguo/gowork/src/yiyun-docker/"
	}
	load()
}

func load() {
	loader := conf.NewLoader()
	loader.Env().Argv()
	_, err := os.Stat(APPPath + "/conf/app.json")
	if err != nil && os.IsExist(err) {
		Panic("app.conf is not exit", err)
	}
	loader.File(APPPath + "/conf/app.json")
	c, err := loader.Load()
	if err != nil {
		Panic("load app.conf error.", err)
	}
	avg = c
	loadFirst()
}

func loadFirst() {
	if Get("port") == "" {
		Set("port", "8080")
	}
}

//Get 获取配置
func Get(k string) interface{} {
	return avg.Get(k)
}

//String 获取配置
func String(k string) string {
	return avg.String(k)
}

//Set 添加配置
func Set(k string, v interface{}) {
	avg.Set(k, v)

}

//Bool 获取
func Bool(k string) bool {
	return avg.Bool(k)
}

//Int 获取
func Int(k string) int {
	return int(avg.Get(k).(float64))
}

//Merge 合并
func Merge(m map[string]interface{}) {
	avg.Merge(m)
}
