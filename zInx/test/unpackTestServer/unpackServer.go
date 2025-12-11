package main

import (
	"fmt"
	"io"
	"lwc/zInx/znet"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}
		go func(conn net.Conn) {
			dp := znet.NewDataPack()
			for {
				headData := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, headData)
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
					msg := msgHead.(*znet.Message) //类型断言,把接口转换为实现类并赋值给msg,因为实现类才有Data字段,接口只有方法
					msg.Data = make([]byte, msg.DataLen)
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}
					fmt.Println("==>Recv Msg:ID=", msg.Id, ", len=", msg.DataLen, ",data=", string(msg.Data))
				}

			}
		}(conn)
	}
}
