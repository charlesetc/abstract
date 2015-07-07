
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
  token Token
  action func (string) ([]string, []string) "outputs a list of things left to parse and a list of results"
  next *Parser
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

func (self Parser) parse(str string) []*Result {
  strings, results := self.action(str)
  outputs := make([]*Result, len(results))

  var i int
  for _, a_string := range results {
    if len(a_string) > 0 {
      outputs[i] = &Result{value: a_string, token: self.token }
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

func (self Parser) chain(other Parser) {
  self.next = &other // or should it take a pointer as an argument?
}

//////////////////////////////////////////////

// Actual Parsing

func isSpace(c rune) bool {
  return c == ' ' || c == '\t' || c == '\n'
}

func makecopy(str string) string {
  out := str
  return out
}

////
// Debugging

func debug_function() func() {
  a := 'a'
  return func () {
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

  p.action = func (str string) ([]string, []string) {
    output := make([]string, 1)

    c, length := utf8.DecodeRune([]byte(str))

    if isSpace(c) {
      output[0] = string(c) // Ideally parse more...
    }

    left_over := []string{str[length:]} // This operation is not optional

    return left_over, output
  }

  a()

  outputs := p.parse("c Hello")

  a()

  for _, val := range outputs {
    fmt.Println(*val)
  }

}
