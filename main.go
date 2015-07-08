package main

import (
	"fmt"
	"unicode/utf8"
	"strings"
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

type Parser struct {
	token  Token
	action func(string) ([]*Result)
	next   *Parser
}

// Methods for Parser


func (self *Parser) tokenize(results []*Result) [][]*TokenValue {
	outputs := make([][]*TokenValue, len(results))
	for i, result := range results {
		outputs[i] = make([]*TokenValue, 1)
		fmt.Println("***")
		fmt.Println(result)
		fmt.Println("***")
		outputs[i][0] = &TokenValue{value: result.value, token: self.token}
	}
	return outputs
}

func (self *Parser) retokenize(tokvals [][]*TokenValue, value string) [][]*TokenValue {
	for i := range tokvals {
		if len(value) > 0 { // Not Empty
			tokvals[i] = append(tokvals[i], &TokenValue{value: value, token: self.token})
		}
	}
	return tokvals
}

func (self *Parser) parse(str string) [][]*TokenValue {
	parent_results := self.action(str)
	outputs := make([][]*TokenValue, 0)
	if self.next != nil { // Pass results on to child to try parsing
		for _, a_result := range parent_results {
			child_results := self.next.parse(a_result.left_over)
			// I think this append will work, but I'm not entirely sure.
			outputs = append(outputs, self.retokenize(child_results, a_result.value)...)
		}
	} else { // Nothing else to parse...
		return self.tokenize(parent_results) // This is an essential line.
	}
	return outputs
}

func (self *Parser) chain(other *Parser) *Parser {
	if self.next == nil {
		self.next = other
	} else {
		self.next.chain(other)
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
  return satisfyString(func (string_b string) string {
    if strings.HasPrefix(string_b, string_a) {
      return string_a
    }
    return "" // The equivalent of 'false', as in 'no parse'
  })
}

func matchNothing() func(string) []*Result {
	return satisfyString(func (string_b string) string {
		return ""
	})
}

func satisfyString(f func(string) string) func(string) []*Result {
	return func(str string) []*Result {
		parsed_string := f(str)
    length := len([]rune(parsed_string))
		if length > 0 { // satisfied
			return []*Result{&Result{left_over: str[length:], value: parsed_string}}
		} // empty parse
    return []*Result{&Result{left_over: str, value: ""}}
	}
}

func satisfyRune(f func(rune) bool) func(string) []*Result {
	return satisfyString(func(s string) string {
    c, _ := utf8.DecodeRuneInString(s)
    if f(c) {
      return string(c)
    } else {
      return ""
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

	q.chain(p).chain(s).chain(r).chain(t)

	outputs := q.parse("c  Hello World")

	for i := range outputs {
		fmt.Println("----")
		for j := range outputs[i] {
			fmt.Println(outputs[i][j])
		}
	}

}
