import pathlib
import pytest
from wabbit.tokenize import main

right_program_path = pathlib.Path(__file__).resolve().parent / "Programs"
right_files = right_program_path.glob("*.wb")
wrong_program_path = pathlib.Path(__file__).resolve().parent / "ErrorLex"
wrong_files = wrong_program_path.glob("*.wb")

def test_right_program():
    # print(right_program_path)
    # print(list(right_files))
    for right_file in right_files:
        main(right_file)


def test_error_program():
    # print(right_program_path)
    # print(list(right_files))
    for file in wrong_files:
        with pytest.raises(Exception):
            main(file)



from wabbit.tokenize import tokenize
from wabbit.program import Program

def test_symbols():
    tokens = list(tokenize("+ - * / < > <= >= == != = && || , ; ( ) { } !"))
    toktypes = [tok.type for tok in tokens]
    assert toktypes == ['PLUS', 'MINUS', 'TIMES', 'DIVIDE',
                        'LT', 'GT', 'LE', 'GE', 'EQ', 'NE', 'ASSIGN',
                        'LAND', 'LOR', 'COMMA', 'SEMI',
                        'LPAREN', 'RPAREN', 'LBRACE', 'RBRACE', 'LNOT', 'EOF'], toktypes

def test_numbers():
    tokens = list(tokenize("123 123.45"))
    toktypes = [tok.type for tok in tokens]
    tokvalues = [tok.value for tok in tokens]
    assert toktypes ==  ['INTEGER', 'FLOAT', 'EOF'], toktypes
    assert tokvalues == ['123', '123.45', 'EOF'],  tokvalues

def test_keywords():
    tokens = list(tokenize("if else while var const break continue print func return true false"))
    toktypes = [tok.type for tok in tokens]
    assert toktypes == ['IF', 'ELSE', 'WHILE', 'VAR', 'CONST', 'BREAK', 'CONTINUE',
                        'PRINT', 'FUNC', 'RETURN', 'TRUE', 'FALSE', 'EOF']

if __name__ == '__main__':
    test_symbols()
    test_numbers()
    test_keywords()