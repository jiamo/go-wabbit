package parser

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"wabbit-go/common"
	"wabbit-go/model"
	"wabbit-go/tokenize"
)

func init() {
	// how to control the log level all the package
	log.SetLevel(log.DebugLevel)
}

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
	current    tokenize.Token // save
	lastIndex  int
	haveErrors bool
}

func NewTokenStream(program *model.Program) (*TokenStream, error) {
	tokens, err := tokenize.Tokenize(program.Source)

	if err != nil {
		log.Errorf("tokenize err %v", err)
		return nil, err
	}
	ts := &TokenStream{
		program:    program,
		gen:        common.SliceToChannel(tokens),
		haveErrors: false,
	}
	ts.lookahead = <-ts.gen
	log.Debugf("beginning %v", ts.lookahead)
	return ts, nil
}

func (ts *TokenStream) Peek(types ...string) *tokenize.Token {
	log.Debugf("Peeking %v while ts.lookahead.Type is %v ", types, ts.lookahead.Type)
	for _, t := range types {
		if ts.lookahead.Type == t {
			ts.current = ts.lookahead
			return &ts.current
		}
	}
	return nil
}

// the lookhead alreay change we can't just using the point
func (ts *TokenStream) Accept(types ...string) *tokenize.Token {
	tok := ts.Peek(types...)
	log.Debugf("Accpet.... %v while head is %v", types, tok)
	if tok != nil {
		ts.lastIndex = ts.lookahead.Index + len(ts.lookahead.Value)
		ts.lookahead = <-ts.gen // this while overwrite it
		log.Debugf("ts.lookahead is %v", ts.lookahead)
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

func (ts *TokenStream) Builder() func(func(constructFunc) model.Node) model.Node {
	startLine := ts.lookahead.Lineno
	startIndex := ts.lookahead.Index

	construct := func(node model.Node) model.Node {
		//print("nodetype=", nodetype)
		//node := nodeInit(args)
		// we got node don't init any more
		ts.program.RecordPosition(node, startLine, startIndex, ts.lastIndex) // ts.lookahead.index
		return node
	}

	return func(do func(constructFunc) model.Node) model.Node {
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
		return do(construct)
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

// Statements is struct
// node statement expression is interface..
func parseStatements(ts *TokenStream) *model.Statements {
	log.Debugf("parseStatements")
	statements := []model.Statement{} //

	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		for ts.Peek("RBRACE", "EOF") == nil {
			// python not None mean none is None
			log.Debugf("before parseStatement")
			statement := parseStatement(ts) //
			statements = append(statements, statement)
		}
		return new(&model.Statements{statements})
	})
	return node.(*model.Statements)
}

func parseStatement(ts *TokenStream) model.Statement {
	// Parse different types of statements
	if ts.Peek("PRINT") != nil {
		return parsePrintStmt(ts)
	} else if ts.Peek("CONST") != nil {
		return parseConstDecl(ts)
	} else if ts.Peek("VAR") != nil {
		return parseVarDecl(ts)
	} else if ts.Peek("RETURN") != nil {
		return parseReturnStmt(ts)
	} else if ts.Peek("IF") != nil {
		return parseIfStmt(ts)
	} else if ts.Peek("WHILE") != nil {
		return parseWhileStmt(ts)
	} else if ts.Peek("BREAK") != nil {
		return parseBreakStmt(ts)
	} else if ts.Peek("CONTINUE") != nil {
		return parseContinueStmt(ts)
	} else if ts.Peek("FUNC") != nil {
		return parseFuncDecl(ts)
	} else {
		return parseExprStmt(ts)
	}
}

// Parse different kinds of statements

func parsePrintStmt(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	log.Debugf("parsePrintStmt")
	node := builder(func(new constructFunc) model.Node {
		ts.Expect("PRINT")
		expr := parseExpression(ts)
		log.Debugf("parsePrintStmt expr is %v", expr)
		ts.Expect("SEMI")
		return new(&model.PrintStatement{expr}) // TODO need to detect Do we need &
	})
	return node.(model.Statement)
}

func parseConstDecl(ts *TokenStream) model.Statement {
	builder := ts.Builder()

	node := builder(func(new constructFunc) model.Node {
		ts.Expect("CONST")
		nameTk, err := ts.Expect("ID")
		if err != nil {

		}
		name := nameTk.Value
		var typ model.Type
		if tok := ts.Accept("ID"); tok != nil {
			typ = &model.NameType{tok.Value}
		}
		ts.Expect("=")
		value := parseExpression(ts)
		ts.Expect("SEMI")
		return new(&model.ConstDeclaration{model.Name{name}, typ, value})
	})

	return node.(model.Statement)
}

func parseVarDecl(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		ts.Expect("VAR")
		name, err := ts.Expect("ID")
		if err != nil {
			log.Debugf("err parseVarDecl")
		}
		var typ model.Type
		if tok := ts.Accept("ID"); tok != nil {
			typ = &model.NameType{tok.Value}
		}
		var value model.Expression
		if tok := ts.Accept("="); tok != nil {
			value = parseExpression(ts)
		}
		ts.Expect("SEMI")
		return new(&model.VarDeclaration{model.Name{name.Value}, typ, value})
	})
	return node.(model.Statement)
}

func parseExprStmt(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		expr := parseExpression(ts)
		ts.Expect("SEMI")
		return new(&model.ExpressionAsStatement{expr})
	})
	return node.(model.Statement)
}

func parseIfStmt(ts *TokenStream) model.Statement {

	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		ts.Expect("IF")
		test := parseExpression(ts)
		ts.Expect("{")
		consequence := parseStatements(ts)
		ts.Expect("}")
		var alternative *model.Statements
		if ts.Accept("ELSE") != nil {
			ts.Expect("{")
			alternative = parseStatements(ts)
			ts.Expect("}")
		}
		// how strange the same struct using different type
		return new(&model.IfStatement{test, *consequence, alternative})
	})
	return node.(model.Statement)

}

func parseWhileStmt(ts *TokenStream) model.Statement {

	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {

		ts.Expect("WHILE")
		test := parseExpression(ts)
		ts.Expect("{")
		body := parseStatements(ts)
		ts.Expect("}")
		// how strange the same struct using different type
		return new(&model.WhileStatement{test, *body})
	})
	return node.(model.Statement)
}

func parseBreakStmt(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {

		ts.Expect("BREAK")
		ts.Expect("SEMI")
		return new(&model.BreakStatement{})
	})
	return node.(model.Statement)
}

func parseContinueStmt(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {

		ts.Expect("CONTINUE")
		ts.Expect("SEMI")
		return new(&model.ContinueStatement{})
	})
	return node.(model.Statement)

}

func parseReturnStmt(ts *TokenStream) model.Statement {
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {

		ts.Expect("RETURN")
		value := parseExpression(ts)
		ts.Expect("SEMI")
		return new(&model.ReturnStatement{value})
	})
	return node.(model.Statement)
}

func parseFuncDecl(ts *TokenStream) model.Statement {

	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {

		ts.Expect("FUNC")
		nameToken, err := ts.Expect("ID")
		if err != nil {
			log.Debugf("")
		}
		name := model.Name{nameToken.Value}
		ts.Expect("(")
		var params []model.Parameter
		for ts.Peek(")") != nil {

			paramsBuilder := ts.Builder()
			node := paramsBuilder(func(newp constructFunc) model.Node {
				pnameToken, err := ts.Expect("ID")
				if err != nil {
					log.Debugf("")
				}
				pname := model.Name{pnameToken.Value}
				ptypeNode, err := ts.Expect("ID")
				if err != nil {
					log.Debugf("")
				}
				ptype := model.NameType{ptypeNode.Value}
				//too complicate to use nest context TODO we may should do at the end
				//param := ts.Builder().Build("Param", pname, ptype)
				//param := model.Parameter{pname, &ptype}
				return newp(&model.Parameter{pname, &ptype})
			})
			param := node.(*model.Parameter)
			params = append(params, *param)
			if ts.Peek(")") != nil {
				ts.Expect(",")
			}
		}
		ts.Expect(")")
		retTypeNode, err := ts.Expect("ID")
		if err != nil {
			log.Debugf("")
		}
		retType := model.NameType{retTypeNode.Value}
		ts.Expect("{")
		body := parseStatements(ts)
		ts.Expect("}")
		return new(&model.FunctionDeclaration{name, params, &retType, *body})
	})
	return node.(model.Statement)

	//return ts.Builder().Build("FuncDecl", name, params, retType, body)
}

func parseExpression(ts *TokenStream) model.Expression {
	return parseAssignExpr(ts)
}

func parseAssignExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseAssignExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseOrExpr(ts)
		for ts.Accept("ASSIGN") != nil {
			right := parseAssignExpr(ts)
			left = new(&model.Assignment{left, right}).(model.Expression)
		}
		return left
	})
	return node.(model.Expression)

}

func parseOrExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseOrExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseAndExpr(ts)
		for ts.Accept("LOR") != nil {
			right := parseAndExpr(ts)
			left = new(&model.LogOr{left, right}).(model.Expression)
		}
		return left
	})
	return node.(model.Expression)

}

func parseAndExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseAndExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseRelExpr(ts)
		for ts.Accept("LAND") != nil {
			right := parseRelExpr(ts)
			left = new(&model.LogAnd{left, right}).(model.Expression)
		}
		return left
	})
	return node.(model.Expression)

}

func parseRelExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseRelExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseAddExpr(ts)

		for {
			tok := ts.Accept("<", "<=", ">", ">=", "==", "!=")
			if tok == nil {
				break
			}
			op := tok.Value
			right := parseAddExpr(ts)
			if op == "<" {
				left = new(&model.Lt{left, right}).(model.Expression)
			} else if op == "<=" {
				left = new(&model.Le{left, right}).(model.Expression)
			} else if op == ">" {
				left = new(&model.Gt{left, right}).(model.Expression)
			} else if op == ">=" {
				left = new(&model.Ge{left, right}).(model.Expression)
			} else if op == "==" {
				left = new(&model.Eq{left, right}).(model.Expression)
			} else if op == "!=" {
				left = new(&model.Ne{left, right}).(model.Expression)
			}
		}

		return left
	})
	return node.(model.Expression)

}

func parseAddExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseAddExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseMulExpr(ts)

		for {
			tok := ts.Accept("+", "-")
			if tok == nil {
				break
			}
			op := tok.Value
			right := parseMulExpr(ts)
			if op == "+" {
				left = new(&model.Add{left, right}).(model.Expression)
			} else if op == "-" {
				left = new(&model.Sub{left, right}).(model.Expression)
			}
		}

		return left
	})
	return node.(model.Expression)

}

func parseMulExpr(ts *TokenStream) model.Expression {
	log.Debugf("parseMulExpr")
	builder := ts.Builder()
	node := builder(func(new constructFunc) model.Node {
		left := parseFactor(ts)

		for {
			tok := ts.Accept("*", "/")
			if tok == nil {
				break
			}
			op := tok.Value
			right := parseFactor(ts)
			if op == "*" {
				left = new(&model.Mul{left, right}).(model.Expression)
			} else if op == "/" {
				left = new(&model.Div{left, right}).(model.Expression)
			}
		}

		return left
	})
	return node.(model.Expression)
}

func parseFactor(ts *TokenStream) model.Expression {
	builder := ts.Builder()
	log.Debugf("parseFactor")
	node := builder(func(new constructFunc) model.Node {
		// Parse expressions
		if tok := ts.Accept("INTEGER"); tok != nil {
			num, err := strconv.Atoi(tok.Value)
			if err != nil {
				log.Errorf("wrong int %v ", tok.Value)
				return nil
			}
			return new(&model.Integer{num})
		} else if tok := ts.Accept("FLOAT"); tok != nil {
			num, err := strconv.ParseFloat(tok.Value, 64)
			if err != nil {
				log.Errorf("wrong float %v", tok.Value)
				return nil
			}
			return new(&model.Float{num})
		} else if tok := ts.Accept("TRUE", "FALSE"); tok != nil {
			return new(&model.NameBool{tok.Value})
		} else if ts.Accept("CHAR"); tok != nil {
			return new(&model.Character{tok.Value})
		} else if tok := ts.Accept("LPAREN"); tok != nil {
			expr := parseExpression(ts)
			ts.Expect("RPAREN")
			return new(&model.Grouping{expr})
		} else if tok := ts.Accept("LBRACE"); tok != nil {
			stmts := parseStatements(ts)
			ts.Expect("RBRACE")
			return new(&model.CompoundExpression{stmts.Statements})
		} else if tok := ts.Accept("+", "-", "!"); tok != nil {
			operand := parseFactor(ts)
			if tok.Value == "+" {
				return new(&model.Pos{operand})
			} else if tok.Value == "-" {
				return new(&model.Neg{operand})
			} else if tok.Value == "!" {
				return new(&model.Not{operand})
			}
		} else if tok := ts.Accept("ID"); tok != nil {
			// Either a variable or function call
			// do we need a parse location.... ?
			if ts.Accept("LPAREN") != nil {
				args := parseArguments(ts)
				ts.Expect("RPAREN")
				return new(&model.FunctionApplication{&model.Name{tok.Value}, args})
			} else {
				return new(&model.Name{tok.Value})
			}
		} else {
			panic(fmt.Sprintf("Unexpected token %v", ts.lookahead))
		}
		return nil
	})
	return node.(model.Expression)
	//panic(fmt.Sprintf("Unexpected token %v", ts.Lookahead))
}

func parseArguments(ts *TokenStream) []model.Expression {

	var args []model.Expression
	for ts.Peek(")") != nil {
		expr := parseExpression(ts)
		args = append(args, expr)
		if ts.Peek(")") != nil {
			ts.Expect(",")
		}
	}
	return args

}

// main function to test on input files
func HandleFile(filename string) (*model.Program, error) {
	prog, err := model.ProgramFromFile(filename)
	if err != nil {
		return nil, err
	}
	err = ParseProgram(prog)
	return prog, err
}
