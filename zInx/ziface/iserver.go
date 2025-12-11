// 服务端的抽象类,定义一些标识服务器的抽象接口类
package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//开启业务服务
	Serve()
	//路由功能:给当前服务注册一个业务路由方法,共客户端连接时自定义业务处理逻辑时使用
	AddRouter(msgId uint32, router IRouter)
	//得到连接管理器
	GetConnMgr() IConnManager
	//服务端创建一个连接后需执行的钩子方法(Hook)的注册
	SetOnConnStart(func(IConnection))
	//服务端断开一个连接前需执行的钩子方法(Hook)的注册
	SetOnConnStop(func(IConnection))
	//连接创建后的钩子方法的调用(需先注册后调用)
	CallOnConnStart(conn IConnection)
	//连接断开前的钩子方法的调用(需先注册后调用)
	CallOnConnStop(conn IConnection)
}
