package main

import (
	"encoding/json"
	"flag"
	"os"
)

type result struct {
	Magnet   string `json:"magnet"`
	Health   uint32 `json:"health"`
	Mimetype string `json:"mimetype"`
}

func main() {
	category := flag.String("category", "", "")
	query := flag.String("query", "", "")
	flag.Parse()

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(result{
		Magnet:   "magnet:?xt=urn:btih:1111111111111111111111111111111111111111&dn=" + *query,
		Health:   42,
		Mimetype: *category,
	})
}
