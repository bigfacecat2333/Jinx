package jinterface

/*
IRequest 接口:
将客户端请求的链接信息和请求的数据包装到一个Request中
*/
type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection

	// GetData 得到请求的消息数据
	GetData() []byte

	// GetMsgId 得到请求的消息ID
	GetMsgId() uint32

	// GetMsgLen 得到请求的消息长度
	GetMsgLen() uint32
}
