# about

主要自用，在`go build`的时候，写入 `-ldflags -X` 信息


# 使用方法

## 0. get

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

## 2. ready exec

use `command` or `config.yml`

command

```text
# aboutbuilder --help
Usage of aboutbuilder:
      --app.description string     build bash file name (default "description of app")
      --app.name string            app name (default "demo")
      --app.version string         app version (default "0.0.0-test")
      --build.output string        output dir (default ".")
      --build.run                  run script generated
      --build.script string        bash file name (default "build.sh")
      --build.source string        source go file or dir (default "main.go")
      --copyright.start uint16     start year of copyright
      --copyright.update uint16    update year of copyright, no use now.
      --copyright.website string   your -website (default "www.example.com")
pflag: help requested
```

config.yml in the dir executing cmd 
```yaml
app:
  name: your_app
  version: 3.4.5
  description: the lopygo app heheheh
build:
  output: 
  script:
  source: ../cmd/demo/
  run: false
copyright:
  start: 2002
  update: 
  website: example.com
```

## 3. exec


```bash
aboutbuilder
```

now you can use cmd "ls -al" show your dir