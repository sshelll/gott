# Gott

可视化指定固定单测执行的命令行工具

## Install

go version >= 1.17，git clone 当前 repo 并执行 `go build` 即可，或是直接使用下方命令：

```sh
go install github@SCU-SJL/gott/@latest
```

go version <= 1.16 可前往 Github Release 页面下载打包好的可执行文件(MacOS Only)



## Usage

在任意目录下直接使用 `gott` 来替换 `go test` 即可

例如:

`go test -v` => `gott -v`

`go test -v -race` => `gott -v -race`

`go test -gcflags=all=-l -race -coverprofile=coverage.out` => `gott -gcflags=all=-l -race -coverprofile=coverage.out`
