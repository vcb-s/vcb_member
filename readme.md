#### 模块初始化顺序
`conf -> helper -> router -> service -> models -> [whatever.go] -> main`

每个文件只能依赖比自己先初始化的文件，不能依赖比自己晚初始化的文件，否则会出现依赖循环文件。

#### depend for vscode
1. 安装扩展依赖：
```bash
export GO111MODULE=on
export GOPROXY=https://goproxy.io
go get -u -v github.com/ramya-rao-a/go-outline
go get -u -v github.com/acroca/go-symbols
go get -u -v github.com/mdempsky/gocode
go get -u -v github.com/rogpeppe/godef
go get -u -v golang.org/x/tools/cmd/godoc
go get -u -v github.com/zmb3/gogetdoc
go get -u -v golang.org/x/lint/golint
go get -u -v github.com/fatih/gomodifytags
go get -u -v golang.org/x/tools/cmd/gorename
go get -u -v sourcegraph.com/sqs/goreturns
go get -u -v golang.org/x/tools/cmd/goimports
go get -u -v github.com/cweill/gotests/...
go get -u -v golang.org/x/tools/cmd/guru
go get -u -v github.com/josharian/impl
go get -u -v github.com/haya14busa/goplay/cmd/goplay
go get -u -v github.com/uudashr/gopkgs/cmd/gopkgs
go get -u -v github.com/davidrjenni/reftools/cmd/fillstruct
```
2. 安装 microsoft 发行的 go 扩展
3. 丢弃刚刚安装时产生的`go.mod`和`go.sum`文件修改
4. 安装项目依赖：`go get`

#### 依赖修正命令
`go mod tidy`