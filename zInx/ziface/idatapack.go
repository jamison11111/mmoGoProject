// 拆包封包的抽象接口
package ziface

type IDataPack interface {
	GetHeadLen() uint32                //获取包头长度
	Pack(msg IMessage) ([]byte, error) //封包方法
	Unpack([]byte) (IMessage, error)   //拆包
}
