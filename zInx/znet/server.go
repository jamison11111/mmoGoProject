// 服务器接口的实现类
package znet

import (
	"errors"
	"fmt"
	"lwc/zInx/utils"
	"lwc/zInx/ziface"
	"net"
)

type Server struct {
	Name        string                        //服务器名称
	IPVersion   string                        //通信协议版本,ipv4/ipv6?
	IP          string                        //绑定的ip地址
	Port        int                           //绑定的端口
	msgHandler  ziface.IMsgHandle             //绑定的路由
	ConnMgr     ziface.IConnManager           //当前Server的连接管理器
	OnConnStart func(conn ziface.IConnection) //该server连接创建时需调用的hook回调函数(比如,欢迎某某玩家上线!)
	OnConnStop  func(conn ziface.IConnection) //该server销毁连接时需调用的hook回调函数(比如,欢送某某玩家下线!)
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务提取为一个函数(该函数在服务端完成读操作后调用)
	fmt.Println("[Conn Handle] CallBackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToCLient error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP: %s, Port %d, is starting\n", s.IP, s.Port)
	fmt.Printf("服务器的名称,版本号，最大连接数，能够接收的最大的数据包的大小分别为 %s %s %d %d", s.Name, utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	//开启一个协程去做监听任务
	go func() {
		//解析监听地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		//真正的进行监听操作
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		//服务端端口开启监听成功后,在正式开始监听连接请求之前,先得由服务端句柄开启业务协程池服务
		s.msgHandler.StartWorkerPool()
		fmt.Println("start Zinx server", s.Name, "successful,now listenning")

		//TODO 每个连接都有一个id,此处应该有一个自动生成id的方法
		var cid uint32
		cid = 0

		for {
			conn, err := listener.AcceptTCP() //监听连接事件，acceptTCP是阻塞方法
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			//将这个连接封装为一个完整的数据结构，包含回调函数
			dealConn := NewConncetion(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server,name", s.Name)
	s.ConnMgr.ClearConn() //关闭所有连接资源
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve()如果在启动服S务是还需处理其他事情,可在这里添加
	select {} //阻塞住(本协程一退出,Start方法里开辟的协程也会一并退出)
}

// 暴露本服务器的构造方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("--->CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("--->CallOnConnStop....")
		s.OnConnStart(conn)
	}
}
