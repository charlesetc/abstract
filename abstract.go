
package main

import (
  "strings"
  "fmt"
)

type Action int

const (
  NONE Action = iota
  AND
  OR
  XOR
  MANY
  MUNCH
)

type Token struct {
  name string
  value string
}

type Result struct {
  tokens []*Token
  left_over string
}

type Lexer struct {
  token string
  action Action
  children []*Lexer
}

func Base() *Lexer {
  tok := new(Lexer)
  tok.action = NONE
  tok.children = []*Lexer{}
  tok.token = ""
  // DOES NOT HAVE A TOKEN.
  return tok
}

func Lex(str string) *Lexer {
  tok := new(Lexer)
  tok.token = str
  tok.action = NONE
  tok.children = []*Lexer{}
  return tok
}

func And(tokens ...*Lexer) *Lexer {
  base := Base()
  base.action = AND
  base.children = tokens
  return base
}

func OneOf(tokens ...*Lexer) *Lexer {
  base := Base()
  base.action = XOR
  base.children = tokens
  return base
}

func Many(token *Lexer) *Lexer {
  base := Base()
  base.action = MANY
  base.children = append(base.children, token)
  return base
}

func Munch(token *Lexer) *Lexer {
  base := Base()
  base.action = MUNCH
  base.children = append(base.children, token)
  return base
}

func Maybe(token *Lexer) *Lexer {
  base := Base()
  base.action = OR    // more efficient than making a new list.
  base.children = append(base.children, token)
  return base
}

func Alias(token *Lexer, str string) *Lexer {
  base := Base()
  base.action = NONE // don't know what this will do
  base.children = append(base.children, token)
  base.token = str
  return base
}

func (self *Lexer) Alias(str string) *Lexer {
  return Alias(self, str)
}

func SingleResult(toks []*Token, left_over string) []*Result {
  return []*Result{&Result{tokens: toks, left_over: left_over}}
}

func NewToken(str string) *Token {
  return &Token{str, str}
}



func (self *Lexer) Compile(str string) ([]*Result) {
  if len(self.children) == 0 {
    if strings.HasPrefix(str, self.token) {
      // Return one result that has this token and the rest of the string.
      return SingleResult([]*Token{NewToken(self.token)}, str[len(self.token):])
    }
    return []*Result{}
  }

  // Otherwise there are children:

  current_list := SingleResult([]*Token{}, str)

  xor_list := []*Result{}

  NextChild:
  for _, child := range self.children {

    many_list := []*Result{}

    ThisChild:
    output_list := []*Result{}

    for _, result := range current_list {

      child_results := child.Compile(result.left_over)
      for _, res := range child_results {
        var the_tokens []*Token
        if self.token != "" {
          value := ""
          for _, tok := range res.tokens {
            value = value + tok.value
          }
          the_tokens = []*Token{&Token{self.token, value}}
        } else {
          the_tokens = res.tokens
        }
        res.tokens = append(result.tokens, the_tokens...)
      }
      output_list = append(output_list, child_results...)

    }

    switch self.action {
    case OR:
      output_list = append(output_list, &Result{[]*Token{}, str})
    case AND:
      if len(output_list) == 0 {
        return []*Result{}
        // Essentially, don't go on to the remaining children.
      }
    case MANY:
      if len(output_list) == 0 {
        output_list = many_list
        break
      } // Otherwise:
      many_list = append(many_list, output_list...)
      current_list = output_list
      goto ThisChild
    case XOR:
      if len(output_list) == 0 {
        continue NextChild
      } else {
        xor_list = append(xor_list, output_list...)
        continue NextChild
      }
    case MUNCH:
      if len(output_list) == 0 {
        output_list = many_list
        break
      } // Otherwise:
      many_list = output_list // Basically save the last one.
      current_list = output_list
      goto ThisChild
    }

    if self.action != XOR {
      current_list = output_list
    }
  }

  if self.action == XOR {
    return xor_list
  }

  return current_list
}

func (l *Lexer) MustCompile(str string) *Result {
  results := l.Compile(str)
  switch len(results) {
  case 0:
    panic("Lexer did not compile.")
  case 1:
    return results[0]
  default:
    for _, res := range results {
      if res.left_over == "" {
        return res
      }
    }
    return results[1]
    // Not a very smart algorithm.
  }
}

func PrintResult(res *Result, answers...bool)  {
  fmt.Print("(")

  for i, tok := range res.tokens {
    if i != 0 {
      fmt.Print(" ")
    }
    fmt.Print(tok.name)
    fmt.Print(":")
    fmt.Print(tok.value)
  }
  fmt.Print(")")
  fmt.Print(res.left_over)
  fmt.Print(" ")
  if len(answers) == 0 {
    fmt.Println()
  }
}

func PrintResults(results []*Result) {
  fmt.Println("[")
  for _, res := range results {
    PrintResult(res)
  }
  fmt.Println("]")
}

func main() {
  // a := Lex("a")
  // c := Lex("c")
  // // b := Lex("b")
  // d := OneOf(c, Maybe(Alias(And(c, a), "wow")))
  digit := OneOf(Lex("1"), Lex("2"), Lex("3"), Lex("4"), Lex("5"), Lex("6"), Lex("7"), Lex("8"), Lex("9"), Lex("0"))
  integer := Alias(Munch(digit), "int")
  float := And(integer, Maybe(And(Lex("."), integer)))

  results := float.Compile("13.20")
  PrintResults(results)
}


// I am aware how ugly this is.
var Digit *Lexer =  OneOf(Lex("0"),
                    Lex("1"), Lex("2"), Lex("3"),
                    Lex("4"), Lex("5"), Lex("6"),
                    Lex("7"), Lex("8"), Lex("9"))
var Lower *Lexer =  OneOf(
                    Lex("a"), Lex("b"), Lex("c"),
                    Lex("d"), Lex("e"), Lex("f"),
                    Lex("g"), Lex("h"), Lex("i"),
                    Lex("j"), Lex("k"), Lex("l"),
                    Lex("m"), Lex("n"), Lex("o"),
                    Lex("p"), Lex("q"), Lex("r"),
                    Lex("s"), Lex("t"), Lex("u"),
                    Lex("v"), Lex("w"), Lex("x"),
                    Lex("y"), Lex("z"))

var Upper *Lexer =  OneOf(
                    Lex("A"), Lex("B"), Lex("C"),
                    Lex("D"), Lex("E"), Lex("F"),
                    Lex("G"), Lex("H"), Lex("I"),
                    Lex("J"), Lex("K"), Lex("L"),
                    Lex("M"), Lex("N"), Lex("O"),
                    Lex("P"), Lex("Q"), Lex("R"),
                    Lex("S"), Lex("T"), Lex("U"),
                    Lex("V"), Lex("W"), Lex("X"),
                    Lex("Y"), Lex("Z"))

var Alpha *Lexer =  OneOf(Upper, Lower)
var Alphanumeric *Lexer = OneOf(Alpha, Digit)
