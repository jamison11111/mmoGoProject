package main

import (
	"fmt"
	"lwc/zInx/ziface"
	"lwc/zInx/znet"
)

// 自定义一个路由并在主函数中传入
type PingRouter struct {
	znet.BaseRouter
}

// func (this *PingRouter) PreHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PreHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	fmt.Println("recv from client :msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//回写
	err := request.GetConnection().SendMsg(1, []byte("ping ping ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error:", err)
	}
}

// func (this *PingRouter) PostHandle(request ziface.IRequest) {
// 	fmt.Println("Call Router PostHandle")
// 	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ...\n"))
// 	if err != nil {
// 		fmt.Println("call back ping ping ping error")
// 	}
// }

// 自定义一个路由并在主函数中传入
type HelloZinxRouter struct {
	znet.BaseRouter
}

func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	fmt.Println("recv from client :msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//回写
	err := request.GetConnection().SendMsg(1, []byte("hello zinx v0.6\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error:", err)
	}
}

func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("DoConnectionBegin is Called ...")

	conn.SetProperty("Name", "林伟朝")
	conn.SetProperty("qq", "1195397685")
	fmt.Println("Set conn name and home properties done!")

	name, _ := conn.GetProperty("Name")
	msg := fmt.Sprintf("恭喜玩家%s成功与服务器建立连接并上线!", name)

	err := conn.SendMsg(2, []byte(msg)) //发给玩家本身,这里不是广播效果
	if err != nil {
		fmt.Println("SendMsg error", err)
	}
}

func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("玩家%d与服务器断开了连接!", conn.GetConnID()) //这个通知只有服务端知道,同样不是广播
}

func main() {
	s := znet.NewServer()
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	s.Serve()
}
