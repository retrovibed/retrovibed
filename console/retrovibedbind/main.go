package main

import "C"
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
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

//export authn_bearer_host
func authn_bearer_host(hostname *C.char) *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	// temporary
	ctransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defaultclient := &http.Client{Transport: ctransport, Timeout: 20 * time.Second}
	// defaultclient = httpx.BindRetryTransport(defaultclient, http.StatusTooManyRequests, http.StatusBadGateway, http.StatusInternalServerError, http.StatusRequestTimeout)

	bearer, err := authn.BearerForHost(ctx, defaultclient, C.GoString(hostname))
	if err != nil {
		log.Println(err)
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
