package yiyun

import (
	"io"

	"golang.org/x/net/websocket"
)

type WebSocketEndpoint struct {
	ws     *websocket.Conn
	output chan []byte
	bin    bool
}

//NewWebSocketEndpoint 创建新的处理
func NewWebSocketEndpoint(ws *websocket.Conn, bin bool) *WebSocketEndpoint {
	return &WebSocketEndpoint{
		ws:     ws,
		output: make(chan []byte),
		bin:    bin,
	}
}

func (we *WebSocketEndpoint) Terminate() {
	Debug("websocket", "Terminated websocket connection")
}

func (we *WebSocketEndpoint) Output() chan []byte {
	return we.output
}

func (we *WebSocketEndpoint) Send(msg []byte) bool {
	var err error
	if we.bin {
		err = websocket.Message.Send(we.ws, msg)
	} else {
		err = websocket.Message.Send(we.ws, string(msg))
	}
	if err != nil {
		Error("websocket", "Cannot send: ", err)
		return false
	}
	return true
}

func (we *WebSocketEndpoint) StartReading() {
	if we.bin {
		go we.read_binary_frames()
	} else {
		go we.read_text_frames()
	}
}

func (we *WebSocketEndpoint) read_text_frames() {
	for {
		var msg string
		err := websocket.Message.Receive(we.ws, &msg)
		if err != nil {
			if err != io.EOF {
				Error("websocket", "Cannot receive: %s", err)
			}
			break
		}
		we.output <- append([]byte(msg), '\n')
	}
	close(we.output)
}

func (we *WebSocketEndpoint) read_binary_frames() {
	for {
		var msg []byte
		err := websocket.Message.Receive(we.ws, &msg)
		if err != nil {
			if err != io.EOF {
				Error("websocket", "Cannot receive: %s", err)
			}
			break
		}
		we.output <- msg
	}
	close(we.output)
}
