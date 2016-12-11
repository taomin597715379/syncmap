package syncmap

import (
	"math/rand"
	"sync"
	"time"
)

const DefaultDeepNumber int = 4096
const seed int = 131

// mapIter include map itself and sync RWMutex
type mapIter struct {
	content map[interface{}]interface{}
	r       *sync.RWMutex
}

// MapWithLock include map of length and mapIter of slice
type SyncMap struct {
	deep  int
	siter []*mapIter
}

// Item style
type Element struct {
	key   interface{}
	value interface{}
}

// New object MaoWithLock and initilizte default  length of map and sync map
func New() *SyncMap {
	return newSyncMap(DefaultDeepNumber)
}

// New object MaoWithLock and initilizte deep parameter length of map and sync map
func NewSyncMap(deep int) *SyncMap {
	return newSyncMap(deep)
}

// new sync map
func newSyncMap(deep int) *SyncMap {
	sm := &SyncMap{}
	sm.deep = deep
	sm.siter = make([]*mapIter, sm.deep)
	for k, _ := range sm.siter {
		sm.siter[k] = &mapIter{}
		sm.siter[k].content = make(map[interface{}]interface{})
		sm.siter[k].r = new(sync.RWMutex)
	}
	return sm
}

// set key-value into map
// int style use remainder algorithm
// string string use rand algorithm
func (sm *SyncMap) Set(key interface{}, value interface{}) {
	var index int
	switch key.(type) {
	case int:
		index = sm.RemainAddress(key.(int)) & sm.deep
	case string:
		index = sm.HashAddress(key.(string))
	}
	sm.siter[index].r.Lock()
	sm.siter[index].content[key] = value
	sm.siter[index].r.Unlock()
}

// get value from key
func (sm *SyncMap) Get(key interface{}) (value interface{}, ok bool) {
	var index int
	switch key.(type) {
	case int:
		index = sm.RemainAddress(key.(int)) & sm.deep
	case string:
		index = sm.HashAddress(key.(string))
	}
	sm.siter[index].r.RLock()
	value, ok = sm.siter[index].content[key]
	sm.siter[index].r.RUnlock()
	return
}

// range all items
func (sm *SyncMap) RangeItems() <-chan Element {
	ch := make(chan Element)
	go func() {
		for _, siter := range sm.siter {
			siter.r.RLock()
			for k, v := range siter.content {
				ch <- Element{key: k, value: v}
			}
			siter.r.RUnlock()
		}
		close(ch)
	}()
	return ch
}

//  delete key - value
func (sm *SyncMap) Delete(key interface{}) {
	var index int
	switch key.(type) {
	case int:
		index = sm.RemainAddress(key.(int)) & sm.deep
	case string:
		index = sm.HashAddress(key.(string))
	}
	sm.siter[index].r.Lock()
	delete(sm.siter[index].content, key)
	sm.siter[index].r.Unlock()
}

// size map
func (sm *SyncMap) Size() int {
	sumDeep := 0
	for _, siter := range sm.siter {
		siter.r.RLock()
		sumDeep += len(siter.content)
		siter.r.RUnlock()
	}
	return sumDeep
}

// remainder location
func (sm *SyncMap) RemainAddress(l int) int {
	return remainder(sm.deep, l)
}

func remainder(numberator int, denominator int) int {
	return denominator & (numberator - 1)
}

// strign hash to int
func strIntHash(key string) int {
	var h int
	for _, c := range key {
		h = h*seed + int(c)
	}
	return h
}

// find a location with the given key
func (sm *SyncMap) HashAddress(key string) int {
	var h int
	for _, c := range key {
		h = h*seed + int(c)
	}
	return h & sm.deep
}

// init function
func init() {
	rand.Seed(time.Now().UnixNano())
}
