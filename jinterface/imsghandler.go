package jinterface

/*
	消息管理抽象接口
*/

type IMsgHandler interface {
	// DoMsgHandler 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	// AddRouter 添加一个路由业务方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
}
