package main

import (
	"fmt"
	"lwc/mmo_game/api"
	"lwc/mmo_game/core"
	"lwc/zInx/ziface"
	"lwc/zInx/znet"
)

// 当客户端与服务器建立连接时候的hook函数(服务端分配给玩家id需要做一次回显)
func OnConnectionAdd(conn ziface.IConnection) {
	//创建一个新玩家(内含玩家id的分配)
	player := core.NewPlayer(conn)
	//此时player的id属性已经填充完毕,直接调用它的同步方法向客户端同步新玩家id消息(底层对应一次socket写操作)
	player.SyncPid()
	/*向其他玩家广播自己上线了的消息,主要是告知自己的出生位置坐标,
	使用msgid=200的broadcast类型消息进行广播。*/
	player.BroadCastStartPosition()
	//将其添加到世界管理器中
	core.WorldMgrObj.AddPlayer(player)
	//将playerid绑定到conn的一个附加属性上
	conn.SetProperty("pid", player.Pid)
	//同步周边玩家上线信息,该函数会使自己出现在周边玩家的视野中（广播200），也会使自己视野可见的玩家出现在自己的视野里（单播消息202）
	player.SyncSurrounding()
	fmt.Println("=====>Player pidId= ", player.Pid, " arrive ===")

}

func main() {
	s := znet.NewServer()

	s.SetOnConnStart(OnConnectionAdd)

	//创建一个世界聊天路由器并将其注册到服务器上
	s.AddRouter(2, &api.WorldChatApi{})
	//注册位置信息广播路由，挂载在服务器上的路由其实就起到相当于时间监听器的作用
	s.AddRouter(3, &api.MoveApi{})
	s.Serve()
}
