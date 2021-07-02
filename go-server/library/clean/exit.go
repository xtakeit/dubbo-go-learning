package clean

import (
	"fmt"
	"os"
)

// ExitErr 执行清理后向标准错误打印传入的err, 并以code 1退出进程
func ExitErr(err error) {
	mu.Lock()
	defer mu.Unlock()
	clean()
	_, _ = fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// Exit 执行清理并以code 0退出进程
func Exit() {
	mu.Lock()
	defer mu.Unlock()
	clean()
	os.Exit(0)
}
