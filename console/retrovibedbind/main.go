package main

import (
	"C"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/cmd/cmdglobalmain"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
)
import "os"

//export authn_bearer
func authn_bearer() *C.char {
	bearer, err := authn.NewBearer()
	if err != nil {
		log.Fatalln(err)
	}
	return C.CString(bearer)
}

//export public_key
func public_key() *C.char {
	encoded, err := os.ReadFile(authn.PublicKeyPath())
	if err != nil {
		log.Fatalln(err)
	}

	return C.CString(string(encoded))
}

//export ips
func ips() *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	db, err := cmdmeta.Database(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	results, err := cmdmeta.Hostnames(ctx, db)
	if err != nil {
		log.Fatalln(err)
	}

	encoded, err := json.Marshal(results)
	if err != nil {
		log.Fatalln(err)
	}

	return C.CString(string(encoded))
}

//export daemon
func daemon(jsonargs *C.char) {
	var args []string
	if err := json.Unmarshal([]byte(C.GoString(jsonargs)), &args); err != nil {
		log.Fatalln(err)
	}

	go cmdglobalmain.Main(args...)
}

func main() {}
