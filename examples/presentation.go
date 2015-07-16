package main

import (
	. "abstract"
	"fmt"
)

func main() {
	// Lexing

	plus := Lex("+")
	minus := Lex("-")
	times := Lex("*")

	paren := OneOfString("(", ")")

	operator := OneOf(plus, minus, times)

	spaces := Munch(Space).Alias("space")
	number := Many(Digit).Alias("number")

	lexer := Many(OneOf(number, paren, operator, spaces))

	result := lexer.MustCompile("2 + (4 * 2) - 3")

	// AST
	tree := AbstractFromResult(result)

	tree.Filter("space")

	tree.Between("(", ")")

	tree.Operator("*", 1, 1)
	tree.Rule(Operator("+", 1, 1), Operator("-", 1, 1))

	fmt.Println(tree)
}
