package clean

import (
	"io"
)

// Push 向清理链中追加io.Closer实例,
func Push(closer io.Closer) {
	mu.Lock()
	defer mu.Unlock()

	chains = append(chains, closer.Close)
}

// PushFunc 向清理链中追加清理函数实例
func PushFunc(closeFunc func() error) {
	mu.Lock()
	defer mu.Unlock()

	chains = append(chains, closeFunc)
}
