//go:build linux

package lib

import (
	"os"
	"syscall"
)

func (o *Open) OpenFileNoCache(dirname string) (*os.File, error) {
	return os.OpenFile(dirname, syscall.O_DIRECT, 0666)
}
