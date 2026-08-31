package retrovibed

//go:generate genieql duckdb --database=dpool.db ./cmd/cmdopts/.migrations
//go:generate genieql bootstrap --queryer=sqlx.Queryer --driver=github.com/marcboeker/go-duckdb duckdb://localhost/dpool.db
//go:generate genieql auto graph -o genieql.gen.go
