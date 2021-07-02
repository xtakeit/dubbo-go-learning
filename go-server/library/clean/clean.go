package clean

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	chains []func() error
	mu     sync.Mutex
)

// init 创建系统信号通道，并启动单独的协程阻塞监听退出信号
// 收到信号后调用Exit: 执行清理并以code 0退出进程
func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-c
		Exit()
	}()
}

// clean 以先入后出的顺序执行压入chains中的清理函数
// 当清理函数返回不为nil的error, 则将对应的error打印到标准错误输出中
func clean() {
	for l := len(chains); l > 0; l = len(chains) {
		if err := chains[l-1](); err != nil {
			_, _ = fmt.Fprint(os.Stderr, err)
		}
		chains = chains[:l-1]
	}
}
