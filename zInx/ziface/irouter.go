package ziface

/*
简而言之,路由可以让服务端开发者自定义业务处理逻辑,而且还可以加前置后置逻辑
*/
type IRouter interface {
	//三个钩子方法,分别对应服务器具体业务逻辑的前置钩子函数,具体业务逻辑函数,后置钩子函数
	PreHandle(request IRequest)
	Handle(request IRequest)
	PostHandle(request IRequest)
}
