
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
  tok := Base()
  tok.action = AND
  tok.children = tokens
  return tok
}


// MAYBE MAKE NEW ONES INSTEAD?
func Many(token *Token) *Token {
  token.action = MANY
  return token
}


func Optionally(token *Token) *Token {
  base := Base()
  base.action = OR
  base.children = append(base.children, token)
  return base
}

                          // []([]tokens, left_over)
                          // 'cause nondeterministic
func (self *Token) Compile(str string) ([][][]string) {
  if len(self.children) == 0 {
    if strings.HasPrefix(str, self.token) {
      return [][][]string{[][]string{[]string{self.token}, []string{str[len(self.token):]}}}
    } // Otherwise: Failed to parse: keep going with just str
    return [][][]string{[][]string{[]string{}, []string{str}}}
  } // Otherwise: there are children

  // By convention, the [][]string is always 2 long
  // and the second parameter is always 1 long.
  current_list := [][][]string{[][]string{[]string{}, []string{str}}}
  output_list := [][][]string{} // add to this and then make it the current_list

  for _, child := range self.children {
    for _, list_of_tokens_and_left_over := range current_list {

      tokens := list_of_tokens_and_left_over[0]
      left_over := list_of_tokens_and_left_over[1][0]

      child_results := child.Compile(left_over)

      if self.action == OR {
        output_list = append(output_list, list_of_tokens_and_left_over)
      }

      for _, child_result := range child_results {
        child_tokens := child_result[0]
        child_left_over := child_result[1][0]

        tokens_so_far := append(tokens, child_tokens...)
        parsed := child_left_over != left_over

        switch self.action {
        case AND, OR:
          if parsed {
            output_list = append(output_)
          }
        case MANY:
        // case XOR:}
          // Do later
        default:
          panic(self.action)
          // This should not happen.
        }

      }
    }
    current_list = output_list
  }
  // There will be at least one child.
  return output_list
}

func main() {

  // a := Tokenize("a")
  // b := Tokenize("b")
  c := Tokenize("c")

  d := And(c, c)
  // d := And(a, b, c)
  // d := And(a, Optionally(And(b, c)), b)


  list_of_tokens := d.Compile("cac")
  fmt.Println(list_of_tokens)

  // // Aliases are built-in, when an upper-level operator w/ multiple children
  // if their token is nil.

  // I also need a "XOR/OneOf" function.

}



// // And
// if child_left_over != left_over { // Parsed!
//
//   // If I have a token, replace the tokens with mine.
//   var going_out_list [][]string
//   if self.token != "" {
//     going_out_list = [][]string{append(tokens_so_far, self.token), []string{child_left_over}}
//   } else {
//     going_out_list = [][]string{tokens_so_far, []string{child_left_over}}
//   }
//   // This might fuck things up:
//   output_list = [][][]string{going_out_list} //append(output_list, going_out_list)
// } else {
//
//   // continue LoopChildren // Might have done it
//   return [][][]string{[][]string{[]string{}, []string{str}}}
// } // I don't append anything here, because it didn't parse.


          // // OR
          // if child_left_over != left_over { // Parsed!
          //   // If I have a token, replace the tokens with mine.
          //   var going_out_list [][]string
          //   if self.token != "" {
          //     going_out_list = [][]string{append(tokens_so_far, self.token), []string{child_left_over}}
          //   } else {
          //     going_out_list = [][]string{tokens_so_far, []string{child_left_over}}
          //   }
          //
          //   // Might cause a bunch of duplicates...
          //   // Easily fixed, but harder to think about.
          //   output_list = append(output_list, going_out_list, list_of_tokens_and_left_over)
          //   fmt.Print("yes -- ")
          //   fmt.Println(child_left_over)
          // } else {
          //   fmt.Print("no --")
          //   fmt.Println(left_over)
          //   output_list = append(output_list, list_of_tokens_and_left_over)
          // }
