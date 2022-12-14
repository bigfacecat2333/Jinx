package jnet

import (
	"Jinx/jinterface"
	"Jinx/utils"
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

	// 当前server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler jinterface.IMsgHandler

	// 该server的连接管理器
	ConnManager jinterface.IConnManager

	// 该server的连接创建时Hook函数 OnConnStart
	OnConnStart func(conn jinterface.IConnection)

	// 该server的连接断开时的Hook函数 OnConnStop
	OnConnStop func(conn jinterface.IConnection)
}

// NewServer 初始化Server模块的方法
func NewServer() jinterface.IServer {
	s := &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
	return s
}

func (s *Server) Start() {
	// 打印一下当前的server的一些信息
	fmt.Printf("[Jinx] Server Name: %s, listenner at IP: %s, Port: %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Jinx] Version: %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	// 启动一个线程去做服务端的监听业务，这样就不会阻塞主线程，希望在Server()中阻塞而不是Start()中阻塞
	go func() {
		// 0 开启消息队列及Worker工作池
		s.MsgHandler.StartWorkerPool()

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

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新连接
			if s.ConnManager.Len() > utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn = ", utils.GlobalObject.MaxConn)
				err := conn.Close()
				if err != nil {
					return
				}
				continue
			}

			// 将处理新连接的业务方法和conn进行绑定(封装成一个类，像一个协议一样)，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息 进行停止或者回收
	fmt.Println("[STOP] Jinx server name ", s.Name)
	s.ConnManager.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {}
}

// AddRouter 给当前的服务注册一个路由方法，供客户端的连接处理使用
func (s *Server) AddRouter(MsgId uint32, router jinterface.IRouter) {
	s.MsgHandler.AddRouter(MsgId, router)
	fmt.Println("Add Router Success!")
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() jinterface.IConnManager {
	return s.ConnManager
}

// SetOnConnStart 注册创建链接之后的钩子方法
func (s *Server) SetOnConnStart(hookFunc func(connection jinterface.IConnection)) {
	s.OnConnStart = hookFunc
}

// SetOnConnStop 注册销毁链接之前的钩子方法
func (s *Server) SetOnConnStop(hookFunc func(connection jinterface.IConnection)) {
	s.OnConnStop = hookFunc
}

// CallOnConnStart 调用创建链接之后的钩子方法
func (s *Server) CallOnConnStart(conn jinterface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("[Start] Call OnConnStart()")
		s.OnConnStart(conn)
	}
}

// CallOnConnStop 调用销毁链接之前的钩子方法
func (s *Server) CallOnConnStop(conn jinterface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("[Stop] Call OnConnStop()")
		s.OnConnStop(conn)
	}
}
