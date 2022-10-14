package jnet

import "Jinx/jinterface"

// BaseRouter 实现了router时，先嵌入这个BaseRouter基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct {
}

// PreHandle 和 PostHandle  方法都为空，可以根据需要重写

// PreHandle 在处理conn业务之前的钩子方法Hook
func (br *BaseRouter) PreHandle(request jinterface.IRequest) {

}

// Handle 在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request jinterface.IRequest) {

}

// PostHandle 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request jinterface.IRequest) {

}
