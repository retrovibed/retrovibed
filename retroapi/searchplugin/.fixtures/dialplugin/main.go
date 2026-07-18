// Command dialplugin is a search-plugin test fixture that reports whether it
// could open a TCP connection to --query, letting registry_test.go assert
// that the registry's firewall actually blocks non-public addresses.
package main

import (
	"encoding/json"
	"flag"
	"net"
	"os"
	"time"

	"github.com/retrovibed/retrovibed/retroapi/ddiscapi"

	// autohijack points net.DefaultResolver at wasinet's virtual sockets;
	// without it a wasip1 build has no network at all.
	_ "github.com/egdaemon/wasinet/wasinet/autohijack"
)

func main() {
	fs := flag.NewFlagSet("plugin", flag.ExitOnError)
	mimetype := fs.String("mimetype", "", "")
	query := fs.String("query", "", "")
	fs.Bool("adult", false, "")
	fs.Parse(os.Args[2:])

	status := "connected"
	if conn, err := net.DialTimeout("tcp", *query, 5*time.Second); err != nil {
		status = "blocked: " + err.Error()
	} else {
		conn.Close()
	}

	json.NewEncoder(os.Stdout).Encode(&ddiscapi.Import{
		Uri:      status,
		Mimetype: *mimetype,
	})
}
