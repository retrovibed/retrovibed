package ginterp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"io"

	"github.com/james-lawrence/genieql"
	"github.com/james-lawrence/genieql/astcodec"
	"github.com/james-lawrence/genieql/astutil"
	"github.com/james-lawrence/genieql/generators"
	"github.com/james-lawrence/genieql/generators/functions"
	"github.com/james-lawrence/genieql/generators/typespec"
	"github.com/james-lawrence/genieql/internal/errorsx"
	"github.com/james-lawrence/genieql/internal/stringsx"
	"github.com/pkg/errors"
)

// InsertBatch configuration interface for generating batch inserts.
type InsertBatch interface {
	genieql.Generator              // must satisfy the generator interface
	Into(string) InsertBatch       // what table to insert into
	Default(...string) InsertBatch // use the database default for the specified columns.
	Conflict(string) InsertBatch   // specify how conflicts should be handled.
	Batch(n int) InsertBatch       // specify a batch insert
}

// NewInsert instantiate a new insert generator. it uses the name of function
// that calls Define as the name of the generated function.
func NewBatchInsert(
	ctx generators.Context,
	name string,
	comment *ast.CommentGroup,
	cf *ast.Field,
	qf *ast.Field,
	tf *ast.Field,
	scanner *ast.FuncDecl,
) InsertBatch {
	return &batch{
		ctx:     ctx,
		name:    name,
		comment: comment,
		qf:      qf,
		cf:      cf,
		tf:      tf,
		scanner: scanner,
		n:       1,
	}
}

func InsertBatchFromFile(cctx generators.Context, name string, tree *ast.File) (InsertBatch, error) {
	var (
		ok          bool
		pos         *ast.FuncDecl
		scanner     *ast.FuncDecl // scanner to use for the results.
		declPattern *ast.FuncType
	)

	if pos = astcodec.FileFindDecl[*ast.FuncDecl](tree, astcodec.FindFunctionsByName(name)); pos == nil {
		return nil, fmt.Errorf("unable to locate function declaration for insert: %s", name)
	}

	// rewrite scanner declaration function.
	if declPattern, ok = pos.Type.Params.List[1].Type.(*ast.FuncType); !ok {
		return nil, errorsx.String("InsertBatch second parameter must be a function type")
	}

	if scanner = functions.DetectScanner(cctx, declPattern); scanner == nil {
		return nil, errors.Errorf("InsertBatch %s - missing scanner", nodeInfo(cctx, pos))
	}

	return NewBatchInsert(
		cctx,
		pos.Name.String(),
		pos.Doc,
		functions.DetectContext(declPattern),
		functions.DetectQueryer(declPattern),
		declPattern.Params.List[len(declPattern.Params.List)-1],
		scanner,
	), nil
}

type batch struct {
	ctx      generators.Context
	n        int // number of records to support inserting
	name     string
	table    string
	conflict string
	defaults []string
	tf       *ast.Field    // type field.
	cf       *ast.Field    // context field, can be nil.
	qf       *ast.Field    // db Query field.
	scanner  *ast.FuncDecl // scanner being used for results.
	comment  *ast.CommentGroup
}

// Into specify the table the data will be inserted into.
func (t *batch) Into(s string) InsertBatch {
	t.table = s
	return t
}

// Default specify the table columns to be given their default values.
func (t *batch) Default(defaults ...string) InsertBatch {
	t.defaults = defaults
	return t
}

// Conflict specify how to handle conflict during an insert.
func (t *batch) Conflict(s string) InsertBatch {
	t.conflict = s
	return t
}

// Batch specify the maximum number of records to insert.
func (t *batch) Batch(size int) InsertBatch {
	t.n = size
	return t
}

func (t *batch) Generate(dst io.Writer) (err error) {
	var (
		cmaps       []genieql.ColumnMap
		queryfields []*ast.Field
		encodings   []ast.Stmt
		explodedecl *ast.FuncDecl
	)
	t.ctx.Println("generation of", t.name, "initiated")
	defer t.ctx.Println("generation of", t.name, "completed")
	t.ctx.Debugln("batch.insert type", t.ctx.CurrentPackage.Name, t.ctx.CurrentPackage.ImportPath, types.ExprString(t.tf.Type))
	t.ctx.Debugln("batch.insert table", t.table)
	t.ctx.Debugln("batch.insert type", t.tf.Names[0])
	t.ctx.Debugln("batch.insert scanner", t.scanner)

	if cmaps, err = generators.ColumnMapFromFields(t.ctx, t.tf); err != nil {
		return errors.Wrap(err, "unable to generate mapping")
	}

	defaulted := genieql.ColumnInfoFilterIgnore(t.defaults...)

	cset := genieql.ColumnMapSet(cmaps)
	defaultedcset := cset.Filter(func(cm genieql.ColumnMap) bool { return defaulted(cm.ColumnInfo) })

	queryfields = generators.QueryFieldsFromColumnMap(t.ctx, defaultedcset.Map(func(idx int, cm genieql.ColumnMap) genieql.ColumnMap {
		local := cm.Local(idx)
		dup := cm
		dup.Field = astutil.Field(astutil.MustParseExpr(t.ctx.FileSet, cm.Definition.ColumnType), local)
		return dup
	})...)

	explodeerrHandler := func(errlocal string) ast.Node {
		explodereturn := make([]ast.Expr, 0, len(queryfields)+1)
		explodereturn = append(explodereturn, astutil.MapFieldsToNameExpr(queryfields...)...)
		explodereturn = append(explodereturn, ast.NewIdent(errlocal))
		return astutil.Return(explodereturn...)
	}

	if _, encodings, _, err = generators.QueryInputsFromColumnMap(t.ctx, t.scanner, explodeerrHandler, defaultedcset...); err != nil {
		return errors.Wrap(err, "unable to transform query inputs")
	}

	errhandling := generators.ScannerErrorHandlingExpr(t.scanner)

	initializesig := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				t.cf,
				t.qf,
				astutil.Field(&ast.Ellipsis{
					Elt: t.tf.Type,
				}, t.tf.Names...),
			},
		},
		Results: t.scanner.Type.Results,
	}

	typename := stringsx.ToPrivate(t.name)
	initialize := functions.NewFn(
		astutil.Return(
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: ast.NewIdent(typename),
					Elts: []ast.Expr{
						&ast.KeyValueExpr{
							Key:   t.cf.Names[0],
							Value: t.cf.Names[0],
						},
						&ast.KeyValueExpr{
							Key:   t.qf.Names[0],
							Value: t.qf.Names[0],
						},
						&ast.KeyValueExpr{
							Key:   ast.NewIdent("remaining"),
							Value: t.tf.Names[0],
						},
					},
				},
			},
		),
	)

	typedecl := typespec.NewType(typename, &ast.StructType{
		Struct: token.Pos(0),
		Fields: &ast.FieldList{
			List: []*ast.Field{
				t.cf,
				t.qf,
				astutil.Field(
					&ast.ArrayType{Elt: t.tf.Type}, ast.NewIdent("remaining")),
				astutil.Field(t.scanner.Type.Results.List[0].Type, ast.NewIdent("scanner")),
			},
		},
	})

	fnrecv := astutil.FieldList(astutil.Field(&ast.StarExpr{X: astutil.Expr(typename)}, ast.NewIdent("t")))

	scansig := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(&ast.StarExpr{X: t.tf.Type}, t.tf.Names...),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(ast.NewIdent("error")),
			},
		},
	}
	scanfn := functions.NewFn(
		astutil.Return(
			astutil.CallExpr(
				&ast.SelectorExpr{
					X: astutil.SelExpr(
						"t", "scanner",
					),
					Sel: ast.NewIdent("Scan"),
				},
				t.tf.Names[0],
			),
		),
	)

	errsig := &ast.FuncType{
		Params: &ast.FieldList{},
		Results: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(ast.NewIdent("error")),
			},
		},
	}

	errfn := functions.NewFn(
		astutil.If(
			nil,
			astutil.BinaryExpr(astutil.SelExpr("t", "scanner"), token.EQL, ast.NewIdent("nil")),
			astutil.Block(
				astutil.Return(ast.NewIdent("nil")),
			),
			nil,
		),
		astutil.Return(
			astutil.CallExpr(
				astutil.SelExpr(
					types.ExprString(
						astutil.SelExpr(
							"t", "scanner",
						),
					),
					"Err",
				),
			),
		),
	)

	closesig := &ast.FuncType{
		Params: &ast.FieldList{},
		Results: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(ast.NewIdent("error")),
			},
		},
	}

	closefn := functions.NewFn(
		astutil.If(
			nil,
			astutil.BinaryExpr(astutil.SelExpr("t", "scanner"), token.EQL, ast.NewIdent("nil")),
			astutil.Block(
				astutil.Return(ast.NewIdent("nil")),
			),
			nil,
		),
		astutil.Return(
			astutil.CallExpr(
				astutil.SelExpr(
					types.ExprString(
						astutil.SelExpr(
							"t", "scanner",
						),
					),
					"Close",
				),
			),
		),
	)

	nextsig := &ast.FuncType{
		Params: &ast.FieldList{},
		Results: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(ast.NewIdent("bool")),
			},
		},
	}

	nextfn := functions.NewFn(
		astutil.DeclStmt(
			astutil.VarList(
				astutil.ValueSpec(ast.NewIdent("bool"), ast.NewIdent("advanced")),
			),
		),
		astutil.If(
			nil, astutil.BinaryExpr(
				astutil.BinaryExpr(astutil.SelExpr("t", "scanner"), token.NEQ, ast.NewIdent("nil")),
				token.LAND,
				astutil.CallExpr(
					&ast.SelectorExpr{
						X:   astutil.SelExpr("t", "scanner"),
						Sel: ast.NewIdent("Next"),
					},
				),
			),
			astutil.Block(
				astutil.Return(ast.NewIdent("true")),
			),
			nil,
		),
		astutil.If(
			nil, astutil.BinaryExpr(
				astutil.BinaryExpr(astutil.CallExpr(ast.NewIdent("len"), astutil.SelExpr("t", "remaining")), token.GTR, astutil.IntegerLiteral(0)),
				token.LAND,
				astutil.BinaryExpr(
					astutil.CallExpr(
						astutil.SelExpr("t", "Close"),
					),
					token.EQL,
					ast.NewIdent("nil"),
				),
			),
			astutil.Block(
				astutil.Assign(
					astutil.ExprList(
						astutil.SelExpr("t", "scanner"),
						astutil.SelExpr("t", "remaining"),
						ast.NewIdent("advanced"),
					),
					token.ASSIGN,
					astutil.ExprList(
						astutil.CallExprEllipsis(
							astutil.SelExpr("t", "advance"),
							astutil.SelExpr("t", "remaining"),
						),
					),
				),
				astutil.Return(
					astutil.BinaryExpr(ast.NewIdent("advanced"), token.LAND, astutil.CallExpr(&ast.SelectorExpr{
						X:   astutil.SelExpr("t", "scanner"),
						Sel: ast.NewIdent("Next"),
					})),
				),
			),
			nil,
		),
		astutil.Return(
			ast.NewIdent("false"),
		),
	)

	advancesig := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				astutil.Field(&ast.Ellipsis{Elt: t.tf.Type}, t.tf.Names[0]),
			},
		},
		Results: &ast.FieldList{
			List: []*ast.Field{
				t.scanner.Type.Results.List[0],
				astutil.Field(&ast.ArrayType{Elt: t.tf.Type}),
				astutil.Field(ast.NewIdent("bool")),
			},
		},
	}

	explodereturn := make([]ast.Expr, 0, len(queryfields)+1)
	explodereturn = append(explodereturn, astutil.MapFieldsToNameExpr(queryfields...)...)
	explodereturn = append(explodereturn, ast.NewIdent("nil"))
	explodestmts := make([]ast.Stmt, 0, len(encodings)+1)
	explodestmts = append(explodestmts, encodings...)
	explodestmts = append(explodestmts, astutil.Return(explodereturn...))

	explodefn := functions.NewFn(
		explodestmts...,
	)

	exploderesults := make([]*ast.Field, len(queryfields), len(queryfields)+1)
	copy(exploderesults, queryfields)
	exploderesults = append(exploderesults, astutil.Field(ast.NewIdent("error"), ast.NewIdent("err")))

	explodesig := &ast.FuncType{
		Params: &ast.FieldList{
			List: []*ast.Field{
				t.tf,
			},
		},
		Results: &ast.FieldList{
			List: exploderesults,
		},
	}

	if explodedecl, err = explodefn.Compile(functions.New("", explodesig)); err != nil {
		return errors.Wrap(err, "failed to generate encoding function")
	}

	generatecasestatement := func(nrecords int, caseexpr ...ast.Expr) *ast.CaseClause {
		var (
			qinputs = astutil.ExprList(
				ast.NewIdent("t.ctx"),
				ast.NewIdent("query"),
			)
			genlocals   []ast.Spec
			assignments []ast.Stmt
			remaining   func(...ast.Expr) ast.Stmt = func(inputs ...ast.Expr) ast.Stmt {
				return astutil.Return(
					astutil.CallExpr(
						t.scanner.Name,
						astutil.CallExpr(
							&ast.SelectorExpr{
								X: astutil.SelExpr(
									"t",
									"q",
								),
								Sel: ast.NewIdent("QueryContext"),
							},
							inputs...,
						),
					),
					&ast.SliceExpr{
						X:   t.tf.Names[0],
						Low: astutil.IntegerLiteral(nrecords),
					},
					ast.NewIdent("true"),
				)
			}
		)

		if nrecords == t.n {
			remaining = func(inputs ...ast.Expr) ast.Stmt {
				return astutil.Return(
					astutil.CallExpr(
						t.scanner.Name,
						astutil.CallExpr(
							&ast.SelectorExpr{
								X: astutil.SelExpr(
									"t",
									"q",
								),
								Sel: ast.NewIdent("QueryContext"),
							},
							inputs...,
						),
					),
					astutil.CallExpr(
						&ast.ArrayType{Elt: t.tf.Type},
						ast.NewIdent("nil"),
					),
					ast.NewIdent("false"),
				)
			}
		}

		for j := 0; j < nrecords; j++ {
			localqf1 := astutil.TransformFields(func(f *ast.Field) *ast.Field {
				return astutil.Field(f.Type, ast.NewIdent(fmt.Sprintf("r%d%s", j, f.Names[0])))
			}, queryfields...)
			inputs := astutil.MapFieldsToNameExpr(localqf1...)
			localqf2 := astutil.MapFieldsToValueSpec(localqf1...)
			localqf3 := astutil.MapValueSpecToSpec(localqf2...)
			assignment := astutil.ExprList(
				inputs...,
			)
			assignment = append(assignment, ast.NewIdent("err"))
			genlocals = append(genlocals, localqf3...)
			assignments = append(
				assignments,
				astutil.If(
					astutil.Assign(
						assignment,
						token.ASSIGN,
						astutil.ExprList(
							astutil.CallExpr(
								ast.NewIdent("transform"),
								&ast.IndexListExpr{
									X:       t.tf.Names[0],
									Indices: astutil.ExprList(astutil.IntegerLiteral(j)),
								},
							),
						),
					),
					astutil.BinaryExpr(
						ast.NewIdent("err"), token.NEQ, ast.NewIdent("nil"),
					),
					astutil.Block(
						astutil.Return(
							errhandling("err"),
							astutil.CallExpr(
								&ast.ArrayType{Elt: t.tf.Type},
								ast.NewIdent("nil"),
							),
							ast.NewIdent("false"),
						),
					),
					nil,
				),
			)

			qinputs = append(qinputs, inputs...)
		}

		queryreplacement := functions.QueryLiteralColumnMapReplacer(t.ctx, t.ctx.Dialect.Insert(nrecords, 0, t.table, t.conflict, cset.ColumnNames(), cset.ColumnNames(), t.defaults), cmaps...)
		casestmts := make([]ast.Stmt, 0, len(assignments)+3)
		casestmts = append(casestmts, astutil.DeclStmt(genieql.QueryLiteral(
			"query",
			queryreplacement,
		)))
		casestmts = append(casestmts, astutil.DeclStmt(astutil.VarList(append(genlocals, astutil.ValueSpec(ast.NewIdent("error"), ast.NewIdent("err")))...)))
		casestmts = append(casestmts, assignments...)
		casestmts = append(casestmts, remaining(qinputs...))

		return astutil.CaseClause(
			caseexpr,
			casestmts...,
		)
	}

	genscanning := func() *ast.BlockStmt {
		stmts := make([]ast.Stmt, 0, t.n)
		stmts = append(stmts, astutil.CaseClause(
			astutil.ExprList(astutil.IntegerLiteral(0)),
			astutil.Return(
				ast.NewIdent("nil"),
				astutil.CallExpr(
					&ast.ArrayType{Elt: t.tf.Type},
					ast.NewIdent("nil"),
				),
				ast.NewIdent("false"),
			),
		))

		for i := 1; len(stmts) < cap(stmts); i++ {
			stmts = append(stmts, generatecasestatement(i, astutil.IntegerLiteral(i)))
		}

		stmts = append(stmts, generatecasestatement(t.n))
		return astutil.Block(stmts...)
	}

	advancefn := functions.NewFn(
		astutil.Assign(
			astutil.ExprList(ast.NewIdent("transform")),
			token.DEFINE,
			astutil.ExprList(astutil.FuncLiteral(explodedecl)),
		),
		astutil.Switch(
			nil,
			astutil.CallExpr(
				ast.NewIdent("len"),
				ast.NewIdent(t.tf.Names[0].String()),
			),
			genscanning(),
		),
	)

	return genieql.NewFuncGenerator(func(dst io.Writer) (err error) {
		if err = generators.GenerateComment(generators.DefaultFunctionComment(t.name), t.comment).Generate(dst); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("New"+stringsx.ToPublic(t.name), initializesig), initialize); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = typespec.CompileInto(dst, typedecl); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("Scan", scansig, functions.OptionRecv(fnrecv)), scanfn); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("Err", errsig, functions.OptionRecv(fnrecv)), errfn); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("Close", closesig, functions.OptionRecv(fnrecv)), closefn); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("Next", nextsig, functions.OptionRecv(fnrecv)), nextfn); err != nil {
			return err
		}

		if err = generators.GapLines(dst, 2); err != nil {
			return err
		}

		if err = functions.CompileInto(dst, functions.New("advance", advancesig, functions.OptionRecv(fnrecv)), advancefn); err != nil {
			return err
		}

		return nil
	}).Generate(dst)
}
