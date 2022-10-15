package jnet

import (
	"Jinx/jinterface"
	"errors"
	"fmt"
	"io"
	"net"
)

// Connection
/*
	连接模块
*/
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 当前链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出/停止的channel
	ExitChan chan bool

	// 该链接处理的方法Router
	Router jinterface.IRouter
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router jinterface.IRouter) jinterface.IConnection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Router:   router,
	}
	return c
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receive buf err ", err)
		//	// 下次仍然有可能读到数据
		//	continue
		//}

		// 创建拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 二进制流 8字节 包括dataLen和msgID
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// unpack,得到msgID和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}

		// 根据dataLen 再次读取data，放在msg.Data中
		if msg.GetMsgLen() > 0 {
			msg.SetData(make([]byte, msg.GetMsgLen()))

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(c.GetTCPConnection(), msg.GetData())
			if err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}

		// 得到当前conn数据的Request请求数据(就是把数据封装到一个Request中)
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 执行注册的路由方法
		go func(request jinterface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

		// 从路由中，找到注册绑定的Conn对应的router调用
	}
}

// Start 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 1. 启动从当前链接的读数据业务
	go c.StartReader()
	// TODO 启动从当前链接的写数据业务,通过router来处理
}

// Stop 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 关闭socket链接
	if err := c.Conn.Close(); err != nil {
		fmt.Println("close conn err = ", err)
		return
	}

	// 回收资源
	close(c.ExitChan)
}

// GetTCPConnection 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		fmt.Println("connection closed when send msg")
		return errors.New("connection closed when send msg")
	}

	// 将data进行封包 binaryMsg = MsgDataLen | MsgID | Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgID)
		return errors.New("pack error msg")
	}

	_, err = c.Conn.Write(binaryMsg)
	if err != nil {
		fmt.Println("send msg error ", err)
		return errors.New("send msg error")
	}
	return nil
}
