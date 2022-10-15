package jnet

import (
	"Jinx/jinterface"
	"fmt"
	"strconv"
)

/*
	消息管理模块实现
*/

type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32]jinterface.IRouter
}

// NewMsgHandler 初始化/创建MsgHandler方法
func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]jinterface.IRouter),
	}
}

// DoMsgHandler 调度/执行对应的Router消息处理方法
func (mh *MsgHandler) DoMsgHandler(request jinterface.IRequest) {
	// 1 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgId(), " is NOT FOUND! Need Register!")
		return
	}

	// 2 根据MsgID调度对应Router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
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
