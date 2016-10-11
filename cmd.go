package yiyun

import (
	"flag"

	"github.com/valyala/fasthttp"
)

//Run 启动入口
func Run(addr string, compress bool) {
	flag.Parse()

	h := RequestHandler
	if compress {
		h = fasthttp.CompressHandler(h)
	}
	Info("ListenAndServer:127.0.0.1", addr)
	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		Error("Error in ListenAndServe", err)
	}
}
