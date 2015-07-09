package main

import (
	"fmt"
	"unicode/utf8"
	"bytes"
	"strings"
	"strconv"
	"errors"
)

type Token int

type Result struct {
	value string
	left_over string
}

type TokenValue struct {
	value string
	token Token
}

type TokenTree struct {
	TokenValue
	branches []*TokenTree
}

type TokenRoot struct {
	branches []*TokenTree
}

type Parser struct {
	token  Token
	action func(string) ([]*Result)
	next   *Parser
}

// Methods for Parser

func NewTokenTree(value string, token Token) *TokenTree {
	toktree := new(TokenTree)
	toktree.value = value
	toktree.token = token
	toktree.branches = make([]*TokenTree, 0)
	return toktree
}

func (self *TokenTree) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("'")
	buffer.WriteString(self.value)
	buffer.WriteString("'")
	buffer.WriteString("[")
	buffer.WriteString(strconv.Itoa(int(self.token)))
	buffer.WriteString("]")

	buffer.WriteString("(")
	for i, tree := range self.branches {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(tree.String())
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (self *TokenRoot) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("TokenRoot (")
	for i := range self.branches {
		if i != 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(self.branches[i].String())
	}
	buffer.WriteString(")")

	return buffer.String()
}

func (self *Parser) tokenize(results []*Result) []*TokenRoot {
	outputs := make([]*TokenRoot, len(results))
	for i, result := range results {
		outputs[i] = new(TokenRoot)
		outputs[i].branches = []*TokenTree{NewTokenTree(result.value, self.token)}
	}
	return outputs
}

func (self *Parser) retokenize(root_list []*TokenRoot, value string) []*TokenRoot {
	for i := range root_list {
		if len(value) > 0 { // Empty strings can be parsed, sure, but they shouldn't be added here.
			root_list[i].branches = append(root_list[i].branches, NewTokenTree(value, self.token))
		}
	}
	return root_list
}

func (self *Parser) BasicParse(str string) []*TokenRoot {
	return self.tokenize(self.action(str))
}

func (self *Parser) Parse(str string) []*TokenRoot {
	if self.next != nil {
		parent_results := self.action(str)
		outputs := make([]*TokenRoot, 0)
		for _, a_result := range parent_results {
			child_results := self.next.Parse(a_result.left_over)
			// I think this append will work, but I'm not entirely sure.
			outputs = append(outputs, self.retokenize(child_results, a_result.value)...)
		}
		return outputs
	}
	return self.BasicParse(str) // End of the line
}

func (self *Parser) Chain(other *Parser) *Parser {
	if self.next == nil {
		self.next = other
	} else {
		self.next.Chain(other)
	}
	return other // faster for chaining.
}

//////////////////////////////////////////////

// Actual Parsing

// Make a map to get strings from tokens

const (
	EOF = iota
	LEFT
	RIGHT
	SPACE
	NUMBER
)

var eof rune


// Wraps a string -> []string, []string function
// and adds the original string to the left_overs
func optionally(f func(string) []*Result) func(string) []*Result  {
  return func (s string) []*Result {
    results := f(s)
    return append(results, &Result{value: "", left_over: s})
  }
}

func matchString(string_a string) func(string) []*Result {
  return satisfyString(func (string_b string) (string, error) {
    if strings.HasPrefix(string_b, string_a) {
      return string_a, nil
    }
		err := "String " + string_a + " did not start with " + string_b + "."
    return "", errors.New(err)
  })
}

func matchNothing() func(string) []*Result {
	return satisfyString(func (string_b string) (string, error) {
		return "", errors.New("Never Satisfied. Never Graduate")
	})
}

func satisfyString(f func(string) (string, error)) func(string) []*Result {
	return func(str string) []*Result {
		parsed_string, err := f(str)
		if err == nil {
			// Satisfied
			return []*Result{&Result{left_over: str[len([]rune(parsed_string)):], value: parsed_string}}
		}
    return []*Result{} // Did not parse; return empty Result list
	}
}

func satisfyRune(f func(rune) bool) func(string) []*Result {
	return satisfyString(func(s string) (string, error) {
    c, _ := utf8.DecodeRuneInString(s)
    if f(c) {
      return string(c), nil
    } else {
      return "", errors.New("Rune was not satisfied")
    }
	})
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func isChar(a rune) func(rune) bool {
	return func(c rune) bool {
		return c == a
	}
}

////
// Debugging

func debug_function() func() {
	a := 'a'
	return func() {
		fmt.Printf("%c\n", a)
		a++
	}
}

////

var a func()

func init() {
	eof = rune(0)
	a = debug_function()
}

func main() {
	var p *Parser

	p = new(Parser)
	p.token = SPACE
	p.action = optionally(satisfyRune(isSpace))

	s := new(Parser)
	s.token = SPACE
	s.action = optionally(satisfyRune(isSpace))


	q := new(Parser)
	q.token = LEFT
	q.action = satisfyRune(isChar('c'))

  r := new(Parser)
  r.token = 5 // Just for testing purposes
  r.action = matchString("Hello") // MatchString is no longer greedy!

	t := new(Parser)
	t.token = 0 // eof
	t.action = matchNothing()

	q.Chain(p).Chain(s).Chain(r)

	root_list := q.Parse("c Hello World")

	for i := range root_list {
		fmt.Println(root_list[i])
	}

}
