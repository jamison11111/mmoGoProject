package znet

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func ClientTest() {
	fmt.Println("Client Test ... start")
	time.Sleep(3 * time.Second) //延时3秒,先启动服务端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err,exit!")
		return //连接失败,退出
	}
	fmt.Println("client start success!")
	for {
		_, err := conn.Write([]byte("hello ZINX lwc")) //向服务端发一条消息(以byte切片的方式发出去)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf) //阻塞式，监听读事件，把内容读进buf
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf(" server call back : %s,cnt= %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}

func TestServer(t *testing.T) {
	s := NewServer()
	go ClientTest()
	s.Serve()
}
