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
	Dispatch(method string, action ActionInterface)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	GetData() map[string]interface{}
	GetCode() int
	IsJSON() bool
	Postpare()
	SetData(map[string]string)
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
	Header     map[string]interface{}
	Session    Session
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
	a.parseHeader()
	a.Session = GetGlobalSessions().GetSession(&ctx.Response, &ctx.Request)
}

//CreateSession 创建session
func (a *Action) CreateSession() {
	a.Session = GetGlobalSessions().CreateSession(&a.Ctx.Response, &a.Ctx.Request)
}
func (a *Action) parseHeader() {
	a.Header = make(map[string]interface{}, 0)
	a.Ctx.Request.Header.VisitAll(func(key, value []byte) {
		a.Header[string(key)] = interface{}(value)
	})
}

func (a *Action) defaultData() {
	a.Data["web_title"] = Get("APPTITLE")
	a.Data["web_url"] = Get("APPURL")

}

//SetData 设置访问参数
func (a *Action) SetData(data map[string]string) {
	for k, v := range data {
		a.Data[k] = v
	}
}

//Dispatch 不使用反射调用控制层方法
func (a *Action) Dispatch(method string, action ActionInterface) {

}

//Prepare 控制器前执行
func (a *Action) Prepare() {

}

//Get get方法
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
	//解析模版
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

//GetFromCookie 从cookie中读取数据
func (a *Action) GetFromCookie(key string) []byte {
	return a.Request.Header.Cookie(key)
}

//CreateCookie 创建cookie
func (a *Action) CreateCookie() *fasthttp.Cookie {
	return fasthttp.AcquireCookie()
}

//SetCookie add cookie to response
func (a *Action) SetCookie(cookie *fasthttp.Cookie) {
	a.Response.Header.SetCookie(cookie)
}

//Redirect 设置客户端重定向
func (a *Action) Redirect(url string) {
	a.Code = 302
	a.Ctx.Response.Header.Add("Location", url)
}

type jsonData struct {
	Message string      `json:"mesage"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
}

//SetJSONData 设置json api返回的数据
func (a *Action) SetJSONData(data interface{}, message, status string) {
	a.Data["json"] = &jsonData{
		Data:    data,
		Message: message,
		Status:  status,
	}
	a.ServserJSON()
}
