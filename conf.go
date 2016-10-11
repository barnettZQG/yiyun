package yiyun

import (
	"os"

	"github.com/jingweno/conf"
)

var (
	APPPath      = ""
	Server       = "yiyun pro server"
	StaticPath   = "/assets"
	StaticSource = "/assets"
	avg          *conf.Conf
)

func init() {
	if APPPath == "" {
		var err error
		APPPath, err = os.Getwd()
		if err != nil {
			Panic("httpserver get current app path error", err)
		}
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
	return avg.Int(k)
}

//Merge 合并
func Merge(m map[string]interface{}) {
	avg.Merge(m)
}
