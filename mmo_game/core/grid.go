package core

import (
	"fmt"
	"sync"
)

/*格子类*/
type Grid struct {
	GID       int //格子id
	MinX      int //实际的横坐标的下界,而不是坐标序号
	MaxX      int
	MinY      int
	MaxY      int
	playerIDs map[int]bool //本格子玩家或npc的集合
	pIDLock   sync.RWMutex
}

func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX, //实际的横坐标的下界,而不是坐标序号
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 将一个玩家或物品添加到当前格子中的方法
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	g.playerIDs[playerID] = true
}

func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()
	delete(g.playerIDs, playerID)
}

func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}
	return
}

// 打印当前各自的信息的日志方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d,minX:%d,maxX:%d,minY:%d,maxY:%d,playerIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
