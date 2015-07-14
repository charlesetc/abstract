
# Abstract
#### Generate abstract syntax trees, easily, in Go.

* No code generation
* No additional languages
* Just a solid parsing library in Go

## Lexical Phase

Abstract provides you with a set of functions, which when combined
generate what is essentially a finite state machine, albeit in a functional manner.

This then produces a sequence of tokens which is then piped into the Tree-Generation phase. It's good to note that these functions have the same power as regex; they can tell if a string matches their configuration, even though that's not their design. Furthermore, they always read from the leftmost character of the string, unlike regex.

### Lexical Functions

`Lex(string) *Lexer`

`Lex` is the most fundamental method provided for parsing with Abstract. It takes a string and returns a Lexer that matches the string. All the other methods take in Lexers that are primarily built `Lex`.

`And(...*Lexer) *Lexer`

`And` is the second most important method. It combines a sequence of Lexers into a single Lexer. For example, in order to parse the string `abc`, you would combine the individual lexers as so: `And(Lex("a"), Lex("b"), Lex("c"))`

`Maybe(*Lexer) *Lexer`

`Maybe` is a nondeterministic part of Abstract. If you say `lexer := Maybe(Lex("a"))`, then one 'track' of parsing will attempt to parse the character `a` and the other will not. See the following examples:

```go
And(Lex("a"), Maybe(Lex("b")), Lex("c")).Compile("ac")  // => [a, c]
And(Lex("a"), Maybe(Lex("b")), Lex("c")).Compile("abc") // => [a, b, c]
```

Note, the return values are not simply arrays of strings. More on that later.

`Many(*Lexer) *Lexer`

The `Many` function takes a single parser and accepts it multiple times. It's not exactly like the Regex version of `*` because it is not greedy. `Many` is similar to `Maybe` for this reason: It's also nondeterministic. Chances are you will use the `Munch` function more often, because it is greedy and deterministic.

`NMany(*Lexer, int, (int)) *Lexer`

The NMany takes a Lexer, an integer, and optionally another integer.
There are two usages:

```go
NMany(Lex("a"), 3) // Will parse up to 3 "a"'s

NMany(Lex("a"), 2, 4) // Will parse between 2 and 4 "a"'s
                      // Note, this includes 4.
```

`Munch(*Lexer) *Lexer`

As stated before, `Munch` is a deterministic greedy version of `Many`. It will parse as many iterations of its lexer as possible. It's usually a good idea to use `Munch` unless you know you want to be nondeterministic and go with `Many`.

`Alias(*Lexer) *Lexer`

# Say more here


## Tree-Generation Phase


## TODO

Make an `NMunch` function which has the same effect as `NMany`.
