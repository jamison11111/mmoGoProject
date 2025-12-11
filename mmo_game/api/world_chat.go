package api

import (
	"fmt"
	"lwc/mmo_game/core"
	"lwc/mmo_game/pb"
	"lwc/zInx/ziface"
	"lwc/zInx/znet"

	"google.golang.org/protobuf/proto"
)

// 客户端发来的msgid=2的聊天消息,由这个世界聊天api路由处理
type WorldChatApi struct {
	znet.BaseRouter
}

func (*WorldChatApi) Handle(request ziface.IRequest) {
	msg := &pb.Talk{}
	//客户端发来的是二进制数据,将其解码到talk结构体里
	err := proto.Unmarshal(request.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	//客户端会在发来的聊天请求中附加一个pid字段属性供服务端查,
	// 这个添加附加字段的操作是hook函数里写的,也即刚建立连接就得加上这个属性。
	//把pid和connection给绑定起来
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("GetProperty pid error ", err)
		request.GetConnection().Stop() //这个连接有问题,服务端直接主动断开之
		return
	}
	//把具体的玩家信息从全局哈希表中查出来
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.Talk(msg.Content)
}
