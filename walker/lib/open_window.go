//go:build windows

package lib

import (
	"os"
)

func (o *Open) Open(dirname string) (*os.File, error) {
	return os.Open(dirname)
}
