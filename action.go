package yiyun

import (
	"strings"

	"github.com/flosch/pongo2"
	"github.com/valyala/fasthttp"
)

//ActionInterface ActionInterface
type ActionInterface interface {
	//Init 初始化
	Init(ctx *fasthttp.RequestCtx, methodName string)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	GetData() map[string]interface{}
	GetCode() int
	IsJSON() bool
	Postpare()
}

//Action 控制器结构
type Action struct {
	// context data
	Ctx        *fasthttp.RequestCtx
	QueryArgs  *fasthttp.Args
	Request    *fasthttp.Request
	Response   *fasthttp.Response
	Data       map[string]interface{}
	Code       int
	isHTML     bool
	isJSON     bool
	Tpl        string
	actionName string
	methodName string
	Args       *fasthttp.Args
}

//Init 初始化
func (a *Action) Init(ctx *fasthttp.RequestCtx, methodName string) {
	a.Ctx = ctx
	a.Request = &ctx.Request
	a.Response = &ctx.Response
	a.QueryArgs = ctx.QueryArgs()
	a.Code = 200
	a.methodName = methodName
	a.Data = make(map[string]interface{}, 0)
	a.defaultData()
}

func (a *Action) defaultData() {
	a.Data["web_title"] = Get("APPTITLE")
	a.Data["web_url"] = Get("APPURL")
}
func (a *Action) Prepare() {

}
func (a *Action) Get() {

}
func (a *Action) Post() {

}
func (a *Action) Delete() {

}
func (a *Action) Put() {

}
func (a *Action) GetData() map[string]interface{} {
	return a.Data
}
func (a *Action) GetCode() int {
	return a.Code
}

//WriteString 返回字符串
func (a *Action) WriteString(data string) {

}

//WriteFile 返回文件
func (a *Action) WriteFile(file string) {

}

//Write 写回
func (a *Action) Write(buffer []byte) {
	a.Ctx.Response.SetBody(buffer)
}

//ServserJSON 以json形式返回
func (a *Action) ServserJSON() {
	a.isJSON = true
}

//ServerHTML 以解析模版返回
func (a *Action) ServerHTML() {
	a.isHTML = true
}

//IsJSON 是否以json返回
func (a *Action) IsJSON() bool {
	return a.isJSON
}

//Postpare Action执行之后执行，解析模版
func (a *Action) Postpare() {
	if a.Tpl == "" {
		a.Tpl = strings.ToUpper(a.methodName)
	}
	tplFile := APPPath + "/view/" + a.Tpl + ".html"
	if !FileIsExit(tplFile) {
		a.Ctx.Error("server error", 500)
		Error("template file is not exit with filename:", a.Tpl)
	} else {
		tp, err := pongo2.FromFile(tplFile)
		if err != nil {
			a.Ctx.Error("server error", 500)
			Error("parse template error.", err)
		} else {
			var tpl = pongo2.Must(tp, err)
			buffer, err := tpl.ExecuteBytes(a.Data)
			if err != nil {
				Error("Parse the tmplate file error.", err)
			} else {
				a.Write(buffer)
			}
		}
		a.Ctx.SetContentType("text/html;charset=utf8")
	}

}
