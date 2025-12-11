package znet

import (
	"fmt"
	"lwc/zInx/utils"
	"lwc/zInx/ziface"
	"strconv"
)

// 这个结构体对象应该是整个服务端全局一份并共享给所有的socket连接使用
type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //存放消息和处理器的映射关系的map
	WorkerPoolSize uint32                    //限定服务器上施行业务逻辑的协程的数量
	TaskQueue      []chan ziface.IRequest    //工作线程相当于消费者,每个消费者消费一条request通道
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId=", request.GetMsgID(), "is not FOUND!")
		return
	}
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api,msgId=" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId=", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQUeue chan ziface.IRequest) {
	fmt.Println("Worker ID= ", workerID, "is started")
	for {
		select {
		case request := <-taskQUeue:
			mh.DoMsgHandler(request)
		}
	}
}

// 这个方法所做的事情是,初始化所有的worker,为每个worker分配一条队列(也即channel通道)
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//先造队列通道
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//再启动工作协程,把通道挂上去
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

/*
request应该发给哪条队列呢？这里涉及到负载均衡,那花样可就多了,这里用最普通的轮询分配方法就行,
复杂的负载均衡算法不是这里的重点
这种算法下,每一条连接和特定的某个业务处理协程之间的关系是固定绑死的
*/
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("将id为", request.GetConnection().GetConnID(), "的连接发来的类型id为", request.GetMsgID(), "的消息分发给workerId为", workerID, "的协程对应的那条队列(channel通道)！")
	mh.TaskQueue[workerID] <- request
}
