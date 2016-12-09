package syncmap

import (
	// "fmt"
	"math/rand"
	"sync"
	"time"
)

const DefaultDeepNumber int = 4096
const seed int = 131

// mapStrIter include map itself and sync RWMutex
type mapStrIter struct {
	content map[string]interface{}
	r       *sync.RWMutex
}

// mapIntIter include map itself and sync RWMutex
type mapIntIter struct {
	content map[int]interface{}
	r       *sync.RWMutex
}

// MapWithLock include map of length and mapIter of slice
type SyncMap struct {
	deep  int
	siter []*mapStrIter
	iiter []*mapIntIter
}

/*
* Item style
 */
type Element struct {
	key   interface{}
	value interface{}
}

/*
* key style
 */
type KeyEle struct {
	keystyle string
	value    interface{}
}

/*
* New object MaoWithLock and initilizte default  length of map and sync map
 */

func New() *SyncMap {
	return newSyncMap(DefaultDeepNumber)
}

/*
* New object MaoWithLock and initilizte deep parameter length of map and sync map
 */
func NewSyncMap(deep int) *SyncMap {
	return newSyncMap(deep)
}

/*
* new sync map
 */
func newSyncMap(deep int) *SyncMap {
	sm := &SyncMap{}
	sm.deep = deep
	sm.siter = make([]*mapStrIter, sm.deep)
	sm.iiter = make([]*mapIntIter, sm.deep)
	for k, _ := range sm.siter {
		sm.siter[k] = &mapStrIter{}
		sm.siter[k].content = make(map[string]interface{})
		sm.siter[k].r = new(sync.RWMutex)
	}
	for k, _ := range sm.iiter {
		sm.iiter[k] = &mapIntIter{}
		sm.iiter[k].content = make(map[int]interface{})
		sm.iiter[k].r = new(sync.RWMutex)
	}
	return sm
}

/*
* set key-value into map
* int style use remainder algorithm
* string string use rand algorithm
 */

func (sm *SyncMap) Set(key interface{}, value interface{}) {
	// var s interface{}
	switch key.(type) {
	case int:
		s := sm.RemainAddress(key.(int))
		s.r.Lock()
		s.content[key.(int)] = value
		s.r.Unlock()
	case string:
		s := sm.HashAddress(key.(string))
		s.r.Lock()
		s.content[key.(string)] = value
		s.r.Unlock()
	}
}

/*
* get value from key
 */

func (sm *SyncMap) Get(key interface{}) (value interface{}, ok bool) {
	switch key.(type) {
	case int:
		s := sm.RemainAddress(key.(int))
		s.r.RLock()
		value, ok = s.content[key.(int)]
		s.r.RUnlock()
	case string:
		s := sm.HashAddress(key.(string))
		s.r.RLock()
		value, ok = s.content[key.(string)]
		s.r.RUnlock()
	}
	return
}

/*
* range all items
 */
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
		for _, iiter := range sm.iiter {
			iiter.r.RLock()
			for k, v := range iiter.content {
				ch <- Element{key: k, value: v}
			}
			iiter.r.RUnlock()
		}
		close(ch)
	}()
	return ch
}

/*
* range all key
 */
func (sm *SyncMap) Rangekeys() <-chan KeyEle {
	ch := make(chan KeyEle)
	go func() {
		for _, siter := range sm.siter {
			siter.r.RLock()
			for k, _ := range siter.content {
				ch <- KeyEle{value: k, keystyle: "string"}
			}
			siter.r.RUnlock()
		}
		for _, iiter := range sm.iiter {
			iiter.r.RLock()
			for k, _ := range iiter.content {
				ch <- KeyEle{value: k, keystyle: "int"}
			}
			iiter.r.RUnlock()
		}
		close(ch)
	}()
	return ch
}

/*
* delete key - value
 */
func (sm *SyncMap) Delete(key interface{}) {
	switch key.(type) {
	case int:
		s := sm.RemainAddress(key.(int))
		s.r.Lock()
		delete(s.content, key.(int))
		s.r.Unlock()
	case string:
		s := sm.HashAddress(key.(string))
		s.r.Lock()
		delete(s.content, key.(string))
		s.r.Unlock()
	}
}

/*
* remainder location
 */

func (sm *SyncMap) RemainAddress(l int) *mapIntIter {
	return sm.iiter[remainder(sm.deep, l)&sm.deep]
}

func remainder(numberator int, denominator int) int {
	return denominator & (numberator - 1)
}

/*
* strign hash to int
 */
func strIntHash(key string) int {
	var h int
	for _, c := range key {
		h = h*seed + int(c)
	}
	return h
}

/*
* find a location with the given key
 */
func (sm *SyncMap) HashAddress(key string) *mapStrIter {
	return sm.siter[strIntHash(key)&sm.deep]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
