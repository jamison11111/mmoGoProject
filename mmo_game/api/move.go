package api

import (
	"fmt"
	"lwc/mmo_game/core"
	"lwc/mmo_game/pb"
	"lwc/zInx/ziface"
	"lwc/zInx/znet"

	"google.golang.org/protobuf/proto"
)

// 玩家移动路由结构体
type MoveApi struct {
	znet.BaseRouter
}

func (*MoveApi) Handle(request ziface.IRequest) {
	//1.将客户端传来的位置移动消息解码（msgId=3）
	msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("客户端发送来的位置同步消息解码失败", err)
		request.GetConnection().Stop()
		return
	}
	//2.获取发送位置同步消息的客户端id
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("获取客户端id失败:", err)
		request.GetConnection().Stop()
		return
	}
	fmt.Println("玩家%d,移动到了(%f,%f,%f,%f)", pid, msg.X, msg.Y, msg.Z, msg.V)
	//3.根据pid得到player对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//4.让player对象发起移动位置信息广播
	player.UpdatePos(msg.X, msg.Y, msg.Z, msg.V)
}
