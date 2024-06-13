package pkg

import (
	"github.com/githgf/ikuai/action"
	"sync"
)

type VlanMap struct {
	sync.RWMutex
	Map map[string]action.VlanData
}

func (l *VlanMap) readMap(key string) (action.VlanData, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *VlanMap) writeMap(key string, value action.VlanData) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}
