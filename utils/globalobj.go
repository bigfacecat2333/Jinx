package utils

import (
	"Jinx/jinterface"
	"encoding/json"
	"os"
)

/*
存储一切有关jinx框架的全局参数，供其他模块使用
一些参数是可以通过jinx.json由用户进行配置
*/

type GlobalObj struct {
	// Server
	TcpServer jinterface.IServer // 当前jinx全局的Server对象
	Host      string             // 当前服务器主机监听的IP
	TcpPort   int                // 当前服务器主机监听的端口号
	Name      string             // 当前服务器的名称

	// Jinx
	Version          string // 当前jinx框架的版本号
	MaxConn          int    // 当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 // 当前框架数据包的最大值
	WorkerPoolSize   uint32 // 当前框架工作池的Goroutine数量
	MaxWorkerTaskLen uint32 // 每个worker对应的消息队列的任务的最大数量
}

// GlobalObject 定义一个全局的对外GlobalObj
var GlobalObject *GlobalObj

// Reload 重新加载用户自定义的参数
func (g *GlobalObj) Reload() {
	file, err := os.ReadFile("conf/jinx.json")
	if err != nil {
		return
	}

	// 将json文件数据解析到struct中
	err = json.Unmarshal(file, &GlobalObject)
	if err != nil {
		return
	}
}

// 提供一个init方法，初始化当前的GlobalObject变量
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:             "JinxServerApp",
		Version:          "V0.9",
		TcpPort:          8999,
		Host:             "0.0.0.0", // 监听所有网卡
		MaxConn:          12000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 应该尝试从conf/jinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}
