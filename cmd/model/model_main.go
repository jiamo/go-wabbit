package main

import "wabbit-go/model"

func main() {
	modelA := &model.BinOpWithOp{
		"+",
		&model.Integer{2},
		&model.BinOpWithOp{
			"*",
			&model.Integer{2},
			&model.Integer{4}}}

	println(model.NodeAsSource(modelA, model.NewContext()))
	println(model.NodeAsSource(
		&model.PrintStatement{&model.Integer{2}},
		model.NewContext()))

	println(model.NodeAsSource(
		&model.Statements{
			[]model.Statement{
				&model.ConstDeclaration{model.Name{"pi"}, nil, &model.Float{3.14159}},
				&model.ConstDeclaration{model.Name{"tau"}, nil, &model.Mul{&model.Float{2.0}, &model.Name{"pi"}}},
				&model.VarDeclaration{model.Name{"radius"}, nil, &model.Float{4.0}},
				&model.VarDeclaration{model.Name{"perimeter"}, &model.NameType{"float"}, nil},
				&model.ExpressionAsStatement{&model.Assignment{&model.Name{"perimeter"}, &model.Mul{&model.Name{"tau"}, &model.Name{"radius"}}}},
				&model.PrintStatement{&model.Name{"perimeter"}},
			},
		},
		model.NewContext(),
	))

	println(model.NodeAsSource(
		&model.PrintStatement{&model.Eq{&model.Integer{2}, &model.Integer{2}}},
		model.NewContext()))

	println(model.NodeAsSource(
		&model.PrintStatement{
			&model.LogOr{&model.NameBool{"true"},
				&model.Grouping{
					&model.Eq{
						&model.Div{&model.Integer{1}, &model.Integer{2}},
						&model.Integer{2},
					}}}},
		model.NewContext()))

	println(model.NodeAsSource(
		&model.Statements{
			[]model.Statement{
				&model.ConstDeclaration{model.Name{"a"}, &model.NameType{"int"}, &model.Integer{2}},
				&model.ConstDeclaration{model.Name{"b"}, &model.NameType{"int"}, &model.Integer{3}},
				&model.IfStatement{
					&model.Lt{&model.Name{"a"}, &model.Name{"b"}},
					model.Statements{[]model.Statement{
						&model.ExpressionAsStatement{&model.Assignment{&model.Name{"minvar"}, &model.Name{"a"}}},
					}},
					&model.Statements{[]model.Statement{
						&model.ExpressionAsStatement{&model.Assignment{&model.Name{"minvar"}, &model.Name{"b"}}},
					}},
				},
			},
		},
		model.NewContext(),
	))

	p3 := &model.Statements{
		[]model.Statement{
			&model.VarDeclaration{model.Name{"n"}, &model.IntegerType{}, &model.Integer{'1'}},
			&model.WhileStatement{
				&model.NameBool{"true"},
				model.Statements{[]model.Statement{
					&model.IfStatement{
						&model.Eq{&model.Name{"n"}, &model.Integer{2}},
						model.Statements{
							[]model.Statement{
								&model.PrintStatement{&model.Name{"n"}},
								&model.BreakStatement{},
							}},
						&model.Statements{
							[]model.Statement{
								&model.ExpressionAsStatement{&model.Assignment{&model.Name{"n"}, &model.Add{&model.Name{"n"}, &model.Integer{'1'}}}},
								&model.ContinueStatement{},
							}},
					},
					&model.ExpressionAsStatement{&model.Assignment{&model.Name{"n"}, &model.Sub{&model.Name{"n"}, &model.Integer{1}}}}},
				}},
			&model.PrintStatement{&model.Name{"n"}}},
	}
	println(model.NodeAsSource(p3, model.NewContext()))

	p4 := &model.Statements{
		[]model.Statement{
			&model.VarDeclaration{model.Name{"x"}, nil, &model.Integer{37}},
			&model.VarDeclaration{model.Name{"y"}, nil, &model.Integer{42}},
			&model.ExpressionAsStatement{
				&model.Assignment{
					&model.Name{"x"},
					&model.CompoundExpression{
						[]model.Statement{
							&model.VarDeclaration{model.Name{"t"}, nil, &model.Name{"y"}},
							&model.ExpressionAsStatement{&model.Assignment{&model.Name{"y"}, &model.Name{"x"}}},
							&model.ExpressionAsStatement{&model.Name{"t"}},
						},
					},
				},
			},
		},
	}
	println(model.NodeAsSource(p4, model.NewContext()))

	p5 := &model.Statements{
		[]model.Statement{
			&model.FunctionDeclaration{
				Name: model.Name{"add"},
				Parameters: []model.Parameter{
					model.Parameter{
						Name: model.Name{"x"},
						Type: &model.IntegerType{},
					},
					model.Parameter{
						Name: model.Name{"y"},
						Type: &model.IntegerType{},
					},
				},
				ReturnType: &model.IntegerType{},
				Body: model.Statements{
					[]model.Statement{
						&model.ReturnStatement{
							Value: &model.Add{
								Left:  &model.Name{"x"},
								Right: &model.Name{"y"},
							},
						},
					},
				},
			},
			&model.VarDeclaration{
				Name: model.Name{"result"},
				Type: nil,
				Value: &model.FunctionApplication{
					&model.Name{"add"},
					[]model.Expression{
						&model.Integer{2},
						&model.Integer{3},
					},
				},
			},
			&model.PrintStatement{
				Value: &model.Name{"result"},
			},
		},
	}

	println(model.NodeAsSource(p5, model.NewContext()))

}
