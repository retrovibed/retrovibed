package main

import (
	"os"

	"github.com/retrovibed/retrovibed/cmd/cmdglobalmain"
)

func main() {
	cmdglobalmain.Main(os.Args[1:]...)
}
