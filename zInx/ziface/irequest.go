package ziface

//把客户端请求定义为一个接口类型
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte //读进来的原始请求信息,也即一个字节数组
	GetMsgID() uint32
}
