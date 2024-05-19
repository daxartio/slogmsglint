package slogmsglint

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "slogmsglint",
	Doc:  "checks that all slog messages are in lowercase",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			methodName := selExpr.Sel.Name
			if !isLogMethod(methodName) {
				return true
			}

			if len(callExpr.Args) < 1 {
				return true
			}

			argIndex := 0
			if strings.HasSuffix(methodName, "Context") {
				argIndex = 1
			} else if methodName == "LogAttrs" || methodName == "Log" {
				argIndex = 2
			}

			if argIndex >= len(callExpr.Args) {
				return true
			}

			basicLit, ok := callExpr.Args[argIndex].(*ast.BasicLit)
			if !ok || basicLit.Kind != token.STRING {
				return true
			}

			msg := strings.Trim(basicLit.Value, "\"")
			if msg != strings.ToLower(msg) {
				pass.Reportf(basicLit.Pos(), "slog message should be in lowercase: %s", msg)
			}

			if strings.HasSuffix(msg, ".") {
				pass.Reportf(basicLit.Pos(), "slog message should not end with a period: %s", msg)
			}

			return true
		})
	}
	return nil, nil
}

func isLogMethod(name string) bool {
	switch name {
	case "Log", "LogAttrs", "Debug", "DebugContext", "Info", "InfoContext", "Warn", "WarnContext", "Error", "ErrorContext":
		return true
	default:
		return false
	}
}
