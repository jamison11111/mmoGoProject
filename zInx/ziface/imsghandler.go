package ziface

//request业务消息处理的抽象接口,,定义了全局handler应该具有的一些行为函数
type IMsgHandle interface {
	DoMsgHandler(request IRequest)          //以非阻塞方式消费一条request消息
	AddRouter(msgId uint32, router IRouter) //为消息注册一个新的处理器
	StartWorkerPool()                       //启动worker工作池
	SendMsgToTaskQueue(request IRequest)    //将request消息发给消息队列,等待这条队列的消费者去处理它
}
