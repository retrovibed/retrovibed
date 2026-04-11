//go:build android

package main

/*
#cgo LDFLAGS: -llog

#include <android/log.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"os"
	"unsafe"
)

// LogcatWriter implements io.Writer and sends data to Android Logcat
type LogcatWriter struct {
	tag      *C.char
	priority C.int
}

func (w *LogcatWriter) Write(p []byte) (n int, err error) {
	msg := C.CString(string(p))
	defer C.free(unsafe.Pointer(msg))
	C.__android_log_write(w.priority, w.tag, msg)
	return len(p), nil
}

// redirectlogs sends Go log output synchronously to Android Logcat.
// Using LogcatWriter directly avoids a goroutine race where log.Fatalln()
// calls os.Exit(1) before a pipe-reading goroutine can flush to logcat.
func redirectlogs() {
	cTag := C.CString("retrovibed")
	// Note: We don't free cTag because the writer uses it for the lifetime of the process
	log.SetFlags(log.Lshortfile | log.LUTC | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%d ", os.Getpid()))
	log.SetOutput(&LogcatWriter{tag: cTag, priority: C.ANDROID_LOG_INFO})
}
