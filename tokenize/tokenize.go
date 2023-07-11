package tokenize

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Token structure
type Token struct {
	Type   string
	Value  string
	Lineno int
	Index  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %s, %d, %d)", t.Type, t.Value, t.Lineno, t.Index)
}

var literals = map[string]string{
	"+":  "PLUS",
	"-":  "MINUS",
	"*":  "TIMES",
	"/":  "DIVIDE",
	"<":  "LT",
	"<=": "LE",
	">":  "GT",
	">=": "GE",
	"=":  "ASSIGN",
	"==": "EQ",
	"!=": "NE",
	"&&": "LAND",
	"||": "LOR",
	"!":  "LNOT",
	";":  "SEMI",
	"(":  "LPAREN",
	")":  "RPAREN",
	"{":  "LBRACE",
	"}":  "RBRACE",
	",":  "COMMA",
}

var keywords = map[string]bool{"print": true, "if": true, "else": true, "var": true, "const": true, "func": true, "while": true, "break": true, "continue": true, "return": true, "true": true, "false": true}

// tokenize function
func Tokenize(text string) ([]Token, error) {
	var err error
	tokens := []Token{}
	n := 0
	lineno := 1
	size := len(text)
	error_array := []string{}
	for n < size {
		if text[n] == ' ' || text[n] == '\t' {
			n++
			continue
		} else if text[n] == '\n' {
			n++
			lineno++
			continue
		} else if n+1 < size && text[n:n+2] == "/*" {
			end := strings.Index(text[n:], "*/")
			if end < 0 {
				error_array = append(error_array, fmt.Sprintf("on line %d: Unterminated comment", lineno))
				n = size
				continue
			} else {
				lineno += strings.Count(text[n:end], "\n")
				n = n + end + 2
				continue
			}
		} else if n+1 < size && text[n:n+2] == "//" {
			end := strings.Index(text[n:], "\n")
			if end < 0 {
				break
			}
			lineno++
			n = n + end + 1
			continue
		} else if unicode.IsDigit(rune(text[n])) {
			start := n
			for n < size && unicode.IsDigit(rune(text[n])) {
				n++
			}
			if n < size && text[n] == '.' {
				n++
				for n < size && unicode.IsDigit(rune(text[n])) {
					n++
				}
				tokens = append(tokens, Token{Type: "FLOAT", Value: text[start:n], Lineno: lineno, Index: start})
			} else {
				tokens = append(tokens, Token{Type: "INTEGER", Value: text[start:n], Lineno: lineno, Index: start})
			}
			continue
		} else if unicode.IsLetter(rune(text[n])) || text[n] == '_' {
			start := n
			for n < size && (unicode.IsLetter(rune(text[n])) || text[n] == '_') {
				n++
			}
			tokType := text[start:n]
			if keywords[tokType] {
				tokens = append(tokens, Token{Type: strings.ToUpper(tokType), Value: tokType, Lineno: lineno, Index: start})
			} else {
				tokens = append(tokens, Token{Type: "ID", Value: tokType, Lineno: lineno, Index: start})
			}
			continue
		} else if text[n] == '\'' {
			start := n
			n++
			for n < size && text[n] != '\'' {
				if text[n] == '\\' {
					n++
				}
				n++
			}
			if n >= size {
				error_array = append(error_array, fmt.Sprintf("on line %d: Unterminated character constant", lineno))
			}
			tokens = append(tokens, Token{Type: "CHAR", Value: text[start : n+1], Lineno: lineno, Index: start})
			n++
			continue
		} else if n+1 < size && literals[text[n:n+2]] != "" {
			tokens = append(tokens, Token{Type: literals[text[n:n+2]], Value: text[n : n+2], Lineno: lineno, Index: n})
			n += 2
			continue
		} else if literals[string(text[n])] != "" {
			tokens = append(tokens, Token{Type: literals[string(text[n])], Value: string(text[n]), Lineno: lineno, Index: n})
			n++
			continue
		} else {
			error_array = append(error_array, fmt.Sprintf("on line %d: Illegal character %s", lineno, string(text[n])))
			n++
		}
	}
	tokens = append(tokens, Token{Type: "EOF", Value: "EOF", Lineno: lineno, Index: n})
	if len(error_array) > 0 {
		err = fmt.Errorf("error number %d\n error details :\n %s", len(error_array), strings.Join(error_array, "\n"))
	}
	return tokens, err
}

// main function to test on input files
func HandleFile(filename string) ([]Token, error) {
	file, err := os.Open(filename) // change filename to your file
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text string

	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}

	return Tokenize(text)

}
