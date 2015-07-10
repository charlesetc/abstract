package main

import (
	"fmt"
	"unicode/utf8"
	"bytes"
	"strings"
	"errors"
	"strconv"
)

type Token int

const (
	OPERATOR Token = ((iota + 1) * -1)
	LEFT_SIDE
	RIGHT_SIDE
)

type Result struct {
	value string
	left_over string
}

type TokenValue struct {
	value string
	token Token
}

type TokenRoot struct {
	branches []*TokenTree
}


type TokenTree struct {
	TokenValue
	branches []*TokenTree
}


type ParserLike interface {
	CanParse(string) bool
	Parse(string) []*TokenRoot
	Chain(ParserLike) ParserLike
}

type Lexer struct {
	token  Token
	action func(string) ([]*Result)
	next   ParserLike
}


//// AbstractParsers

type AbstractParser struct {
	parse_action func(string, *AbstractParser) []*TokenRoot
	next ParserLike
}

// I think compilers should be transparent to other compilers
// when it comes to whether or not they'll work.
func (self *AbstractParser) CanParse(str string) bool {
	return self.next.CanParse(str) // Is this right?
}

func (self *AbstractParser) Parse(str string) []*TokenRoot {
	return self.parse_action(str, self)
}


func (self *AbstractParser) Chain(other ParserLike) ParserLike {
	if self.next == nil {
		self.next = other
	} else {
		self.next.Chain(other)
	}
	return other // faster for chaining.
}

func NewAbstractParser(f func (string, *AbstractParser) []*TokenRoot) *AbstractParser {
	c := new(AbstractParser)
	c.parse_action = f
	return c
}

func Optionally(parser ParserLike) *AbstractParser {
	return NewAbstractParser(func (str string, self *AbstractParser) []*TokenRoot {
		if parser.CanParse(str) {
			parser.Chain(self.next) // obvs
			return parser.Parse(str)
		}
		// Else can't parse
		return self.next.Parse(str)
	})
}

//// Compilers

// Compilers just deal with the tokens and the syntax tree.

func NewCompiler(f func ([]*TokenRoot) []*TokenRoot) *AbstractParser {
	return NewAbstractParser(func (str string, self *AbstractParser) []*TokenRoot {
		roots := self.next.Parse(str) // Basically just wraps the function.
		return f(roots)								// Not sure if it's worth it, but whatever
	})
}

func Operator(token Token) *AbstractParser {
	return NewCompiler(func (roots []*TokenRoot) []*TokenRoot {
		for _, root := range roots {
			root.Traverse(func (tree *TokenTree) {
				var first int
				for i, child := range tree.branches {
					if child.token == token {
						first = i
						break
					}
				}
				elem := new(TokenTree)
				elem.token = OPERATOR
				elem.branches = make([]*TokenTree, 3)

				elem.branches[0] = new(TokenTree)
				elem.branches[0].branches = tree.branches[:first]

				elem.branches[1] = new(TokenTree)
				elem.branches[1].token = token

				// This weird bit of array functions gets an empty array when
				// there are none to the right

				elem.branches[2] = new(TokenTree)
				array := tree.branches[first:] // If this doesn't work, pass by reference
				array = append(array, new(TokenTree))
				array = array[:len(array)-1] // Get rid of that element I added
				elem.branches[2].branches = array

				*tree = *elem

				// Well....
				// Values are being dropped and there are duplicate trees?
				// but whatever! It might work in the future!

			})
		}
		return roots // for chaining
	})
}

//// Token Trees

func (self Token) String() string {
	return strconv.Itoa(int(self))
}

func NewTokenTree(value string, token Token) *TokenTree {
	toktree := new(TokenTree)
	toktree.value = value
	toktree.token = token
	toktree.branches = make([]*TokenTree, 0)
	return toktree
}

func (self *TokenTree) Traverse(f func(*TokenTree))  {
	for _, child := range self.branches {
		child.Traverse(f) // #depthfirstsearch #bottomup?
	}
	f(self)
}

// func (self *TokenRoot) Tree() *TokenTree {
// 	tree := new(TokenTree)
// 	tree.branches = self.branches
// 	return tree
// }

func (self *TokenRoot) Traverse(f func(*TokenTree)) {
	for _, child := range self.branches {
		child.Traverse(f) // #depthfirstsearch #bottomup?
	}
}

func (self *TokenTree) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("'")
	buffer.WriteString(self.value)
	buffer.WriteString("'")
	buffer.WriteString("[")
	buffer.WriteString(self.token.String())
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

//// Lexers

func (self *Lexer) tokenize(results []*Result) ([]*TokenRoot, error) {
	outputs := make([]*TokenRoot, len(results))
	for i, result := range results {
		outputs[i] = new(TokenRoot)
		outputs[i].branches = []*TokenTree{NewTokenTree(result.value, self.token)}
	}
	var err error
	if len(results) == 0 {
		err = errors.New("Did not parse token " + self.token.String())
	} else {
		err = nil
	}
	return outputs, err
}

func (self *Lexer) retokenize(root_list []*TokenRoot, value string) []*TokenRoot {
	for i := range root_list {
		if len(value) > 0 { // Empty strings can be parsed, sure, but they shouldn't be added here.
			root_list[i].branches = append(root_list[i].branches, NewTokenTree(value, self.token))
		}
	}
	return root_list
}

func (self *Lexer) CanParse(str string) (val bool) {
	result_list := self.action(str)
	val = len(result_list) > 0
	if self.next != nil {
		val = val && self.next.CanParse(str) // Get the entire chain!
	}
	return
}

func (self *Lexer) Parse(str string) []*TokenRoot {
	parent_results := self.action(str)
	if self.next != nil {
		outputs := make([]*TokenRoot, 0)
		for _, a_result := range parent_results {
				child_results := self.next.Parse(a_result.left_over)
			// I think this append will work, but I'm not entirely sure.
			outputs = append(outputs, self.retokenize(child_results, a_result.value)...)
		}
		return outputs
	}
	root_list, _ := self.tokenize(parent_results)
	return root_list
}

func (self *Lexer) Chain(other ParserLike) ParserLike {
	if self.next == nil {
		self.next = other
	} else {
		self.next.Chain(other)
	}
	return other // faster for chaining.
}


//////////////////////////////////////////////

//// Functional Parsing

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
	var p *Lexer

	p = new(Lexer)
	p.token = SPACE
	p.action = satisfyRune(isSpace)

	s := new(Lexer)
	s.token = SPACE
	s.action = satisfyRune(isSpace)


	q := new(Lexer)
	q.token = LEFT
	q.action = satisfyRune(isChar('c'))

  r := new(Lexer)
  r.token = 5 // Just for testing purposes
  r.action = matchString("Hello") // MatchString is no longer greedy!

	t := new(Lexer)
	t.token = 0 // eof
	t.action = matchNothing()

	plus_op := new(Lexer)
	plus_op.token = 9 // Plus
	plus_op.action = satisfyRune(isChar('+'))

	c := Operator(9)

	c.Chain(q).Chain(p).Chain(plus_op).Chain(Optionally(s)).Chain(r)

	root_list := q.Parse("c + Hello World")

	for i := range root_list {
		fmt.Println(root_list[i])
	}

}
