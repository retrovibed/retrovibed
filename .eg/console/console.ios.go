package console

import (
	"bytes"
	"context"
	"encoding/hex"
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
			shell.Newf("git clone --depth 1 --branch v1.4.3 https://github.com/duckdb/duckdb.git %s", duckdbsrc).Lenient(true),
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

// RetagSimulator rewrites LC_BUILD_VERSION platform from macOS to iOS-simulator
// across each object in a static archive at src, producing dst. Cgo c-archive
// output from GOOS=darwin is tagged macOS; Xcode's simulator linker refuses it.
// cmdsize is unchanged so a direct byte substitution is sufficient and avoids
// vtool's header-pad limitations.
func RetagSimulator(src, dst string) eg.OpFn {
	return func(ctx context.Context, _ eg.Op) error {
		macOS, _ := hex.DecodeString("320000001800000001000000")
		iosSim, _ := hex.DecodeString("320000001800000007000000")

		work, err := os.MkdirTemp("", "retag-sim-")
		if err != nil {
			return err
		}
		defer os.RemoveAll(work)

		if err := shell.Run(ctx, shell.Newf("cd %s && ar x %s", work, src)); err != nil {
			return err
		}

		entries, err := os.ReadDir(work)
		if err != nil {
			return err
		}

		for _, e := range entries {
			if filepath.Ext(e.Name()) != ".o" {
				continue
			}
			p := filepath.Join(work, e.Name())
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			if !bytes.Contains(data, macOS) {
				continue
			}
			if err := os.WriteFile(p, bytes.ReplaceAll(data, macOS, iosSim), 0644); err != nil {
				return err
			}
		}

		return shell.Run(ctx, shell.Newf("cd %s && ar rcs %s *.o", work, dst))
	}
}
