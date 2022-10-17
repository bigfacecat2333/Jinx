package jinterface

/*
	消息管理抽象接口
*/

type IMsgHandler interface {
	// DoMsgHandler 调度/执行对应的Router消息处理方法
	DoMsgHandler(request IRequest)

	// AddRouter 添加一个路由业务方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)

	// StartWorkerPool 启动一个Worker工作池(开启工作池的动作只能发生一次，一个Jinx框架只能有一个Worker工作池)
	StartWorkerPool()

	// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
	SendMsgToTaskQueue(request IRequest)
}
