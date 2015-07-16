package main

import (
	. "abstract"
	"fmt"
)

func main() {
	spaces := Maybe(Munch(Space)).Garbage() // So funny
	plus := Lex("+")
	minus := Lex("-")
	times := Lex("*")
	over := Lex("/")
	operator := OneOf(plus, minus, times, over)
	number := And(spaces, Many(Digit), spaces).Alias("number")
	lexer := And(Many(And(number, operator)), number)
	str := "2 - 4 + 3 * 2"
	result := lexer.MustCompile(str)
	tree := AbstractParent(result.Tokens())

	tree.Rule(Operator("*", 1, 1), Operator("/", 1, 1))
	tree.Rule(Operator("+", 1, 1), Operator("-", 1, 1))
	tree.Operator("**", 1, 1) // This also works now.

	spaces = Maybe(Munch(Space).Garbage()).Alias("space")
	word := And(spaces, Munch(Alpha), spaces)
	number = And(spaces, Many(Digit), spaces).Alias("number")
	left := Lex("(")
	right := Lex(")")
	lexer = Many(OneOf(operator, number, word, left, right))
	result = lexer.MustCompile("2+(3*2)")
	tree = AbstractFromResult(result)
	tree.Between("(", ")")
	tree.Operator("*", 1, 1)
	tree.Operator("+", 1, 1)

	fmt.Println(tree)

}
