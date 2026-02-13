package loglint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func NewAnalyzer(opts *Options) *analysis.Analyzer {
	if opts == nil {
		opts = DefaultOptions()
	}

	return &analysis.Analyzer{
		Name:     "loglint",
		Doc:      "checks slog/zap log calls using a set of rules",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			ins := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

			ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
				callExpr := n.(*ast.CallExpr)

				lc, ok := ExtractLogCall(pass, callExpr)
				if !ok {
					return
				}

				for _, rule := range rules {
					rule(pass, opts, &lc)
				}
			})

			return nil, nil
		},
	}
}
