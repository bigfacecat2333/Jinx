package jnet

import "Jinx/jinterface"

type Request struct {
	// 已经和客户端建立好的 链接
	conn jinterface.IConnection

	// 客户端请求的数据
	msg jinterface.IMessage
}

func (r *Request) GetConnection() jinterface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetMsgLen() uint32 {
	return r.msg.GetMsgLen()
}
