package offheap

import (
	"fmt"

	"github.com/shurcooL/go-goon"
)

var verbose bool

func VPrintf(format string, a ...interface{}) {
	if verbose {
		fmt.Printf(format, a...)
	}
}

func VDump(i interface{}) {
	if verbose {
		goon.Dump(i)
	}
}
