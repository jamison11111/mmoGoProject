// 消息类型的接口类
package ziface

type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
