package main

import (
	"C"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"unsafe"

	"github.com/retrovibed/retrovibed/retroapi/authn"
	"github.com/retrovibed/retrovibed/retroapi/netmonx"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdglobalmain"
	"github.com/retrovibed/retrovibed/shallows/cmd/cmdmeta"
	"github.com/retrovibed/retrovibed/shallows/meta/identityssh"
	"golang.org/x/oauth2"
)
import "github.com/retrovibed/retrovibed/shallows/cmd/cmdopts"

//export build_version
func build_version() *C.char {
	version, err := cmdopts.BuildVersion()
	if err != nil {
		log.Println(err)
	}
	return C.CString(version)
}

//export oauth2_bearer
func oauth2_bearer() *C.char {
	var (
		err   error
		token *oauth2.Token
	)

	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	chttp := &http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       authn.TLSConfig(),
		},
	}

	signer, err := authn.SSHSigner()
	if err != nil {
		log.Println("failed to create oauth2 bearer token", err)
		return C.CString("")
	}

	if token, err = authn.Oauth2Bearer(ctx, signer, chttp, "", authn.UserDisplayName()); err != nil {
		log.Println("failed to create oauth2 bearer token", err)
		return C.CString("")
	}

	return C.CString(token.AccessToken)
}

//export authn_bearer
func authn_bearer() *C.char {
	bearer, err := authn.NewBearer(cmdopts.JWTSecret)
	if err != nil {
		log.Fatalln(err)
	}
	return C.CString(bearer)
}

//export authn_bearer_host
func authn_bearer_host(hostname *C.char) *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	shostname := C.GoString(hostname)
	ctransport := &http.Transport{
		TLSClientConfig: authn.TLSConfig(),
	}
	defaultclient := &http.Client{Transport: ctransport, Timeout: 20 * time.Second}
	defaultclient = authn.RetryClient(defaultclient)

	bearer, err := authn.BearerForHost(ctx, defaultclient, shostname)
	if err != nil {
		log.Println(err)
		return C.CString("")
	}

	return C.CString(bearer.AccessToken)
}

//export public_key
func public_key() *C.char {
	encoded, err := os.ReadFile(authn.PublicKeyPath())
	if os.IsNotExist(err) {
		return C.CString("")
	} else if err != nil {
		log.Fatalln(err)
	}

	return C.CString(string(encoded))
}

//export username
func username() *C.char {
	return C.CString(string(cmdglobalmain.Username()))
}

// returns an empty string on success, non empty contains the error.
//
//export unseed
func unseed() *C.char {
	if err := authn.Unseeded(); err != nil {
		return C.CString(err.Error())
	}

	return C.CString("")
}

// returns an empty string on success, non empty contains the error.
//
//export seed
func seed(s *C.char) *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	db, err := cmdopts.DatabaseMeta(ctx)
	if err != nil {
		log.Println("failed to connect to db", err)
		return C.CString(err.Error())
	}
	defer db.Close()

	id, err := authn.Seeded(ctx, C.GoString(s), false, authn.PrivateKeyPath())
	if err != nil {
		log.Println("failed to seed identity", err)
		return C.CString(err.Error())
	}

	if err = identityssh.InitializeAdmin(ctx, db, id.PublicKey()); err != nil {
		log.Println("unable to import ssh identity", err)
		return C.CString(err.Error())
	}

	return C.CString("")
}

//export ips
func ips() *C.char {
	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()
	db, err := cmdopts.DatabaseMeta(ctx)
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

//export netmon_metered
func netmon_metered() C.int {
	if netmonx.Metered() {
		return 1
	}
	return 0
}

//export netmon_set_metered
func netmon_set_metered(b C.int) C.int {
	_b := b != 0
	netmonx.SetMetered(_b)
	return b
}

//export gsetenv
func gsetenv(key *C.char, value *C.char) {
	os.Setenv(C.GoString(key), C.GoString(value))
}

//export egdaemon
func egdaemon(jsonargs *C.char) {
	var args []string
	if err := json.Unmarshal([]byte(C.GoString(jsonargs)), &args); err != nil {
		log.Fatalln(err)
		return
	}

	go cmdglobalmain.Main(args...)
}

//export logging
func logging() {
	redirectlogs()
}

//export checkpointdb
func checkpointdb() {
	// this method is to force checkpoint the database on system initialization.
	// this prevents a bunch of duckdb issues from impacting startup due to bad shutdowns.
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	db, err := cmdopts.DatabaseMeta(ctx)
	if err != nil {
		log.Println("failed to connect to database", err)
		return
	}
	defer db.Close()
}

//export validatecert
func validatecert(hostname *C.char, certData *C.uchar, certLen C.int) C.int {
	shostname := C.GoString(hostname)
	scert := C.GoBytes(unsafe.Pointer(certData), certLen)
	ctx, done := context.WithTimeout(context.Background(), time.Second)
	defer done()
	if err := authn.ValidateCertificate(ctx, shostname, scert); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func main() {}
