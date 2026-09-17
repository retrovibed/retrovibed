package retrovibed

//go:generate genieql duckdb --extension=icu --extension=vss --extension=inet --database=dpool.db ./cmd/cmdopts/.migrations
//go:generate genieql bootstrap --queryer=sqlx.Queryer --driver=github.com/marcboeker/go-duckdb duckdb://localhost/dpool.db
//go:generate genieql auto graph -o genieql.gen.go

//go:generate genieql duckdb --extension=icu --extension=vss --extension=inet --database=cache.db ./cmd/cmdopts/.migrations.cache
//go:generate genieql bootstrap --output-file=cache.config --queryer=sqlx.Queryer --driver=github.com/marcboeker/go-duckdb duckdb://localhost/cache.db
//go:generate genieql auto graph --config cache.config --tags retrovibed.cache -o genieql.gen.go
