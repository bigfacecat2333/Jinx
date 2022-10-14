package Net

import (
	"Jinx/JInterface"
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
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP: %s, Port: %d, is starting\n", s.IP, s.Port)
	go func() { // 启动一个线程去做服务端的监听业务，这样就不会阻塞主线程，希望在Server()中阻塞而不是Start()中阻塞
		// 1. 获取一个TCP的addr (创建一个套接字/句柄) 用于监听(localAddr) ,封装了包括bind(), inet_aton(), htons()等系统调用
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 2. 监听服务器的地址, listenner的作用是监听客户端的连接请求(是一个fd的列表)
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Jinx server succ, ", s.Name, "succ, now listening...")
		// 3. 阻塞(区别于io复用中的阻塞, 客户端返回才会消耗cpu)的等待客户端连接，处理客户端连接业务（读写）
		for {
			connect, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务，做一个最基本的最大512字节的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := connect.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}
					fmt.Printf("recv client buf %s, cnt %d\n", buf, cnt)
					// 断言
					if _, err := connect.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err", err)
						continue
					}
				}
			}()
		}
	}()

}

func (s *Server) Stop() {

}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态

}

// NewServer 初始化Server模块的方法
func NewServer(name string) JInterface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      0,
	}
	return s
}
