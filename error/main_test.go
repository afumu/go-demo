package main

import (
	"fmt"
	"io"
	"testing"
)

func TestName(t *testing.T) {

}

// 这种方式就不太好，需要对err有太多的判断了
func WriteResp(w io.Writer, body io.Reader) error {
	_, err := fmt.Fprint(w, "aaa")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "bbb")
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(w, "ccc")
	if err != nil {
		return err
	}

	_, err = io.Copy(w, body)
	return err
}

type errWriter struct {
	io.Writer
	err error
}

func (ew errWriter) Write(buf []byte) (int, error) {
	if ew.err != nil {
		return 0, ew.err
	}
	n, err := ew.Writer.Write(buf)
	return n, err
}

// 这种方式比较好
func WriteResp2(w io.Writer, body io.Reader) error {
	ew := &errWriter{Writer: w}
	fmt.Fprint(ew, "aaa")

	fmt.Fprint(ew, "bbb")

	fmt.Fprint(ew, "ccc")

	io.Copy(ew, body)
	return ew.err
}
