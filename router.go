package yiyun

import (
	"reflect"
	"strings"
)

//Rule 路由规则 扩张string方法
type Rule string

//Isvalid 验证路由规则是否有效
// path: /member/:id/:m
func (p Rule) Isvalid() bool {
	if p == "" {
		return false
	}
	if p[0] != '/' {
		return false
	}
	for i, b := range p.String() {
		if !((b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '/' || b == ':' || (b >= '0' && b <= '9' && i > 0)) {
			return false
		}
	}
	return true
}

func (p Rule) String() string {
	return string(p)
}

//IsPass 路由是否符合此规则
//url: /member/1/helloword
func (p Rule) IsPass(url string, data map[string]string) bool {

	if p.String() == url {
		return true
	}
	ps := strings.Split(p.String(), "/")
	urls := strings.Split(url, "/")
	if len(ps) != len(urls) {
		return false
	}
	for i, psT := range ps {
		if psT == "" {
			if psT != urls[i] {
				return false
			}
			continue
		}
		if psT[0] != ':' {
			if psT != urls[i] {
				return false
			}
		}
		if psT[0] == ':' {
			if urls[i] == "" {
				return false
			}
			//路由数据封装
			data[psT] = urls[i]
		}
	}
	log.Debug("path:", p.String(), "url:", url)
	return true
}

//RouterInfo 路由信息
type RouterInfo struct {
	Rule       Rule
	Action     reflect.Type
	MethodName string
	MethodType string
	Data       map[string]string
}

var routers []*RouterInfo

func init() {
	routers = make([]*RouterInfo, 0)
}

//Router registry routers
func Router(p string, action ActionInterface, methodType, methodName string) {
	path := Rule(p)
	if !path.Isvalid() {
		log.Panic("router path:", path, "is not valid")
	}
	if action == nil {
		log.Panic("action is not nil")
	}
	ac := reflect.TypeOf(action)
	info := &RouterInfo{
		Rule:       path,
		Action:     ac,
		MethodName: methodName,
		MethodType: methodType,
		Data:       make(map[string]string),
	}
	routers = append(routers, info)
	Debug("Add router:", p)
}

//AnalyzeRouter 分析路由，判断准人规则，取出参数
func (r *RouterInfo) AnalyzeRouter(path string) bool {
	return r.Rule.IsPass(path, r.Data)
}

//GetRouter 获取路由
func GetRouter(path string, methodType string) (info *RouterInfo) {
	if routers == nil {
		return nil
	}
	for _, router := range routers {
		if router.AnalyzeRouter(path) {
			if strings.ToUpper(router.MethodType) == strings.ToUpper(methodType) {
				return router
			}
		}
	}
	return nil
}
