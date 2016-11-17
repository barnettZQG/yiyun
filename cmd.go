package yiyun

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/websocket"
)

//Run 启动入口
func Run(addr string) {
	h := RequestHandler
	Info("ListenAndServer:127.0.0.1", addr)
	if err := fasthttp.ListenAndServe(addr, h); err != nil {
		Error("Error in ListenAndServe", err)
	}
}

//RunWebSocket 启动websocket服务监听
func RunWebSocket(addr string, h websocket.Handler) {
	Info("WebSocket ListenAndServer:127.0.0.1", addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		Error("Error in ListenAndServe", err)
	}
}
