package ginterp

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"

	"github.com/james-lawrence/genieql/internal/envx"
)

func WasiPackage() *build.Package {
	return &build.Package{
		Dir:           envx.String("", "GENIEQL_WASI_PACKAGE_DIR"),
		Name:          envx.String("", "GENIEQL_WASI_PACKAGE_NAME"),
		ImportComment: envx.String("", "GENIEQL_WASI_PACKAGE_IMPORT_COMMENT"),
		Doc:           envx.String("", "GENIEQL_WASI_PACKAGE_DOC"),
		ImportPath:    envx.String("", "GENIEQL_WASI_PACKAGE_IMPORT_PATH"),
		Root:          envx.String("", "GENIEQL_WASI_PACKAGE_ROOT"),
		SrcRoot:       envx.String("", "GENIEQL_WASI_PACKAGE_SRC_ROOT"),
		PkgRoot:       envx.String("", "GENIEQL_WASI_PACKAGE_PKG_ROOT"),
		PkgTargetRoot: envx.String("", "GENIEQL_WASI_PACKAGE_PKG_TARGET_ROOT"),
		BinDir:        envx.String("", "GENIEQL_WASI_PACKAGE_BIN_DIR"),
		Goroot:        envx.Boolean(false, "GENIEQL_WASI_PACKAGE_GO_ROOT"),
		PkgObj:        envx.String("", "GENIEQL_WASI_PACKAGE_PKG_OBJ"),
		AllTags:       envx.Strings(nil, "GENIEQL_WASI_PACKAGE_ALL_TAGS"),
		ConflictDir:   envx.String("", "GENIEQL_WASI_PACKAGE_CONFLICT_DIR"),
		BinaryOnly:    envx.Boolean(false, "GENIEQL_WASI_PACKAGE_BINARY_ONLY"),
		GoFiles:       envx.Strings(nil, "GENIEQL_WASI_PACKAGE_GO_FILES"),
	}
}

func LoadFile() (*ast.File, error) {
	fset := token.NewFileSet()
	fp := envx.String("", "GENIEQL_WASI_FILEPATH")
	// log.Println("LOADING FILE", fp)
	// fsx.PrintDir(os.DirFS("."))
	// fsx.PrintDir(os.DirFS(filepath.Dir(fp)))
	// fsx.PrintString(fp)
	// fsx.PrintString(filepath.Join(filepath.Dir(fp), "src", "main.go"))
	return parser.ParseFile(fset, fp, nil, parser.ParseComments)
}
