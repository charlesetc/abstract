
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
  MANY
)

type Result struct {
  tokens []string
  left_over string
}

type Token struct {
  token string
  action Action
  children []*Token
}

func Base() *Token {
  tok := new(Token)
  tok.action = NONE
  tok.children = []*Token{}
  tok.token = ""
  // DOES NOT HAVE A TOKEN.
  return tok
}

func Tokenize(str string) *Token {
  tok := new(Token)
  tok.token = str
  tok.action = NONE
  tok.children = []*Token{}
  return tok
}

func And(tokens ...*Token) *Token {
  base := Base()
  base.action = AND
  base.children = tokens
  return base
}

func Many(token *Token) *Token {
  base := Base()
  base.action = MANY
  base.children = append(base.children, token)
  return base
}

func Optionally(token *Token) *Token {
  base := Base()
  base.action = OR    // more efficient than making a new list.
  base.children = append(base.children, token)
  return base
}

func SingleResult(toks []string, left_over string) []*Result {
  return []*Result{&Result{tokens: toks, left_over: left_over}}
}




func (self *Token) Compile(str string) ([]*Result) {
  if len(self.children) == 0 {
    if strings.HasPrefix(str, self.token) {
      // Return one result that has this token and the rest of the string.
      return SingleResult([]string{self.token}, str[len(self.token):])
    }
    return []*Result{}
  }

  // Otherwise there are children:

  current_list := SingleResult([]string{}, str)

  for _, child := range self.children {

    many_list := []*Result{}

    ThisChild:
    output_list := []*Result{}

    for _, result := range current_list {

      child_results := child.Compile(result.left_over)
      for _, res := range child_results {
        res.tokens = append(result.tokens, res.tokens...)
      }
      output_list = append(output_list, child_results...)

    }

    switch self.action {
    case OR:
      output_list = append(output_list, &Result{[]string{}, str})
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
    }

    current_list = output_list
  }

  return current_list
}


func PrintResults(results []*Result) {
  fmt.Print("[")
  for _, res := range results {
    fmt.Print(" (")

    for i, str := range res.tokens {
      if i != 0 {
        fmt.Print(" ")
      }
      fmt.Print(str)
    }
    fmt.Print(")")
    fmt.Print(res.left_over)
    fmt.Print(" ")
  }
  fmt.Println("]")
}


func main() {
  // a := Tokenize("a")
  c := Tokenize("c")
  b := Tokenize("b")
  d := Many(And(c, Optionally(b)))
  list_of_tokens := d.Compile("cbcbcb")
  PrintResults(list_of_tokens)
}
