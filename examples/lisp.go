package main

import (
	. "abstract"
	"fmt"
)

func main() {

	/* Lexing */
	left_paren := Lex("(")
	right_paren := Lex(")")
	left_square := Lex("[")
	right_square := Lex("]")
	paren := OneOf(left_paren, right_paren, left_square, right_square)
	spaces := Munch(Space).Alias("space")
	punct := OneOfString("_", "-")
	operator := OneOfString("+", "-", "*", "-", ">", "<")
	word := Munch(OneOf(Alpha, punct)).Alias("word")
	number := Munch(Digit) //And(Munch(Digit), Maybe(And(Lex("."), Munch(Digit)))).Alias("number")

	lexer := Many(OneOf(paren, word, number, spaces, operator))
	result := lexer.MustCompile("(def + hi (fn [2 + 3 - 1] this) this)")

	/* AST */
	tree := AbstractFromResult(result)
	tree.Filter("space")
	tree.Between("[", "]")
	tree.Between("(", ")")
	lists := tree.Select("[]")
	lists.Operator("+", 1, 1)
	lists.Operator("-", 1, 1)

	fmt.Println(tree)
}
