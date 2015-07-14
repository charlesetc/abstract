
package abstract

import (
  "testing"
)

//// Lexical Tests

var a, b, c *Lexer

func init()  {
  a = Lex("a")
  b = Lex("b")
  c = Lex("c")
}

func TestAnd(t *testing.T)  {
  // Most basic test.
  lexer := And(a, b)
  result := lexer.MustCompile("abc")
  if !CompareTokens(result.tokens, []string{"a", "b"}) {
    t.Error("Basic AND functionality")
  }

  // Multple Parameters
  lexer = And(a, b, c)
  results := lexer.Compile("abc")
  // PrintResults(results)
  result = lexer.MustCompile("abc")
  if !CompareTokens(result.tokens, []string{"a", "b", "c"}) {
    t.Error("AND with multiple parameters")
  }

  // And doesn't parse
  results = lexer.Compile("abd")
  if len(results) != 0 {
    t.Error("AND gives false positive")
  }
}

func TestOr(t *testing.T) {
  lexer := Maybe(a)
  results := lexer.Compile("abd")
  if len(results) < 2 {
    t.Error("Maybe returns too few possible tracks")
  } else if len(results) > 2 {
    t.Error("Maybe returns too many possible tracks")
  }
  result := lexer.MustCompile ("aab")

  str1 := "abc"
  str2 := "ac"
  lexer = And(a, Maybe(b), c)
  results = lexer.Compile(str1)
  if len(results) != 1 {
    t.Error("There's a problem with Maybe")
  }
  results = lexer.Compile(str2)
  if len(results) != 1 {
    t.Error("There's a problem with Maybe")
  }
  result = lexer.MustCompile(str1)
  if !CompareTokens(result.tokens, []string{"a", "b", "c"}) {
    t.Error("Maybe gives wrong tokens")
  }
  result = lexer.MustCompile(str2)
  if !CompareTokens(result.tokens, []string{"a", "c"}) {
    t.Error("Maybe gives wrong tokens")
  }
}

func TestOneOf(t *testing.T) {
  lexer := OneOf(a, b)
  results := lexer.Compile("ab")
  if len(results) != 1 {
    t.Error("OneOf doesn't parse correctly")
  }
  results = lexer.Compile("ba")
  if len(results) != 1 {
    t.Error("OneOf doesn't parse correctly")
  }
  results = lexer.Compile("cba")
  if len(results) != 0 {
    t.Error("OneOf doesn't parse correctly")
  }
}

func TestMany(t *testing.T) {
  lexer := Many(a)
  results := lexer.Compile("baaaaaaa")
  if len(results) != 0 {
    t.Error("Many gives false positives")
  }
  results = lexer.Compile("aaa")
  if len(results) != 3 {
    t.Error("Many doesn't parse correctly")
  }
  result := lexer.MustCompile("aaa")
  if len(result.tokens) != 3 {
    t.Error("Terrible algorithm here.")
  }
}

func TestNMany(t *testing.T)  {
  lexer := NMany(a, 2)
  results := lexer.Compile("baaaaa")
  if len(results) != 0 {
    t.Error("NMany gives false positives")
  }
  results = lexer.Compile("aaaaa")
  if len(results) != 2 {
    t.Errorf("NMany gives %d results instead of 2.", len(results))
  }

  lexer = NMany(a, 2, 3)
  results = lexer.Compile("a")
  if len(results) != 0 {
    t.Error("NMany gives false positives")
  }
  results = lexer.Compile("aa")
  if len(results) != 2 {
    t.Error("NMany doesn't work with multiple parameters")
  }
  results = lexer.Compile("aaa")
  if len(results) != 3 {
    t.Error("NMany doesn't work with multiple parameters")
  }
}

func TestMunch(t *testing.T) {
  lexer := Munch(a)
  results := lexer.Compile("baaaaaaa")
  if len(results) != 0 {
    t.Error("Munch gives false positives")
  }
  results = lexer.Compile("aaa")
  if len(results) != 1 {
    t.Error("Munch doesn't parse correctly")
  }
  if len(results[0].tokens) != 3 {
    t.Error("Munch doesn't parse correctly")
  }
}

func TestLexicalIntegration(t *testing.T)  {
  period := Lex(".").Alias("period")
  at := Alias(Lex("@"), "atsign")
  alphahyphen := OneOf(Alphanumeric, Lex("-"))
  punct :=  OneOf(
            Lex("!"),
            Lex("#"),
            Lex("$"),
            Lex("%"),
            Lex("&"),
            Lex("'"),
            Lex("*"),
            Lex("+"),
            Lex("-"),
            Lex("/"),
            Lex("="),
            Lex("?"),
            Lex("^"),
            Lex("_"),
            Lex("`"),
            Lex("{"),
            Lex("|"),
            Lex("}"),
            Lex("~"))
  char := OneOf(punct, Alpha)
  // Periods can't be at the beginning or end,
  // and they can't appear twice in a row.
  email := And(Many(And(char, Maybe(period))), char, at, alphahyphen, Many(alphahyphen), period, NMany(alphahyphen, 2, 4), Eof)
  address := "test@test.com"
  if !email.Match(address) {
    t.Errorf("Email did not match %s", address)
  }
  address = "test.test@test.com"
  if !email.Match(address) {
    t.Errorf("Email did not match %s", address)
  }
  address = "Test._test@gtest-test.i2fo"
  if !email.Match(address) {
    t.Errorf("Email did not match %s", address)
  }
  address = "Test._test@gtest-test.co"
  if !email.Match(address) {
    t.Errorf("Email did not match %s", address)
  }
  address = "test@test.c"
  if email.Match(address) {
    t.Errorf("Email matched %s and should not have", address)
  }
  address = "test@test.cdfsa"
  if email.Match(address) {
    t.Errorf("Email matched %s and should not have", address)
  }
}

func TestGarbage(t *testing.T) {
  lexer := And(a,b,c.Garbage()).Alias("hello")
  result := lexer.MustCompile("abc")
  if result.Tokens()[0].Value != "ab"{
    t.Errorf("Either Alias() or Garbage() doesn't work.")
  }
}

//// AST Testing

func TestBasicOperator(t *testing.T)  {
  result := And(a, b, c).MustCompile("abc")
  tree := AbstractParent(result.Tokens())
  tree.Operator("b", 1, 1)
  if len(tree.Children) != 1 {
    t.Error("Operator is not grouping properly")
  }
}

func TestMultipleOperators(t *testing.T)  {
  result := And(a, b, c, b, c).MustCompile("abcbc")
  tree := AbstractParent(result.Tokens())
  tree.Operator("b", 1, 1)
  if len(tree.Children) != 1 {
    t.Error("Operator is not grouping properly with multiple operations")
  }
}

func TestLeftAssociativity(t *testing.T)  {
  result := And(a, b, c, b, c).MustCompile("abcbc")
  tree := AbstractParent(result.Tokens())
  tree.Operator("b", 1, 1)
  if tree.Children[0].Children[0].Children[0].Children[0].Children[0].Token.Value != "a" {
    t.Error("Operator is not by default left-associtive")
  }
}

// Testing Helper Functions

func CompareTokens(tokens []*Token, strings []string) bool {
  for i, tok := range tokens {
    if tok.Value != strings[i] {
      return false
    }
  }
  return true
}
