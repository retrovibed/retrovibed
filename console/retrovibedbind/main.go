package main

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/retrovibed/retrovibed/authn"
	"github.com/retrovibed/retrovibed/cmd/cmdglobalmain"
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
	log.Println(os.Environ())
	return C.CString(authn.PublicKeyPath())
}

//export ips
func ips() *C.char {
	return C.CString(cmdglobalmain.Hostname())
}

//export daemon
func daemon(jsonargs *C.char) {
	log.Println("DERP DERP", cmdglobalmain.Hostname())
	var args []string
	if err := json.Unmarshal([]byte(C.GoString(jsonargs)), &args); err != nil {
		log.Fatalln(err)
	}
	log.Println("DERP DERP", args)
	go cmdglobalmain.Main(args...)
}

func main() {}

func debug(envs ...string) {
	for _, e := range envs {
		_ = log.Output(2, fmt.Sprintln(e))
	}
}
