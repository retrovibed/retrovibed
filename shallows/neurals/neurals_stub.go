//go:build !(retrovibed && neural)

package neurals

import (
	"log"
	"sync"
)

var _stubOnce sync.Once

func predict(_ *Text, input string) (string, error) {
	_stubOnce.Do(func() { log.Println("warning neural stub running") })
	return input, nil
}
