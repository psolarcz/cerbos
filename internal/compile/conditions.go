// Copyright 2021-2024 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package compile

import (
	"github.com/google/cel-go/cel"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	policyv1 "github.com/cerbos/cerbos/api/genpb/cerbos/policy/v1"
	runtimev1 "github.com/cerbos/cerbos/api/genpb/cerbos/runtime/v1"
	"github.com/cerbos/cerbos/internal/conditions"
)

func Condition(cond *policyv1.Condition) (*runtimev1.Condition, error) {
	mc := &moduleCtx{unitCtx: &unitCtx{errors: new(ErrorList)}, fqn: "UNKNOWN", sourceFile: "UNKNOWN"}
	cc := compileCondition(mc, "unknown", cond, false)
	return cc, mc.error()
}

func compileCondition(modCtx *moduleCtx, parent string, cond *policyv1.Condition, markReferencedVariablesAsUsed bool) *runtimev1.Condition {
	if cond == nil {
		return nil
	}

	switch c := cond.Condition.(type) {
	case *policyv1.Condition_Match:
		return compileMatch(modCtx, parent, c.Match, markReferencedVariablesAsUsed)
	default:
		modCtx.addErrWithDesc(errScriptsUnsupported, "Unsupported feature in %s", parent)
		return nil
	}
}

func compileMatch(modCtx *moduleCtx, parent string, match *policyv1.Match, markReferencedVariablesAsUsed bool) *runtimev1.Condition {
	if match == nil {
		return nil
	}

	switch t := match.Op.(type) {
	case *policyv1.Match_Expr:
		expr := &runtimev1.Expr{Original: t.Expr, Checked: compileCELExpr(modCtx, parent, t.Expr, markReferencedVariablesAsUsed)}
		return &runtimev1.Condition{Op: &runtimev1.Condition_Expr{Expr: expr}}
	case *policyv1.Match_All:
		exprList := compileMatchList(modCtx, parent, t.All.Of, markReferencedVariablesAsUsed)
		return &runtimev1.Condition{Op: &runtimev1.Condition_All{All: exprList}}
	case *policyv1.Match_Any:
		exprList := compileMatchList(modCtx, parent, t.Any.Of, markReferencedVariablesAsUsed)
		return &runtimev1.Condition{Op: &runtimev1.Condition_Any{Any: exprList}}
	case *policyv1.Match_None:
		exprList := compileMatchList(modCtx, parent, t.None.Of, markReferencedVariablesAsUsed)
		return &runtimev1.Condition{Op: &runtimev1.Condition_None{None: exprList}}
	default:
		modCtx.addErrWithDesc(errUnexpectedErr, "Unknown match operation in %s: %T", parent, t)
		return nil
	}
}

func compileCELExpr(modCtx *moduleCtx, parent, expr string, markReferencedVariablesAsUsed bool) *exprpb.CheckedExpr {
	celAST, issues := conditions.StdEnv.Compile(expr)
	if issues != nil && issues.Err() != nil {
		modCtx.addErrWithDesc(newCELCompileError(expr, issues), "Invalid expression in %s", parent)
		return nil
	}

	checkedExpr, err := cel.AstToCheckedExpr(celAST)
	if err != nil {
		modCtx.addErrWithDesc(err, "Failed to convert AST of `%s` in %s", expr, parent)
		return nil
	}

	if markReferencedVariablesAsUsed {
		modCtx.variables.Use(parent, checkedExpr)
	}

	return checkedExpr
}

func compileMatchList(modCtx *moduleCtx, parent string, matches []*policyv1.Match, markReferencedVariablesAsUsed bool) *runtimev1.Condition_ExprList {
	exprList := make([]*runtimev1.Condition, len(matches))
	for i, m := range matches {
		exprList[i] = compileMatch(modCtx, parent, m, markReferencedVariablesAsUsed)
	}

	return &runtimev1.Condition_ExprList{Expr: exprList}
}
