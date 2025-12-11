package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"lwc/zInx/utils"
	"lwc/zInx/ziface"
)

// 包含连接的属性的结构体
type Connection struct {
	TcpServer    ziface.IServer //当前conn属于哪个Server(也即本connection是由哪个server造出来的),在conn初始化的时候添加即可,如果想和连接管理器交互则先获得server字段即可
	Conn         *net.TCPConn
	ConnID       uint32 //连接id,本质上就是所谓的sessionID
	isClosed     bool
	handleAPI    ziface.HandFunc        //当前连接绑定的处理方法，可以注册一条handler执行链到一个连接上
	MsgHandler   ziface.IMsgHandle      //当前连接的处理方法
	ExitBuffChan chan bool              //信号通道,用于传达该连接已退出或停止的信号
	msgChan      chan []byte            //无缓冲通道(也即不设容量),用于读,写两个goroutine之间的消息通信
	msgBuffChan  chan []byte            //有缓冲通道,(也即设容量),也是用于读,写两个goroutine之间的消息通信
	property     map[string]interface{} //存放连接属性的map
	propertyLock sync.RWMutex           //保护上面这个map的锁
}

func NewConncetion(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte), //msgChan初始化
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}), //初始化属性map
	}
	c.TcpServer.GetConnMgr().Add(c) //连接和连接管理器需要交互的场景之1
	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		dp := NewDataPack()
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil { //一直阻塞直到请求头切片读满
			fmt.Println("read msg head error ", err)
			c.ExitBuffChan <- true
			continue
		}
		//拆包(读协议头和读实际数据分开解耦是必要的,解包器只负责读出解包所需字段,具体要怎么利用这些字段由业务层来说了算)
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			c.ExitBuffChan <- true
			continue
		}
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data) //这里没有和测试脚本一样用断言将接口转为实际的实现类不同，后续可看看是否有影响

		req := Request{
			conn: c,
			msg:  msg,
		}

		//这行代码的上面相当于是处理读事件的协程,数据读进来后,会有专门的新的处理业务逻辑的handle协程来做操作
		//相当于读操作和业务操作解耦,提高读写事件的io能力,这里的router其实就可以类比为netty框架里的handler的概念
		//这里的Connection对应的协程对应于netty里的响应读写事件的eventgroupworker线程
		//业务逻辑协程,相当于读,业务,写三件事情都解耦了,这种架构可以用最少的协程连接数来承担最大的通信并发量
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req) //else的话退化回没有协程池,无限创建业务处理协程的情况。
		}
	}
}

func (c *Connection) Start() {
	//阻塞等待客户端消息的读协程,读到消息后会交给注册的router handler协程去处理
	go c.StartReader()
	//router handler协程将业务处理完毕后,会把处理完毕的消息通过通道传递给下面的这个写协程
	go c.StartWriter()
	//连接创建完成后,先执行一遍之前注册到服务器上的钩子方法
	c.TcpServer.CallOnConnStart(c)
	//这样一来,读协程只负责循环接收消息，解包封装成request传递给router handler协程,handler协程处理完毕后会把
	//处理结果封包然后通过通道传递给
	for {
		select {
		case <-c.ExitBuffChan:
			//开启监听读,接着陷入循环阻塞,直到收到退出信号(后续可以添加一些资源回收的操作,如有必要的话哈)
			return
		}
	}
}

func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭连接前的亡语,其实就是回调函数
	c.TcpServer.CallOnConnStop(c)
	c.Conn.Close()
	c.ExitBuffChan <- true

	c.TcpServer.GetConnMgr().Remove(c) //连接需和连接管理器交互的场景之2

	close(c.ExitBuffChan)
	close(c.msgChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg.")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data)) //先构造msg,再把msg转为字节数组
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	c.msgChan <- msg //业务协程将消息通过无缓冲通道传递给写协程
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg.")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data)) //先构造msg,再把msg转为字节数组
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}
	c.msgBuffChan <- msg //业务协程将消息通过有缓冲通道传递给写协程
	return nil
}

// 写goroutine函数,阻塞等待读goroutine往管道里发消息,从通道里接收到的是handler处理过并且编码好了的消息
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send Data error:,", err, " Conn Writer exit")
				return
			}
		case data, ok := <-c.msgBuffChan: //有缓冲管道,仅在管道内无值且管道为非关闭状态下,会阻塞。
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error:,", err, " Conn Writer exit")
					return
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				break
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
