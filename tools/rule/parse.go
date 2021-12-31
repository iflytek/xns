package rule

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
)

func parseRule(rule string) error {

	expr,err:=parser.ParseExpr(rule)
	if err != nil{
		return err
	}

	switch expr {

	}
}

func parseExpr(expr ast.Expr)(Rule,error){
	switch expr.(type) {
	case *ast.BinaryExpr:
		be:=expr.(*ast.BinaryExpr)
		switch be.Op {
		case token.EQL:
			rule:=&Matches{}
		}
	case ast.Field:


	}
}
