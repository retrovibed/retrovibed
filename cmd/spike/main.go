package main

import (
	"encoding/binary"
	"log"
	"math/rand/v2"
	"time"
)

func TimestampTransactionID() string {
	var b [binary.MaxVarintLen64]byte
	_ = binary.PutUvarint(b[:], rand.Uint64())
	v2 := b[:2]
	n := binary.PutVarint(b[:], time.Now().UnixNano())
	b[2] = v2[0]
	b[3] = v2[1]
	return string(b[:n])
}
func main() {
	v2 := rand.UintN(256)
	log.Println("DERP DERP DERP", []byte(TimestampTransactionID()))
	log.Println("DERP DERP DERP", v2, byte(v2>>8), byte(v2))
}
