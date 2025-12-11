module lwc/mmo_game

go 1.24.4

require lwc/zInx v0.0.0 //必须有个虚拟版本号

require (
	github.com/golang/protobuf v1.5.4 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

//本地依赖模块相互导入的方法,强行为本地模块指定本地路径映射
replace lwc/zInx => ../zInx/zinx
