//go:build ios

package main

/*
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>

void nslog_write(const char *msg) {
    NSLog(@"%s", msg);
}
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

type NSLogWriter struct{}

func (w *NSLogWriter) Write(p []byte) (n int, err error) {
	msg := C.CString(string(p))
	defer C.free(unsafe.Pointer(msg))
	C.nslog_write(msg)
	return len(p), nil
}

func redirectlogs() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))
	log.SetOutput(&NSLogWriter{})
}
