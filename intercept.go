package yiyun

import (
	"strings"

	"github.com/valyala/fasthttp"
)

type Intercept interface {
	Intercept(*fasthttp.Request, *fasthttp.Response) bool
}

func init() {
	intercepts = make(map[string][]Intercept, 0)
}

var intercepts map[string][]Intercept

//AddIntercept 应用拦截器
func AddIntercept(path string, in Intercept) {
	intercepts[path] = append(intercepts[path], in)
}

//RunIntercept 执行拦截
func RunIntercept(ctx *fasthttp.RequestCtx) bool {
	path := string(ctx.Path())
	ins := getIntercepts(path)
	if ins == nil {
		return false
	}
	for _, in := range ins {
		if !in.Intercept(&ctx.Request, &ctx.Response) {
			return true
		}
	}
	return false
}

func getIntercepts(path string) []Intercept {
	for k, v := range intercepts {
		if isPass(path, k) {
			return v
		}
	}
	return nil
}

//source:来源路径 path:配置路径
func isPass(source string, path string) bool {

	if path == source {
		return true
	}
	ps := strings.Split(path, "/")
	urls := strings.Split(source, "/")
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
		}
	}
	return true
}
