package znet

import (
	"errors"
	"sync"
	"zinx/src/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 私有 connID -> conn
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	cm := &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
		connLock:    sync.RWMutex{},
	}
	return cm
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 共享资源map加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connections, conn.GetConnID()) //map的删除元素方式...
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.Unlock()
	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("no conn in connection map")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 停止并删除全部的连接信息, 注意用range遍历
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
}
