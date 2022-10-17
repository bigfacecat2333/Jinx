package jinterface

/*
IServer 定义一个服务器接口
*/
type IServer interface {
	// Start 启动服务器
	Start()

	// Stop 停止服务器
	Stop()

	// Serve 运行服务器
	Serve()

	// AddRouter 路由功能：给当前服务注册一个路由方法，供客户端的连接处理使用
	AddRouter(MsgId uint32, router IRouter)

	// GetConnMgr 得到链接管理
	GetConnMgr() IConnManager

	// SetOnConnStart 注册创建链接之后的钩子方法
	SetOnConnStart(func(connection IConnection))

	// SetOnConnStop 注册销毁链接之前的钩子方法
	SetOnConnStop(func(connection IConnection))

	// CallOnConnStart 调用创建链接之后的钩子方法
	CallOnConnStart(connection IConnection)

	// CallOnConnStop 调用销毁链接之前的钩子方法
	CallOnConnStop(connection IConnection)
}
