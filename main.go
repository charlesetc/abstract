package main

import (
	"fmt"
	"unicode/utf8"
)

type Token int

// This is the real token
// It has both the type and the
// text associated with it
type Result struct {
	value string
	token Token
}

type Parser struct {
	token  Token
	action func(string) ([]string, []string) "outputs a list of things left to parse and a list of results"
	next   *Parser
}

const (
	EOF = iota
	LEFT
	RIGHT
	SPACE
	NUMBER
)

var eof rune

// Methods for Parser

func (self *Parser) parse(str string) []*Result {
	strings, results := self.action(str)
	outputs := make([]*Result, len(results))

	a()

	var i int
	for _, a_string := range results {
		if len(a_string) > 0 {
			outputs[i] = &Result{value: a_string, token: self.token}
			i++
		}
	}

	for _, a_string := range strings {
		if self.next != nil {
			results := self.next.parse(a_string)
			outputs = append(outputs, results...)
		}
	}

	return outputs
}

func (self *Parser) chain(other *Parser) {
	self.next = other // or should it take a pointer as an argument?
}

//////////////////////////////////////////////

// Actual Parsing


// Wraps a string -> []string, []string function
// and adds the original string to the left_overs
func optionally(f func(string) ([]string, []string)) func(string) ([]string, []string)  {
  return func (s string) ([]string, []string) {
    left_over, output := f(s)
    return append(left_over, s), output
  }
}

func matchString(string_a string) func(string) ([]string, []string) {
  return satisfyString(func (string_b string) string {
    if string_a == string_b {
      return string_a
    }
    return "" // The equivalent of 'false' that I'm using
  })
}

func satisfyString(f func(string) string) func(string) ([]string, []string) {
	return func(str string) ([]string, []string) {
    var left_over, output []string
		result := f(str)
    length := len([]byte(result))
		if length > 0 { // satisfied
			output = []string{result}
			left_over = []string{str[length:]}
		} else { // don't parse
			output = []string{}
			left_over = []string{str}
		}
    return left_over, output
	}
}

func satisfyRune(f func(rune) bool) func(string) ([]string, []string) {
	return satisfyString(func(s string) string {
    c, _ := utf8.DecodeRuneInString(s)
    if f(c) {
      return string(c)
    } else {
      return ""
    }
	})
}

// This version doesn't reuse code:

// func satisfyRune(f func(rune) bool) func(string) ([]string, []string) {
// 	return func(str string) (left_over []string, output []string) {
// 		c, length := utf8.DecodeRune([]byte(str))
// 		if f(c) { // satisfied
// 			output := []string{string(c)}
// 			left_over := []string{str[length:]}
// 		} else { // don't parse
// 			output := []string{}
// 			left_over := []string{str}
// 		}
// 		return
// 	}
// }

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func isChar(a rune) func(rune) bool {
	return func(c rune) bool {
		return c == a
	}
}

func makecopy(str string) string {
	out := str
	return out
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

	q := new(Parser)
	q.token = LEFT
	q.action = satisfyRune(isChar('c'))

  r := new(Parser)
  r.token = 5 // Just for testing purposes
  r.action = matchString("Hello")



	q.chain(r) // REVIEW: Change this.
  r.chain(p)

	outputs := q.parse("cHello")

	for _, val := range outputs {
		fmt.Println(*val)
	}

}
