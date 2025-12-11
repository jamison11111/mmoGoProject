package znet

import (
	"errors"
	"fmt"
	"lwc/zInx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex //连接的读写锁,map做多线程修改时起到互斥保护作用
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully:conn num=", connMgr.Len())
}

// 仅议出管理,连接本身仍可正常运行
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), "successfully:conn num=", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	connMgr.connLock.RLock() //读读不互斥,读写互斥
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	for connID, conn := range connMgr.connections {
		conn.Stop()
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All Connections successfully:conn num=", connMgr.Len())
}
