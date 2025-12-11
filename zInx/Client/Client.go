package main

import (
	"fmt"
	"io"
	"lwc/zInx/znet"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second) //延时3秒,先启动服务端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit!")
		return //连接失败,退出
	}
	fmt.Println("client start success!")
	for {
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(1, []byte("zInx V0.6 Client Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //阻塞直到切片读满
		if err != nil {
			fmt.Println("read head error")
			break
		}
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*znet.Message) //断言类型强转
			msg.Data = make([]byte, msg.GetDataLen())
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("==>Recv Msg:ID=", msg.Id, ",len=", msg.DataLen, ",data=", string(msg.GetData()))
		}
		time.Sleep(1 * time.Second)
	}
}
