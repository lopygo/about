# about

主要自用，在`go build`的时候，写入 `-ldflags -X` 信息


# 使用方法

## 0. build

```bash
go get -u github.com/lopygo/about/cmd/aboutbuilder
```

## 1. 写原代码

如 `cmd/demo` 所示
```go
package main

import (
	"fmt"

	"github.com/lopygo/about/info"
)

func main() {

	i, err := info.FromInput()
	if err != nil {
		fmt.Println(err)
		return
	}

	i.ShowAll()
}

```

