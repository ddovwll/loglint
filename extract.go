package loglint

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var slogMsgIndex = map[string]int{
	"Debug": 0, "Info": 0, "Warn": 0, "Error": 0,
	"DebugContext": 1, "InfoContext": 1, "WarnContext": 1, "ErrorContext": 1,
	"Log": 2, "LogAttrs": 2,
}

var zapMsgIndex = map[string]int{
	"Debug": 0, "Info": 0, "Warn": 0, "Error": 0, "DPanic": 0, "Panic": 0, "Fatal": 0,
	"Debugf": 0, "Infof": 0, "Warnf": 0, "Errorf": 0, "DPanicf": 0, "Panicf": 0, "Fatalf": 0,
	"Debugw": 0, "Infow": 0, "Warnw": 0, "Errorw": 0, "DPanicw": 0, "Panicw": 0, "Fatalw": 0,
	"Debugln": 0, "Infoln": 0, "Warnln": 0, "Errorln": 0, "DPanicln": 0, "Panicln": 0, "Fatalln": 0,
}

func ExtractLogCall(pass *analysis.Pass, call *ast.CallExpr) (LogCall, bool) {
	fn, pkgPath, ok := calledFunc(pass, call)
	if !ok || fn == nil || fn.Pkg() == nil {
		return LogCall{}, false
	}

	var (
		msgIdx int
		isLog  bool
	)

	switch pkgPath {
	case "log/slog":
		msgIdx, isLog = slogMsgIndex[fn.Name()]
	case "go.uber.org/zap":
		msgIdx, isLog = zapMsgIndex[fn.Name()]
	default:
		return LogCall{}, false
	}

	if !isLog || msgIdx < 0 || msgIdx >= len(call.Args) {
		return LogCall{}, false
	}

	msgExpr := call.Args[msgIdx]
	msg, isConst, isLiteral := extractString(pass, msgExpr)

	return LogCall{
		Msg:          msg,
		MsgIsNotVar:  isConst,
		MsgIsLiteral: isLiteral,
		MsgExprPos:   msgExpr.Pos(),
		MsgExprEnd:   msgExpr.End(),
	}, true
}

func calledFunc(pass *analysis.Pass, call *ast.CallExpr) (*types.Func, string, bool) {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		if obj, ok := pass.TypesInfo.Uses[fun].(*types.Func); ok && obj.Pkg() != nil {
			return obj, obj.Pkg().Path(), true
		}
	case *ast.SelectorExpr:
		if sel := pass.TypesInfo.Selections[fun]; sel != nil {
			if obj, ok := sel.Obj().(*types.Func); ok && obj.Pkg() != nil {
				return obj, obj.Pkg().Path(), true
			}
		}
		if obj, ok := pass.TypesInfo.Uses[fun.Sel].(*types.Func); ok && obj.Pkg() != nil {
			return obj, obj.Pkg().Path(), true
		}
	}
	return nil, "", false
}

func extractString(pass *analysis.Pass, expr ast.Expr) (msg string, isConst bool, isLiteral bool) {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return "", false, false
		}
		return s, true, true
	}

	if tv, ok := pass.TypesInfo.Types[expr]; ok && tv.Value != nil && tv.Value.Kind() == constant.String {
		return constant.StringVal(tv.Value), true, false
	}

	return "", false, false
}
