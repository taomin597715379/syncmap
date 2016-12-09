package syncmap

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

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

func TestSet(t *testing.T) {
	s := New()
	s.Set(1, "1gfrebvfdsbf")
	s.Set("1", "1gfrebvfdsbf")
	s.Set(3, "avdsavdsav")
	s.Set("3", GetRandomString(10))
	for ele := range s.RangeItems() {
		fmt.Println(ele.key, ele.value)
	}
	for ele := range s.Rangekeys() {
		fmt.Println(ele.value, ele.keystyle)
	}
	s.Delete(1)
	fmt.Println(s.Get(1))
	s.Delete("1")
	fmt.Println(s.Get("1"))
	for ele := range s.Rangekeys() {
		fmt.Println(ele.value, ele.keystyle)
	}
}

func BenchmarkSyncMapIntKey(b *testing.B) {
	s := New()
	var wg sync.WaitGroup
	// 10个 goroutine 打数据
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

func BenchmarkSimpleMapIntKey(b *testing.B) {
	var s = make(map[int]interface{})
	l := new(sync.RWMutex)
	var wg sync.WaitGroup
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				l.Lock()
				s[rand.Intn(5000)] = GetRandomString(10)
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkSyncMapStringKey(b *testing.B) {
	s := New()
	var wg sync.WaitGroup
	// 10个 goroutine 打数据
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

func BenchmarkSimpleMapStringKey(b *testing.B) {
	s := New()
	var wg sync.WaitGroup
	// 10个 goroutine 打数据
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
