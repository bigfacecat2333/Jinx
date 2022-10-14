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

// PreHandle Test PreHandle
func (r *PingRouter) PreHandle(request jinterface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		return
	}
}

// Handle Test Handle
func (r *PingRouter) Handle(request jinterface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		return
	}
}

// PostHandle Test PostHandle
func (r *PingRouter) PostHandle(request jinterface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		return
	}
}

func main() {
	// 1. 创建一个server句柄，使用Jinx的api
	s := jnet.NewServer("Jinx V0.2")
	// 2. 给当前jinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3. 启动server
	s.Serve()
}
