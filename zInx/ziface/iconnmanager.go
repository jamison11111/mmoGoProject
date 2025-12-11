/*
连接管理模块的抽象接口
*/
package ziface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection) //从管理模块中移除一个连接,并非关闭一个连接,被移除的连接只是不再接受ConnManager管理
	Get(connID uint32) (IConnection, error)
	Len() int   //获取管理的总连接的个数
	ClearConn() //删除并停止所管理的全部连接
}
