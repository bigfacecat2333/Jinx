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

// DoConnectionBegin 创建连接之后的Hook函数
func DoConnectionBegin(conn jinterface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
		return
	}

	// 设置链接属性
	fmt.Println("Set conn Name, 冰冷稻草人...")
	conn.SetProperty("Name", "冰冷稻草人")
	conn.SetProperty("Address", "https://github.com/bigfacecat2333")
	conn.SetProperty("Token", "1234567890")
}

// DoConnectionStop 销毁连接之前的Hook函数
func DoConnectionStop(conn jinterface.IConnection) {
	fmt.Println("DoConnectionStop is Called ...")
	fmt.Println("ConnID = ", conn.GetConnID(), " is lost...")

	// 获取链接属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if address, err := conn.GetProperty("Address"); err == nil {
		fmt.Println("Address = ", address)
	}
	if token, err := conn.GetProperty("Token"); err == nil {
		fmt.Println("Token = ", token)
	}
}

func main() {
	// 1. 创建一个server句柄，使用Jinx的api
	s := jnet.NewServer()
	// 2. 注册hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionStop)
	// 3. 给当前jinx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// 4. 启动server
	s.Serve()
}
