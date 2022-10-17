package jnet

import (
	"Jinx/jinterface"
	"Jinx/utils"
	"fmt"
	"strconv"
)

/*
	消息管理模块实现
*/

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]jinterface.IRouter

	// 负责Worker取任务的消息队列, 每个worker都有一个taskQueue, 一个taskQueue有多个request任务
	TaskQueue []chan jinterface.IRequest

	// 业务工作Worker池的worker数量 和 TaskQueue一一对应
	WorkerPoolSize uint32
}

// NewMsgHandler 初始化/创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]jinterface.IRouter),
		TaskQueue:      make([]chan jinterface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request jinterface.IRequest) {
	// 1 从request中找到msgID
	router, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), " is NOT FOUND! Need Register!")
		return
	}

	// 2 根据MsgID调度对应Router业务即可
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// AddRouter 添加一个路由业务方法，供客户端的链接处理使用
func (mh *MsgHandler) AddRouter(msgID uint32, router jinterface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id已经注册
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	// 2 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success!")
}

// StartWorkerPool 启动一个Worker工作池(开启工作流程)
func (mh *MsgHandler) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个Worker被启动
		// 当前的Worker对应的channel消息队列，开辟空间，第0个Worker用第0个channel
		// TaskQueue 保存的是能够处理的请求的最大数量(管道大小)
		mh.TaskQueue[i] = make(chan jinterface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan jinterface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue, 由Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request jinterface.IRequest) {
	// 1 将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来进行分配
	// 根据ConnID来进行取模，得到的值就是当前的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgId(),
		" to WorkerID = ", workerID)

	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
