package offheap

import (
	"fmt"

	"github.com/shurcooL/go-goon"
)

var Verbose bool

func VPrintf(format string, a ...interface{}) {
	if Verbose {
		fmt.Printf(format, a...)
	}
}

func VDump(i interface{}) {
	if Verbose {
		goon.Dump(i)
	}
}
