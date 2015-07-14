
package abstract

import (
  "bytes"
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
  NMANY
  OPERATOR // Only used in Syntax Tree part.
)

type Token struct {
  Name string
  Value string
}

type operator struct {
  name string
  left int
  right int
}

func Operator(name string, left int, right int) *operator {
  return &operator{name, left, right}
}

type Result struct {
  tokens []*Token
  left_over string
}

func (r *Result) Tokens() []*Token {
  return r.tokens
}

type Lexer struct {
  token string
  action Action
  to int // Only used for nmany
  from int // Only used for nmany
  children []*Lexer
}

func base() *Lexer {
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
  b := base()
  b.action = AND
  b.children = tokens
  return b
}

func OneOf(tokens ...*Lexer) *Lexer {
  b := base()
  b.action = XOR
  b.children = tokens
  return b
}

func Many(token *Lexer) *Lexer {
  b := base()
  b.action = MANY
  b.children = append(b.children, token)
  return b
}

func NMany(token *Lexer, numbers... int) *Lexer {
  var from, to int
  switch len(numbers) {
  case 1:
    from = 0
    to = numbers[0]
  case 2:
    from = numbers[0] + 1
    to = numbers[1]
  default:
    panic("NMany requires (1..2) numbers after the lexer.\nExample: NMany(myLexer, 4, 0) // Between 0 to 4")
  }
  b := base()
  b.action = NMANY
  b.to = to
  b.from = from
  b.children = append(b.children, token)
  return b
}

func Munch(token *Lexer) *Lexer {
  b := base()
  b.action = MUNCH
  b.children = append(b.children, token)
  return b
}

func Maybe(token *Lexer) *Lexer {
  b := base()
  b.action = OR    // more efficient than making a new list.
  b.children = append(b.children, token)
  return b
}

func Alias(token *Lexer, str string) *Lexer {
  b := base()
  b.action = NONE // don't know what this will do
  b.children = append(b.children, token)
  b.token = str
  return b
}

func (self *Lexer) Alias(str string) *Lexer {
  return Alias(self, str)
}

func (self *Lexer) Garbage() *Lexer {
  b := base()
  b.children = append(b.children, self)
  b.token = "abstract://garbage"
  return b
}

func singleResult(toks []*Token, left_over string) []*Result {
  return []*Result{&Result{tokens: toks, left_over: left_over}}
}

func newToken(str string) *Token {
  return &Token{str, str}
}

func (self *Token) String() string {
  return self.Name + ":" + self.Value
}

func (self *Lexer) Compile(str string) ([]*Result) {
  if len(self.children) == 0 {
    s := string(append([]byte(str), 0))
    if strings.HasPrefix(s, self.token) {
      // Return one result that has this token and the rest of the string.
      if self.token == string([]byte{0}) {
        return singleResult([]*Token{newToken(self.token)},str)
      }
      return singleResult([]*Token{newToken(self.token)}, str[len(self.token):])
    }
    return []*Result{}
  }

  // Otherwise there are children:

  current_list := singleResult([]*Token{}, str)

  xor_list := []*Result{}

  NextChild:
  for _, child := range self.children {

    many_list := []*Result{}
    many_count := 1

    ThisChild:
    output_list := []*Result{}

    for _, result := range current_list {

      child_results := child.Compile(result.left_over)
      for _, res := range child_results {

        var the_tokens []*Token
        if self.token != "" {

          value := ""

          if self.token != "abstract://garbage" {
            for _, tok := range res.tokens {
              value = value + tok.Value
            }
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
    case NMANY, MANY:
      if len(output_list) == 0 || (self.action == NMANY && many_count > self.to) {
        if self.action == NMANY && many_count < self.from {
          output_list = []*Result{}
          break
        }
        output_list = many_list
        break
      } // Otherwise:
      many_list = append(many_list, output_list...)
      many_count++
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

func (l *Lexer) Match(str string) bool {
  return len(l.Compile(str)) > 0
}



func PrintResult(res *Result, answers...bool)  {
  fmt.Print("(")

  for i, tok := range res.tokens {
    if i != 0 {
      fmt.Print(" ")
    }
    fmt.Print(tok.Name)
    fmt.Print(":")
    fmt.Print(tok.Value)
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

// Example:
// func main() {
//   digit := OneOf(Lex("1"), Lex("2"), Lex("3"), Lex("4"), Lex("5"), Lex("6"), Lex("7"), Lex("8"), Lex("9"), Lex("0"))
//   integer := Alias(Munch(digit), "int")
//   float := And(integer, Maybe(And(Lex("."), integer)))
//
//   results := float.Compile("13.20")
//   PrintResults(results)
// }


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
var Eof *Lexer = Lex(string([]byte{0}))
var Space = OneOf(Lex(" "), Lex("\n"), Lex("\t"))


//// Abstract Syntax Trees.

type Abstract struct {
  Token *Token
  Children []*Abstract
}

func (self *Abstract) String() string {
  var out bytes.Buffer
  if self.Token != nil {
    out.WriteString(self.Token.String())
  }
  out.WriteString("[")
  for i, child := range self.Children {
    if i != 0 {
      out.WriteString(" ")
    }
    out.WriteString(child.String())
  }
  out.WriteString("]")
  return out.String()
}


func AbstractFromToken(token *Token) *Abstract {
  return &Abstract{Token: token, Children: []*Abstract{}}
}

func AbstractWithName(str string) *Abstract {
  t := new(Token)
  t.Name = str
  return AbstractFromToken(t)
}

func AbstractParent(tokens []*Token) *Abstract {
  children := make([]*Abstract, len(tokens))
  for i, tok := range tokens {
    children[i] = AbstractFromToken(tok)
  }
  a := new(Abstract)
  a.Children = children
  // Leaving token nil.
  return a
}

// Assumes side-effects
func (self *Abstract) Walk(f func(*Abstract)) {
  for _, child := range self.Children {
    child.Walk(f)
  }
  f(self)
}

func (self *Abstract) printChildren() {
  fmt.Print(self.Token)
  fmt.Print("(")
  for i, child := range self.Children {
    if i != 0 {
      fmt.Print(" ")
    }
    fmt.Print(child.Token)
  }
  fmt.Println(")")
}

// Default left-associative // right-associative -> reverse list.
func (self *Abstract) Rule(ops ...*operator) {
  self.Walk(func (abstract *Abstract) {
    var left_number int
    var right_number int
    var name string
    var matched bool

    i := 0
    for i < len(abstract.Children) {

      child := abstract.Children[i]
      matched = false

      for _, op := range ops {
        if child.Token.Name == op.name {
          left_number = op.left
          right_number = op.right
          name = op.name
          matched = true
          break
        }
      }

      if matched {
        if i < left_number {
          panic(fmt.Sprintf("Rule %s needs %d tokens to its left.", name, left_number))
        } else if i > (len(abstract.Children) - right_number) {
          panic(fmt.Sprintf("Rule %s needs %d tokens to its right.", name, right_number))
        }

        left := AbstractWithName("abstract://right")
        right := AbstractWithName("abstract://left")

        alternative_children := make([]*Abstract, len(abstract.Children))
        copy(alternative_children, abstract.Children)

        left.Children = alternative_children[i-left_number:i]
        right.Children = alternative_children[i+1:i+1+right_number]
        child.Children = []*Abstract{left, right}

        abstract.Children = append(append(abstract.Children[:i-left_number], abstract.Children[i]), abstract.Children[i+1+right_number:]...)

        i = i - (len(left.Children) + len(right.Children) - 1)
      }
      i++
    }
  })
}
