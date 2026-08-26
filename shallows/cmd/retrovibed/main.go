package main

import (
	"os"

	"github.com/retrovibed/retrovibed/shallows/cmd/cmdglobalmain"
)

func main() {
	cmdglobalmain.Main(os.Args[1:]...)
}
