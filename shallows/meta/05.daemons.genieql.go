//go:build genieql.generate
// +build genieql.generate

package meta

import (
	genieql "github.com/james-lawrence/genieql/ginterp"
)

//easyjson:json
func Daemon(gql genieql.Structure) {
	gql.From(
		gql.Table("meta_daemons"),
	)
}
