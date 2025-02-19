package deeppool

//go:generate genieql duckdb --database=dpool.db ./cmd/shallows/.migrations
//go:generate genieql bootstrap --queryer=sqlx.Queryer --driver=github.com/marcboeker/go-duckdb duckdb://localhost/dpool.db
