package jnet

import (
	"Jinx/jinterface"
	"Jinx/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

// Connection
/*
	连接模块
*/
type Connection struct {
	// 当前Conn属于哪个Server
	TcpServer jinterface.IServer

	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 当前链接的ID
	ConnID uint32

	// 当前的链接状态
	isClosed bool

	// 告知当前链接已经退出/停止的channel
	ExitChan chan bool

	// 无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	// 消息管理模块 用来绑定MsgID和对应的处理业务API(router)关系
	MsgHandler jinterface.IMsgHandler

	// 链接属性集合
	property map[string]interface{}

	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

// NewConnection 初始化链接模块的方法
func NewConnection(TcpServer jinterface.IServer, conn *net.TCPConn, connID uint32, handler jinterface.IMsgHandler) jinterface.IConnection {
	c := &Connection{
		TcpServer:  TcpServer,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		MsgHandler: handler,
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// StartWriter 启动写数据的Goroutine,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println("[Writer is exit] connID = ", c.ConnID, "remote addr is ", c.RemoteAddr().String())

	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

// StartReader 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("[Reader is exit] connID = ", c.ConnID, "remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
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

		// 通过WorkerPool来处理
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// Start 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 1. 启动从当前链接的读数据业务
	go c.StartReader()
	// 2. 启动从当前链接的写数据业务,通过router来处理
	go c.StartWriter()

	// 按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应的hook函数
	c.TcpServer.CallOnConnStart(c)
}

// Stop 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID = ", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed {
		return
	}

	c.isClosed = true

	// 调用开发者注册的销毁链接之前的钩子函数
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket链接
	if err := c.Conn.Close(); err != nil {
		fmt.Println("close conn err = ", err)
		return
	}

	// 关闭Writer
	c.ExitChan <- true

	// 将当前链接从connMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

	// 将数据发送给管道
	c.msgChan <- binaryMsg
	return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	c.property[key] = value
	c.propertyLock.Unlock()
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	if value, ok := c.property[key]; ok {
		c.propertyLock.RUnlock()
		return value, nil
	} else {
		c.propertyLock.RUnlock()
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	delete(c.property, key)
	c.propertyLock.Unlock()
}
