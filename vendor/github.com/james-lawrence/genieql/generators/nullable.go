package generators

import (
	"go/ast"
	"go/types"

	"github.com/james-lawrence/genieql"
	"github.com/pkg/errors"
)

// tdRegistry type definition registry
type tdRegistry func(s string) (genieql.ColumnDefinition, error)

func composeTypeDefinitionsExpr(definitions ...tdRegistry) genieql.LookupTypeDefinition {
	return func(e ast.Expr) (d genieql.ColumnDefinition, err error) {
		for _, registry := range definitions {
			if d, err = registry(types.ExprString(e)); err == nil {
				return d, nil
			}
		}

		return d, errors.Errorf("failed to locate type information for expr %s", types.ExprString(e))
	}
}

func composeTypeDefinitions(definitions ...tdRegistry) tdRegistry {
	return func(e string) (d genieql.ColumnDefinition, err error) {
		for _, registry := range definitions {
			if d, err = registry(e); err == nil {
				return d, nil
			}
		}

		return d, errors.Errorf("failed to locate type information for %s", e)
	}
}
