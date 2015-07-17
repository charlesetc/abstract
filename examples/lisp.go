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
	integer := Munch(Digit).Alias("integer")
	float := And(Munch(Digit), Lex("."), Munch(Digit)).Alias("float")

	number := OneOf(integer, float) // Wish it could be FirstOf



	lexer := Many(OneOf(paren, word, number, operator, spaces))


	result := lexer.MustCompile("2.0 + 2")
	tree := AbstractFromResult(result)

	tree.Filter("space")
	tree.Select("float").Apply(Many(OneOf(integer, Lex("."), integer)))

	fmt.Println(tree)

	result = lexer.MustCompile("(def + hi (fn [2 + 3 - 1] this) this)")

	/* AST */
	tree = AbstractFromResult(result)
	tree.Filter("space")
	tree.Between("[", "]")
	tree.Between("(", ")")
	lists := tree.Select("[]")
	lists.Operator("+", 1, 1)
	lists.Operator("-", 1, 1)

	fmt.Println()

	// fmt.Println(tree)
}
