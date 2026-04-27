//go:build ios

package main

/*
#include <os/log.h>
#include <stdlib.h>

static void oslog_write(const char *msg) {
	os_log_with_type(OS_LOG_DEFAULT, OS_LOG_TYPE_DEFAULT, "%{public}s", msg);
}
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

type OSLogWriter struct{}

func (w *OSLogWriter) Write(p []byte) (n int, err error) {
	msg := C.CString(string(p))
	defer C.free(unsafe.Pointer(msg))
	C.oslog_write(msg)
	return len(p), nil
}

func redirectlogs() {
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))
	log.SetOutput(&OSLogWriter{})
}
