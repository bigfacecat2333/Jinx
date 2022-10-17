package main

import (
	"Jinx/jinterface"
	"Jinx/jnet"
	"fmt"
)

/*
	基于Jinx框架来开发的服务器端应用程序
*/

// PingRouter Ping test 自定义路由（用户）
type PingRouter struct {
	jnet.BaseRouter
}

// Handle Test Handle
func (r *PingRouter) Handle(request jinterface.IRequest) {
	fmt.Println("Call PingRouter Handle...")

	fmt.Println("receive from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

type HelloRouter struct {
	jnet.BaseRouter
}

func (r *HelloRouter) Handle(request jinterface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")

	fmt.Println("receive from client: msgID = ", request.GetMsgId(), ", data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello...Hello...Hello..."))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	// 1. 创建一个server句柄，使用Jinx的api
	s := jnet.NewServer()
	// 2. 给当前jinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 3. 启动server
	s.Serve()
}
