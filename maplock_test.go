package syncmap

import (
	"bytes"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	deep := 1024
	if s1 := New(); s1 == nil {
		t.Error("New() is nil")
	}
	if s2 := NewSyncMap(deep); s2 == nil {
		t.Error("NewSyncMap() is nil")
	}
}

func TestSet(t *testing.T) {
	s := New()
	s.Set(1, 1)
	s.Set("1", 1)
	if s.Size() != 2 {
		t.Error("Set() Number should equal 2")
	}
}

func TestGet(t *testing.T) {
	s := New()
	v, ok := s.Get("this key not exists now")
	if ok {
		t.Error("ok should be false")
	}
	if v != nil {
		t.Error("valus should be nil for missing now")
	}
	s.Set(1, 1)
	v, ok = s.Get(1)
	if !ok {
		t.Error("ok should be true")
	}
	if v != 1 {
		t.Error("v should be an integer 1")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	s.Set("hello", "world")
	s.Delete("hello")
	if _, ok := s.Get("hello"); ok {
		t.Error("delete fail")
	}
}

func TestSize(t *testing.T) {
	s := New()
	for i := 0; i < 42; i++ {
		s.Set(i, 42-i)
	}
	if s.Size() != 42 {
		t.Error("size return the wrong of number ")
	}
}

func TestRangeItems(t *testing.T) {
	var expectedIter = make([]interface{}, 0)
	s := New()
	for i := 0; i < 42; i++ {
		s.Set(i, 42-i)
	}
	for iter := range s.RangeItems() {
		expectedIter = append(expectedIter, iter)
	}
	if len(expectedIter) != 42 {
		t.Error("Rnage item return the wrong of number ")
	}
}

//生成随机字符串
func GetRandomString(l int64) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"
	b := []byte(str)
	buffer := bytes.NewBuffer(nil)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := int64(0); i < l; i++ {
		buffer.WriteString(string(b[r.Intn(len(b))]))
	}
	return buffer.String()
}

// 10个 goroutine 打数据
func BenchmarkSyncMapIntKey(b *testing.B) {
	s := New()
	var wg sync.WaitGroup
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				s.Set(rand.Intn(5000), GetRandomString(10))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// 10个 goroutine 打数据
func BenchmarkSyncMapStringKey(b *testing.B) {
	s := New()
	var wg sync.WaitGroup
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				s.Set(GetRandomString(10), GetRandomString(10))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
