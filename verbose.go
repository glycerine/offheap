package offheap

import (
	"fmt"

	"github.com/shurcooL/go-goon"
)

var verbose bool

func vprintf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func vdump(i interface{}) {
	if verbose {
		goon.Dump(i)
	}
}
