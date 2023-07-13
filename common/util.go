package common

import "sync"

func SliceToChannel[T any](tokens []T) chan T {
	tokenChan := make(chan T)
	go func() {
		for _, token := range tokens {
			tokenChan <- token
		}
		close(tokenChan)
	}()
	return tokenChan
}

type ChainMap struct {
	m      sync.Map
	parent *ChainMap
}

func NewChainMap() *ChainMap {
	return &ChainMap{}
}

func (cm *ChainMap) GetValue(key interface{}) (value interface{}, ok bool) {
	value, ok = cm.m.Load(key)
	if ok {
		return value, ok
	}
	if cm.parent != nil {
		return cm.parent.GetValue(key)
	}
	return nil, false
}

func (cm *ChainMap) SetValue(key, value interface{}) {
	cm.m.Store(key, value)
}

func (cm *ChainMap) NewChild() *ChainMap {
	return &ChainMap{parent: cm}
}
