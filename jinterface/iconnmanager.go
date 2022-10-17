package jinterface

/*
	IConnManager 连接管理抽象层
*/

type IConnManager interface {
	// Add 添加链接
	Add(conn IConnection)

	// Remove 删除链接
	Remove(conn IConnection)

	// Get 根据ConnID获取链接
	Get(connID uint32) (IConnection, error)

	// Len 得到当前链接总数
	Len() int

	// ClearConn 清除并终止所有链接
	ClearConn()
}
