package parser

import (
	"errors"
	"fmt"
	"wabbit-go/common"
	"wabbit-go/model"
	"wabbit-go/tokenize"
)

type EOF struct {
	Type   string
	Value  string
	Lineno string
	Index  int
}

var eof = EOF{"EOF", "EOF", "EOF", -1}

type TokenStream struct {
	program    *model.Program
	gen        chan tokenize.Token
	lookahead  tokenize.Token
	lastIndex  int
	haveErrors bool
}

func NewTokenStream(program *model.Program) (*TokenStream, error) {
	tokens, err := tokenize.Tokenize(program.Source)
	if err != nil {
		return nil, err
	}
	ts := &TokenStream{
		program:    program,
		gen:        common.SliceToChannel(tokens),
		haveErrors: false,
	}
	ts.lookahead = <-ts.gen
	return ts, nil
}

func (ts *TokenStream) Peek(types ...string) *tokenize.Token {
	for _, t := range types {
		if ts.lookahead.Type == t {
			return &ts.lookahead
		}
	}
	return nil
}

func (ts *TokenStream) Accept(types ...string) *tokenize.Token {
	tok := ts.Peek(types...)
	if tok != nil {
		ts.lastIndex = ts.lookahead.Index + len(ts.lookahead.Value)
		ts.lookahead = <-ts.gen
	}
	return tok
}

func (ts *TokenStream) Expect(types ...string) (*tokenize.Token, error) {
	tok := ts.Accept(types...)
	if tok == nil {
		return nil, errors.New(fmt.Sprintf("Syntax error at %s expect %v but got %v", ts.lookahead.Value, types, tok))
	}
	return tok, nil
}

func (ts *TokenStream) Synchronize(types ...string) {
	for ts.Accept(types...) == nil {
		ts.lookahead = <-ts.gen
	}
}

// type nodeInitFunc func(args ...interface{}) model.Node
// pass it back
type constructFunc func(model.Node) model.Node

func (ts *TokenStream) Builder() func(func(constructFunc) model.Node) {
	startLine := ts.lookahead.Lineno
	startIndex := ts.lookahead.Index

	construct := func(node model.Node) model.Node {
		//print("nodetype=", nodetype)
		//node := nodeInit(args)
		// we got node don't init any more
		ts.program.RecordPosition(node, startLine, startIndex, ts.lastIndex) // ts.lookahead.index
		return node
	}

	return func(do func(constructFunc) model.Node) {
		//defer func() {
		//	if err := recover(); err != nil {
		//		if err, ok := err.(SyntaxError); ok {
		//			fmt.Printf("%d: %v", startLine, err)
		//			ts.haveErrors = true
		//		} else {
		//			panic(err)
		//		}
		//	}
		//}()
		do(construct)
	}
}

func ParseProgram(program *model.Program) error {
	// Code this to recognize any Wabbit program and build AST
	ts, err := NewTokenStream(program)
	if err != nil {
		return err
	}
	statements := parseStatements(ts)
	_, err = ts.Expect("EOF")
	if err != nil {
		return err
	}
	program.Model = statements
	return nil
}

func parseStatements(ts *TokenStream) model.Node {
	statements := []model.Statement{}
	ts.Builder()(func(new constructFunc) model.Node {
		for ts.Peek("RBRACE", "EOF") != nil {
			statement := parseStatement(ts)
			statements = append(statements, statement)
		}
		return new(&model.Statements{statements})
	})
}

func parseStatement(ts *TokenStream) Node {
	// Parse different types of statements
	if ts.Peek("PRINT") {
		return parsePrintStmt(ts)
	} else if ts.Peek("CONST") {
		return parseConstDecl(ts)
	} else if ts.Peek("VAR") {
		return parseVarDecl(ts)
	} else if ts.Peek("RETURN") {
		return parseReturnStmt(ts)
	} else if ts.Peek("IF") {
		return parseIfStmt(ts)
	} else if ts.Peek("WHILE") {
		return parseWhileStmt(ts)
	} else if ts.Peek("BREAK") {
		return parseBreakStmt(ts)
	} else if ts.Peek("CONTINUE") {
		return parseContinueStmt(ts)
	} else if ts.Peek("FUNC") {
		return parseFuncDecl(ts)
	} else {
		return parseExprStmt(ts)
	}
}

// Parse different kinds of statements

func parsePrintStmt(ts *TokenStream) Node {
	ts.Expect("PRINT")
	expr := parseExpression(ts)
	ts.Expect(";")
	return ts.Builder().Build("PrintStmt", expr)
}

func parseConstDecl(ts *TokenStream) Node {
	ts.Expect("CONST")
	name := ts.Expect("ID").Value
	var typ string
	if tok := ts.Accept("ID"); tok != nil {
		typ = tok.Value
	}
	ts.Expect("=")
	value := parseExpression(ts)
	ts.Expect(";")
	return ts.Builder().Build("ConstDecl", name, typ, value)
}

func parseVarDecl(ts *TokenStream) Node {
	ts.Expect("VAR")
	name := ts.Expect("ID").Value
	var typ string
	if tok := ts.Accept("ID"); tok != nil {
		typ = tok.Value
	}
	var value Node
	if ts.Accept("=") {
		value = parseExpression(ts)
	}
	ts.Expect(";")
	return ts.Builder().Build("VarDecl", name, typ, value)
}

func parseExprStmt(ts *TokenStream) Node {
	expr := parseExpression(ts)
	ts.Expect(";")
	return ts.Builder().Build("ExprStmt", expr)
}

func parseIfStmt(ts *TokenStream) Node {
	ts.Expect("IF")
	test := parseExpression(ts)
	ts.Expect("{")
	consequence := parseStatements(ts)
	ts.Expect("}")
	var alternative []Node
	if ts.Accept("ELSE") {
		ts.Expect("{")
		alternative = parseStatements(ts)
		ts.Expect("}")
	}
	return ts.Builder().Build("IfStmt", test, consequence, alternative)
}

func parseWhileStmt(ts *TokenStream) Node {
	ts.Expect("WHILE")
	test := parseExpression(ts)
	ts.Expect("{")
	body := parseStatements(ts)
	ts.Expect("}")
	return ts.Builder().Build("WhileStmt", test, body)
}

func parseBreakStmt(ts *TokenStream) Node {
	ts.Expect("BREAK")
	ts.Expect(";")
	return ts.Builder().Build("BreakStmt")
}

func parseContinueStmt(ts *TokenStream) Node {
	ts.Expect("CONTINUE")
	ts.Expect(";")
	return ts.Builder().Build("ContinueStmt")
}

func parseReturnStmt(ts *TokenStream) Node {
	ts.Expect("RETURN")
	value := parseExpression(ts)
	ts.Expect(";")
	return ts.Builder().Build("ReturnStmt", value)
}

func parseFuncDecl(ts *TokenStream) Node {
	ts.Expect("FUNC")
	name := ts.Expect("ID").Value
	ts.Expect("(")
	var params []Node
	for !ts.Peek(")") {
		pname := ts.Expect("ID").Value
		ptype := ts.Expect("ID").Value
		param := ts.Builder().Build("Param", pname, ptype)
		params = append(params, param)
		if !ts.Peek(")") {
			ts.Expect(",")
		}
	}
	ts.Expect(")")
	retType := ts.Expect("ID").Value
	ts.Expect("{")
	body := parseStatements(ts)
	ts.Expect("}")
	return ts.Builder().Build("FuncDecl", name, params, retType, body)
}

func parseExpression(ts *TokenStream) Node {
	return parseAssignExpr(ts)
}

func parseAssignExpr(ts *TokenStream) Node {
	left := parseOrExpr(ts)
	for ts.Peek("=") {
		ts.Next()
		right := parseAssignExpr(ts)
		left = ts.Builder().Build("AssignExpr", left, right)
	}
	return left
}

func parseOrExpr(ts *TokenStream) Node {
	left := parseAndExpr(ts)
	for ts.Peek("||") {
		ts.Next()
		right := parseAndExpr(ts)
		left = ts.Builder().Build("OrExpr", left, right)
	}
	return left
}

func parseAndExpr(ts *TokenStream) Node {
	left := parseRelExpr(ts)
	for ts.Peek("&&") {
		ts.Next()
		right := parseRelExpr(ts)
		left = ts.Builder().Build("AndExpr", left, right)
	}
	return left
}

func parseRelExpr(ts *TokenStream) Node {
	left := parseAddExpr(ts)

	for ts.Peek("<", "<=", ">", ">=", "==", "!=") {
		op := ts.Next().Value
		right := parseAddExpr(ts)
		left = ts.Builder().Build("RelExpr", left, op, right)
	}

	return left
}

func parseAddExpr(ts *TokenStream) Node {
	left := parseMulExpr(ts)

	for ts.Peek("+", "-") {
		op := ts.Next().Value
		right := parseMulExpr(ts)
		left = ts.Builder().Build("AddExpr", left, op, right)
	}

	return left
}

func parseMulExpr(ts *TokenStream) Node {
	left := parseFactor(ts)

	for ts.Peek("*", "/") {
		op := ts.Next().Value
		right := parseFactor(ts)
		left = ts.Builder().Build("MulExpr", left, op, right)
	}

	return left
}

func parseFactor(ts *TokenStream) Node {
	// Parse expressions
	if ts.Peek("INTEGER") {
		return ts.Builder().Build("Integer", ts.Next().Value)
	} else if ts.Peek("FLOAT") {
		return ts.Builder().Build("Float", ts.Next().Value)
	} else if ts.Peek("TRUE", "FALSE") {
		return ts.Builder().Build("Bool", ts.Next().Value)
	} else if ts.Peek("CHAR") {
		return ts.Builder().Build("Char", ts.Next().Value)
	} else if ts.Peek("(") {
		ts.Next()
		expr := parseExpression(ts)
		ts.Expect(")")
		return ts.Builder().Build("ParenExpr", expr)
	} else if ts.Peek("{") {
		ts.Next()
		stmts := parseStatements(ts)
		ts.Expect("}")
		return ts.Builder().Build("CompoundExpr", stmts...)
	} else if ts.Peek("+", "-", "!") {
		op := ts.Next().Value
		operand := parseFactor(ts)
		return ts.Builder().Build("UnaryExpr", op, operand)
	} else if ts.Peek("ID") {
		// Either a variable or function call
		id := ts.Next().Value
		if ts.Peek("(") {
			ts.Next()
			args := parseArguments(ts)
			ts.Expect(")")
			return ts.Builder().Build("CallExpr", id, args...)
		} else {
			return ts.Builder().Build("Variable", id)
		}
	}

	panic(fmt.Sprintf("Unexpected token %v", ts.Lookahead))
}

func parseArguments(ts *TokenStream) []Node {
	var args []Node
	for !ts.Peek(")") {
		expr := parseExpression(ts)
		args = append(args, expr)
		if !ts.Peek(")") {
			ts.Expect(",")
		}
	}
	return args
}

//// Helper functions
//
//func newNode(typ string, children ...Node) *Node {
//	return &Node{typ, children}
//}
//
//type Node struct {
//	Type     string
//	Children []Node
//}
//
//func (n *Node) Pos() (int, int) {
//	// Unimplemented
//	return 0, 0
//}
//
//type Program struct {
//	Source string
//	AST    []Node
//}
//
//func (p *Program) Pos(node Node, startLine, startIndex, endIndex int) {
//	// Unimplemented
//}
//
//func tokenize(source string) chan *Token {
//	// Lexical analysis, unimplemented
//	return make(chan *Token)
//}

func Main() {
	// Test parsing
	source := `
  const PI = 3.14;
  var r number;
  print area(r); 
`
	program := &Program{source, nil}
	ParseProgram(program)
	fmt.Println(program.AST)
}
