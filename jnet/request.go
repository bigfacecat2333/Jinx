package jnet

import "Jinx/jinterface"

type Request struct {
	// 已经和客户端建立好的 链接
	conn jinterface.IConnection

	// 客户端请求的数据
	data []byte
}

func (r *Request) GetConnection() jinterface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
