// 管理当前世界所有玩家的管理器
package core

import (
	"sync"
)

type WorldManager struct {
	AoiMgr  *AOIManager       //地图管理器
	Players map[int32]*Player //当前世界的所有玩家的集合哈希表
	pLock   sync.RWMutex
}

// 这个世界管理器应该是所有Player共享的模块,因此需定义一个该管理器的全局对象
// 相当于java里的静态变量,实际开发中一般用单例设计模式来保证变量的全局唯一
// 而不是直接用全局变量
var WorldMgrObj *WorldManager

// 当core包被引用并第一次使用的时候,core下所有go文件内的init函数会被自动执行
// 初始化世界管理器(间接地初始化了地图管理器,全局单例。)
func init() {
	WorldMgrObj = &WorldManager{
		Players: make(map[int32]*Player),
		AoiMgr:  NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
	}
}

// 增删查player的三个方法
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()
	//添加到具体的格子中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Z)
}

func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	wm.pLock.Lock()
	delete(wm.Players, pid) //TODO:似乎没有把玩家从网格中删除,没关系吗？
	wm.pLock.Unlock()
}

func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock() //TODO:似乎没有把玩家从网格中删除,没关系吗？
	return wm.Players[pid]
}

func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	players := make([]*Player, 0)
	for _, v := range wm.Players {
		players = append(players, v)
	}
	return players
}
