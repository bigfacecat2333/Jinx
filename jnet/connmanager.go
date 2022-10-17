package jnet

import (
	"Jinx/jinterface"
	"errors"
	"log"
	"sync"
)

type ConnManager struct {
	connections map[uint32]jinterface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                      // 保护连接集合的读写锁
}

// NewConnManager 初始化链接管理
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]jinterface.IConnection),
	}
}

// Add 添加链接
func (cm *ConnManager) Add(conn jinterface.IConnection) {
	// 保护共享资源Map 加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 将conn添加到ConnManager中
	cm.connections[conn.GetConnID()] = conn

	log.Println("connection add to ConnManager successfully: conn num = ", cm.Len())
}

// Remove 删除链接
func (cm *ConnManager) Remove(conn jinterface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除连接信息
	delete(cm.connections, conn.GetConnID())

	log.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", cm.Len())
}

// Get 根据ConnID获取链接
func (cm *ConnManager) Get(connID uint32) (jinterface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	// 根据connID从连接中获取链接
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

// Len 得到当前链接总数
func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

// ClearConn 清除并终止所有链接
func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn，并停止conn的工作
	for connID, conn := range cm.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connections, connID)
	}

	log.Println("Clear All Connections successfully: conn num = ", cm.Len())
}
