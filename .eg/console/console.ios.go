package console

import (
	"context"
	"os"
	"path/filepath"

	"github.com/egdaemon/eg/runtime/wasi/eg"
	"github.com/egdaemon/eg/runtime/wasi/shell"
)

const inetCMakeConfig = `duckdb_extension_load(inet
    SOURCE_DIR ${CMAKE_CURRENT_LIST_DIR}/../../inet
    LOAD_TESTS
)
`

// EnsureDuckDBSource clones DuckDB and inet extension sources, writes the
// extension_config_local.cmake so inet is built as a static extension.
func EnsureDuckDBSource(duckdbsrc, duckdbinet string) eg.OpFn {
	return func(ctx context.Context, op eg.Op) error {
		if err := shell.Op(
			shell.Newf("git clone --depth 1 --branch v1.5.3 https://github.com/duckdb/duckdb.git %s", duckdbsrc).Lenient(true),
			shell.Newf("git clone --depth 1 --branch v1.4-andium https://github.com/duckdb/duckdb_inet.git %s", duckdbinet).Lenient(true),
		)(ctx, op); err != nil {
			return err
		}

		return os.WriteFile(filepath.Join(duckdbsrc, "extension", "extension_config_local.cmake"), []byte(inetCMakeConfig), 0644)
	}
}

func BuildIOS(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter build ios --release --no-codesign"),
	)
}

func BuildIOSSimulator(ctx context.Context, _ eg.Op) error {
	runtime := flutterRuntimev2(shell.Runtime())
	return shell.Run(
		ctx,
		runtime.New("flutter build ios --debug --simulator"),
	)
}
