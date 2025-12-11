package ziface

import "net"

//包含连接方法的接口
type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgId uint32, data []byte) error     //发包方法(编码+发包),将打包好的数据发到无缓冲的通信管道去给writer协程消费,发送时通道易阻塞,会阻塞我们的发送协程,也就是业务协程,单条连接上的消息并发量很大时体验很差
	SendBuffMsg(msgId uint32, data []byte) error //带缓冲的发送消息接口,对应带缓冲的消息通信管道,不会那么容易阻塞业务协程,体验更佳。
	//为连接绑定一些用户的自定义属性,这样业务协程就可以取出他们进行调用了
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

//定义一个函数类型,这个类型的名称是HandFunc,函数的形参和返回值也指定了
//这是个回调函数,用于服务端做实际的业务并Write回消息给客户端
//第二个参数是服务端读到的字节数组,第三个参数是这个字节数组内的有效内容的长度
type HandFunc func(*net.TCPConn, []byte, int) error
