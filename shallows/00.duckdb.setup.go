package retrovibed

//go:generate genieql duckdb --database=dpool.db ./cmd/cmdmeta/.migrations
//go:generate genieql bootstrap --queryer=sqlx.Queryer --driver=github.com/marcboeker/go-duckdb duckdb://localhost/dpool.db
