package yiyun

import (
	"encoding/json"
	"os"
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
)

//RequestHandler 主处理入口
func RequestHandler(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	p := Path(string(ctx.Path()))
	if ok := p.isStatic(); ok {
		StaticSourceDispatch(p, ctx)
	} else {
		ActionDispatch(ctx)
	}
	Info("Request url:", string(p), " duration time :", time.Now().Sub(start), "method:", string(ctx.Method()))
}

//ActionDispatch 控制器方法调度
func ActionDispatch(ctx *fasthttp.RequestCtx) {
	//执行拦截
	if RunIntercept(ctx) {
		return
	}
	path := string(ctx.Path())
	method := string(ctx.Method())
	router := GetRouter(path, method)
	if router == nil {
		ErrorReJSON(ctx, "请求路由不存在", 404)
		return
	}
	actionType := router.Action
	vc := reflect.New(actionType)
	action, ok := vc.Interface().(ActionInterface)
	if !ok {
		ErrorReJSON(ctx, "action is not ActionInterface", 404)
		return
	}
	if action == nil {
		ErrorReJSON(ctx, "请求路由不存在", 404)
		return
	}
	action.Init(ctx, router.MethodName)
	action.SetData(router.Data)
	action.Prepare()
	if ctx.IsDelete() {
		action.Delete()
	}
	if ctx.IsGet() {
		action.Get()
	}
	if ctx.IsPost() {
		action.Post()
	}
	if ctx.IsPut() {
		action.Put()
	}
	if fn, ok := reflect.TypeOf(action).MethodByName(router.MethodName); ok {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(action)
		fn.Func.Call(params)
		if action.IsJSON() {
			js := action.GetData()["json"]
			if js == nil {
				ErrorReJSON(ctx, "server response is empty", 500)
			} else {
				Success(ctx, js, action.GetCode())
			}
		} else {
			action.Postpare()
		}
	} else {
		ctx.Error("Router Action Error", 404)
	}

}

//StaticSourceDispatch 静态文件调度器
//只允许get操作
func StaticSourceDispatch(p Path, ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("Don't allow this method ,please use get", 403)
	} else {
		staticFile := p.GetStaticFile()
		fi, err := os.Stat(staticFile)
		if err != nil && os.IsExist(err) {
			ctx.NotFound()
		} else if fi != nil {
			if ok := fi.IsDir(); ok {
				ctx.NotFound()
			} else {
				fasthttp.ServeFile(ctx, staticFile)
			}
		} else {
			ctx.NotFound()
		}
	}
}

//ErrorReJSON 错误返回
func ErrorReJSON(ctx *fasthttp.RequestCtx, message string, code int) {
	ctx.Response.SetStatusCode(code)
	ctx.SetContentType("text/json; charset=utf8")
	responseJSON := &ResponseJSON{
		Code:    code,
		Message: message,
		Server:  Server,
	}
	content, err := json.Marshal(responseJSON)
	if err != nil {
		Error("marshal response json error.", err)
	}
	ctx.Response.SetBodyString(string(content))
}

//Success 成功返回
func Success(ctx *fasthttp.RequestCtx, data interface{}, code int) {
	ctx.Response.SetStatusCode(code)
	ctx.SetContentType("text/json; charset=utf8")
	content, err := json.Marshal(data)
	if err != nil {
		Error("marshal response json error.", err)
	}
	ctx.Response.SetBodyString(string(content))
}
