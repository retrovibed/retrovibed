package main

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/cmd/cmdglobalmain"
	"github.com/retrovibed/retrovibed/cmd/cmdmeta"
)

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
	return C.CString(authn.PublicKeyPath())
}

// json array of ip addresses
//
//export ips
func ips() *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	db, err := cmdmeta.Database(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	return C.CString(cmdglobalmain.Hostname())
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

func debug(envs ...string) {
	for _, e := range envs {
		_ = log.Output(2, fmt.Sprintln(e))
	}
}
