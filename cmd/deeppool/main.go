package main

import (
	"log"
	"os"

	"github.com/james-lawrence/torrent"
	"github.com/james-lawrence/torrent/metainfo"
)

func main() {
	src, err := os.Open("hello.world.txt")
	if err != nil {
		log.Fatalln(err)
	}
	minfo, err := metainfo.NewFromReader(src, metainfo.OptionDisplayName("hello.world.txt"))
	if err != nil {
		log.Fatalln(err)
	}
	md1, err := torrent.NewFromInfo(*minfo)
	// md1, err := torrent.NewFromReader(src, torrent.OptionDisplayName("hello.world.txt"))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("magnet uri:", torrent.NewMagnet(md1).String())
	md2, err := torrent.NewFromFile("hello.world.txt")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("magnet uri:", torrent.NewMagnet(md2).String())
}
