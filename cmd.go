package yiyun

import "github.com/valyala/fasthttp"

//Run 启动入口
func Run(addr string) {
	h := RequestHandler
	Info("ListenAndServer:127.0.0.1", addr)
	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		Error("Error in ListenAndServe", err)
	}
}
