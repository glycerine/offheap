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

// p is a shortcut for a call to fmt.Printf that implicitly starts
// and ends its message with a newline.
func p(format string, stuff ...interface{}) {
	fmt.Printf("\n "+format+"\n", stuff...)
}
