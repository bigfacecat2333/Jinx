package main

import "Jinx/Net"

/*
	基于Jinx框架来开发的服务器端应用程序
*/

func main() {
	// 1. 创建一个server句柄，使用Jinx的api
	s := Net.NewServer("Jinx V0.1")
	// 2. 初始化server的一些参数
	s.Serve()
}
