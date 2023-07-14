package model

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Program struct {
	Source     string
	Model      Node
	HaveErrors bool
	Db         map[int]Locator
}

func NewProgram(source string) *Program {
	return &Program{
		Source:     source,
		HaveErrors: false,
		Db:         make(map[int]Locator),
	}
}

func (p *Program) ErrorMessage(message string) {
	fmt.Println(message)
	p.HaveErrors = true
}

func ProgramFromFile(filename string) (*Program, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewProgram(string(content)), nil
}

func (p *Program) RecordPosition(node Node, lineno, start, end int) {
	// when calling Id it should in Map
	p.Db[node.Id()] = NewLocator(p.Source, lineno, start, end)
}

func (p *Program) Location(node Node) Locator {
	return p.Db[node.Id()]
}

type Locator struct {
	SourceCode string
	Lineno     int
	Start      int
	End        int
}

func NewLocator(sourceCode string, lineno, start, end int) Locator {
	return Locator{
		SourceCode: sourceCode,
		Lineno:     lineno,
		Start:      start,
		End:        end,
	}
}

func (l Locator) Source() string {
	return l.SourceCode[l.Start:l.End]
}

func (l Locator) LineContext(start, end int) string {
	if start == 0 {
		start = l.Start
	}
	if end == 0 {
		end = l.End
	}
	s := start
	for s >= 0 && l.SourceCode[s] != '\n' {
		s--
	}
	if l.SourceCode[s] == '\n' {
		s++
	}

	e := end - 1
	for e < len(l.SourceCode) && l.SourceCode[e] != '\n' {
		e++
	}
	if l.SourceCode[e] != '\n' {
		e++
	}

	return l.SourceCode[s:e] + "\n" + strings.Repeat(" ", start-s) + strings.Repeat("^", end-start)
}
