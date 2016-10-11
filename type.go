package yiyun

import "strings"

//ResponseJSON http返回数据结构
type ResponseJSON struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Server  string      `json:"server"`
	Data    interface{} `json:"data"`
}

//Path 访问路径
type Path string

func (p Path) isStatic() bool {
	if ok := strings.HasPrefix(string(p), StaticPath); ok {
		return true
	}
	return false
}

//GetStaticFile 根据路径查找实际路径
func (p Path) GetStaticFile() (filePath string) {
	filePath = APPPath + StaticSource + strings.Replace(string(p), StaticPath, "", 1)
	//log.Debug("Get static file source,path:", filePath)
	return
}
