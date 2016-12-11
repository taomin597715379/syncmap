# syncmap

[![GoDoc](https://godoc.org/github.com/taomin597715379/syncmap?status.svg)](https://godoc.org/github.com/taomin597715379/syncmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/taomin597715379/syncmap)](https://goreportcard.com/report/github.com/taomin597715379/syncmap)
[![Licenses](https://img.shields.io/badge/license-bsd-orange.svg)](https://opensource.org/licenses/BSD-3-Clause)

## Introduction

Syncmap will be the traditional map and hash together, thus greatly through the map of the concurrent performance. According to incomplete testing found that this approach than the traditional way to improve the performance of twice. Users can not care about the lock problem, as in the absence of concurrent use of the same map.

## Usage

Install with:

```bash
go get github.com/taomin597715379/syncmap
```

Example:
```go
import(
	"fmt"
	"github.com/taomin597715379/syncmap"
)

func main() {
	s :=syncmap.New()
	s.Set(1,1)
	v,ok := s.Get(1)
	fmt.Println(v, ok) // 1, ok

	v, ok := s.Get("this key not exists now")
	fmt.Println(v, ok) // nil, false
}
```
## Todo

- Rich variety of API
- Classification statistics for different types of Key
- Optimize the hash function to further improve concurrency