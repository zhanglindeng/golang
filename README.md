# golang
go language code note

## library

### [uuid](https://github.com/landjur/golibrary.git)
### [micro](https://github.com/micro/micro.git)
### [hprose](https://github.com/hprose/hprose-go.git)
### [hprose](https://github.com/hprose)
### [gin-limiter](https://github.com/julianshen/gin-limiter)
### [Go 可视化性能分析工具](http://colobu.com/2017/03/02/a-short-survey-of-golang-pprof/)

### iota
```
type Level int

const (
	// 严重程度从低到高
	LevelDebug     Level = iota
	LevelInfo
	LevelNotice
	LevelWarning
	LevelError
	LevelCritical
	LevelAlert
	LevelEmergency
)
```

### [post json](http://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go)

### [echo session](https://github.com/ipfans/echo-session)

## test
- 测试单个文件，要带上源文件 `go test -v abc_test.go abc.go`
- 测试单个方法 `go test -v -test.run TestAbc`
- [参考网页](http://blog.csdn.net/shenlanzifa/article/details/51451814)
