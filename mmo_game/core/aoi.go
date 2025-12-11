/*aoi管理模块,管理所有的grid格子*/
package core

import (
	"fmt"
)

// 一些常量的定义,同包内的其他go文件可以直接使用
const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y     = 20
)

type AOIManager struct {
	MinX  int //整个区域的左边界的坐标
	MaxX  int
	CntsX int //x方向格子的数量
	MinY  int //整个区域的下边界坐标
	MaxY  int
	CntsY int           //y方向格子的数量
	grids map[int]*Grid //key为格子id,value为格子对象,根据格子id和x轴上格子数量可以反推出格子的横纵坐标
}

// 给AOI管理器初始化所有的格子,相当于初始化一整个正方形游戏地图
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX, //x方向格子的数量
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY, //y方向各自的数量
		grids: make(map[int]*Grid),
	}
	//给AOI管理器初始化所有的格子,相当于初始化一整个正方形游戏地图
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			gid := y*cntsX + x //根据格子的横纵坐标和地图的格子数量来计算格子id
			//,根据横纵坐标序号可以算出每个Grid格子的坐标实际上下界
			aoiMgr.grids[gid] = NewGrid(gid, aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(), aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

// 得到每个各自在x轴方向的长度度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

// 得到每个格子在y轴方向的高度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印地图的所有相关信息,相当于java中的tostring方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\nminX:%d,maxX:%d,cntsX:%d,minY:%d,maxY:%d,cntsY:%d,\nGrids in AOI Manager:\n", m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}
	return s
}

// 根据格子的gID得到与其相邻的周边的九宫格的信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//先判断地图上有没有这个id的格子
	if _, ok := m.grids[gID]; !ok {
		return
	}
	//go中的返回值变量不用重复定义,可以直接使用,这里的返回值切片指针grids是用来存九宫格元素的,可以直接用
	grids = append(grids, m.grids[gID])
	idx := gID % m.CntsX //求取当前格子在x轴上的序号
	if idx > 0 {
		grids = append(grids, m.grids[gID-1]) //此时左边必有格子
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1]) //此时右边必有格子
	}
	//获取x轴上的所有格子的格子id并扔到一个整型切片里面
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	for _, v := range gidsX {
		idy := v / m.CntsX //y轴上的序号
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX]) //此时上方必有格子
		}
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX]) //此时下方必有格子
		}
	}
	return
}

// 根据横纵坐标计算横纵坐标序号,根据序号求出格子id
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()
	return gy*m.CntsX + gx
}

// 先求格子序号,再求九宫格,最后再求九宫格内所有的玩家,这些玩家是当前玩家视野可见的
func (m *AOIManager) GetPIDByPos(x, y float32) (playerIDs []int) {
	gID := m.GetGIDByPos(x, y)
	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		fmt.Printf("===>grid ID:%d,pids:%v===", v.GID, v.GetPlayerIDs())
	}
	return
}

// 获取当前格子id所代表的那个格子上的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 将一个PlayerID添加到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 通过横纵坐标添加一个player到某格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	m.grids[gID].Add(pID)
}

// 通过横纵坐标删除某个格子里的某个玩家id(下线)
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	m.grids[gID].Remove(pID)
}
