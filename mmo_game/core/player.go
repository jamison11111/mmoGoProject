package core

import (
	"fmt"
	"lwc/mmo_game/pb"
	"lwc/zInx/ziface"
	"math/rand"
	"sync"

	"google.golang.org/protobuf/proto"
)

type Player struct {
	Pid  int32
	Conn ziface.IConnection //当前玩家的连接
	X    float32
	Y    float32
	Z    float32
	V    float32
}

// 以下这两个全局变量(相当于java中的静态变量)用来作为每个玩家的id生成器
var PidGen int32 = 1
var IdLock sync.Mutex

func NewPlayer(conn ziface.IConnection) *Player {
	IdLock.Lock()
	id := PidGen
	PidGen++
	IdLock.Unlock()
	//初始位置相对来说还是比较随机的
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(160 + rand.Intn(10)),
		Y:    0,
		Z:    float32(134 + rand.Intn(17)),
		V:    0,
	}
	return p
}

/*
发送消息的函数,服务端通过player这个方法调用者来获得与这个player的socket连接,然后把要发给这个
player的消息通过这个连接发出去,其中这些消息先得通过proto协议序列化,序列化的数据传到客户端后,客
户端程序通过与之匹配的相同的一个proto协议将这些二进制消息给反序列化为C#的结构体,然后渲染到unity
引擎上产生相应的效果proto.Message是所有由proto协议编译出来的go文件消息的共同父类,这意味着里面
可以放任何proto编译出来的go消息结构体类型
*/
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	//序列化要发给客户端的消息
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	//从方法调用者处获取与客户端的socket连接
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	//调用Zinx框架的SendMsg发包
	/*
		重新捋一遍发包逻辑(也就是SendMsg方法在底层会做的事情),
		首先把消息类型和消息本身拼成一个tlv结构体(该结构体包含
		消息类型,消息长度和消息本身三个类型的字段)接着把他序列
		化为字节数组最后通过socekt的write方法写出去
	*/
	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}
}

/*
这是真正的协议方法,也是第一个协议方法,使用方为服务器,MsgID为1,
协议名称(消息名称或成go结构体名称)为SyncPid,用于将生成的玩家id发给客户端玩家
*/
func (p *Player) SyncPid() {
	//构造一条proto.Message的子类数据,然后直接发出去
	data := &pb.SyncPid{
		Pid: p.Pid,
	}
	p.SendMsg(1, data)
}

/*
由服务端主动发起,向其他玩家广播自己上线了的消息,主要是告知自己的出生位置坐标,
使用msgid=200的BroadCast类型消息进行广播。协议名称(消息名称或成go结构体名称)为BroadCast
*/
func (p *Player) BroadCastStartPosition() {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			&pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}
	p.SendMsg(200, msg)
}

// 由服务器来帮助player把它发的消息content给转发给所有登录到服务器上的在线玩家
func (p *Player) Talk(content string) {
	msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1, //代表聊天广播
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	players := WorldMgrObj.GetAllPlayers()
	for _, player := range players {
		//底层会把消息利用proto序列化,然后在tlv协议封装
		//客户端有配套的反序列化步骤和路由器,收到序号为200的消息后,客户端会自动
		//把它按照特定的规则解析渲染到前端聊天框上,从而实现世界聊天的效果
		player.SendMsg(200, msg)
	}
}
