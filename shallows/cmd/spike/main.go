package main

import (
	"encoding/binary"
	"log"
	"math"
	"math/rand/v2"
	"time"
)

func StorageProfit(users int, factor float64) float64 {
	return (100 * (factor/(math.Log2(float64(users+1))) + (1 - factor)))
}

func Derp(users int, factor float64) {
	costper := StorageProfit(users, factor)
	log.Println(users, ":", costper, costper*float64(1000))
}

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
	const factor = 0.98
	log.Println("DERP DERP DERP", []byte(TimestampTransactionID()))
	for i := 1; i < 100; i++ {
		Derp(i, factor)
	}

	Derp(999, factor)
	Derp(1000, factor)
	Derp(10000, factor)
	Derp(100000, factor)
	Derp(1000000, factor)
}
