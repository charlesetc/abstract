
package main

import (
  . "abstract"
  "fmt"
)

func main() {
  left_paren := Lex("(")
  right_paren := Lex(")")
  left_square := Lex("[")
  right_square := Lex("]")
  paren := OneOf(left_paren, right_paren, left_square, right_square)
  spaces := Munch(Space)
  punct := OneOfString("_", "-")
  operator := OneOfString("+", "-", "*", "-", ">", "<")
  word := Munch(OneOf(Alpha, punct)).Alias("word")
  number := And(Munch(Digit), Maybe(And(Lex("."), Munch(Digit)))).Alias("number")

  lexer := Many(OneOf(paren, word, number, spaces, operator))
  result := lexer.MustCompile("(def hi (fn [2 3] this))")
  tree := AbstractFromResult(result)
  tree.Filter(" ")
  tree.Between("[", "]")
  tree.Between("(", ")")
  tree.Operator("word", 0, 2)
  fmt.Println(tree)
}
