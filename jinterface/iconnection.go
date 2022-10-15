package jinterface

import "net"

/*
IConnection 连接接口
*/
type IConnection interface {
	// Start 启动连接，让当前的连接开始工作
	Start()

	// Stop 停止连接，结束当前连接工作
	Stop()

	// GetTCPConnection 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接模块的连接ID
	GetConnID() uint32

	// RemoteAddr 获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr

	// SendMsg Send 发送数据，将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error
}
