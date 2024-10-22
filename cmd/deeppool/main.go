package main

import (
	"log"
	"os"

	"github.com/james-lawrence/torrent"
)

func main() {
	src, err := os.Open("hello.world.txt")
	if err != nil {
		log.Fatalln(err)
	}
	md, err := torrent.NewFromReader(src)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("magnet uri:", torrent.NewMagnet(md).String())
}
