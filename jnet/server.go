package jnet

import (
	"Jinx/jinterface"
	"fmt"
	"net"
)

// Server 定义一个server的服务器模块，实例化IServer
type Server struct {
	// 服务器的名称
	Name string

	// 服务器绑定的IP版本
	IPVersion string

	// 服务器监听的IP
	IP string

	// 服务器监听的端口
	Port int

	// 当前的server添加一个router，server注册的链接对应的处理业务
	Router jinterface.IRouter
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)

	// 启动一个线程去做服务端的监听业务，这样就不会阻塞主线程，希望在Server()中阻塞而不是Start()中阻塞
	go func() {

		// 1. 获取一个TCP的addr (创建一个套接字/句柄) 用于监听(localAddr) ,封装了包括bind(), inet_aton(), htons()等系统调用
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 2. 监听服务器的地址, listener的作用是监听客户端的连接请求(是一个socket_fd的列表)
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Jinx server success, ", s.Name, "success, now listening...")

		var cid uint32
		cid = 0

		// 3. 阻塞(区别于io复用中的阻塞, 客户端返回才会消耗cpu)的等待客户端连接，处理客户端连接业务（读写）
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定(封装成一个类，像一个协议一样)，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

// AddRouter 给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(router jinterface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Success!")
}

// NewServer 初始化Server模块的方法
func NewServer(name string) jinterface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
