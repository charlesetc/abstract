
# Abstract
#### Generate abstract syntax trees, easily, in Go.

* No code generation
* No additional languages
* Just a solid parsing library in Go

## Example First?
Let's get right to it, here's how you parse arithmetic expressions:

```go
plus := Lex("+")
minus := Lex("-")
times := Lex("*")
over := Lex("/")
spaces := Maybe(Munch(Space)).Garbage()

operator := OneOf(plus, minus, times, over)
number := And(spaces, Many(Digit), spaces).Alias("number")
mainlexer := And(Many(And(number, operator)), number)
result := lexer.MustCompile("2 + 5 - 3 * 8")

tree.Rule(Operator("*", 1, 1), Operator("/", 1, 1))
tree.Rule(Operator("+", 1, 1), Operator("-", 1, 1))

fmt.Println(tree)
```

Now all Abstract does is generate a syntax tree. Evaluating it is highly dependent on your needs. Curious to know how it works?


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

The `Many` function takes a single parser and accepts it multiple times. It's not exactly like the Regex version of `*` because it is not greedy. `Many` is similar to `Maybe` for this reason: It's also nondeterministic. Chances are you will use the `Munch` function more often, because it is greedy and deterministic. Furthermore, `Many` requires that there is at least one for a successful parse. If you want zero or more go with `Maybe(Many(...))`

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

`Alias` renames and groups a series of tokens. Let's say you're parsing a number than may or may not have underscores in the middle of it. (Like OCaml!) This would be done as so:

```go
part_of_number := Many(Digit)
underscore := Lex("_")
number := And(part_of_number, Maybe(Many(And(underscore, Digit)))
```

This is useful, but you're still going to get a list of tokens that looks somewhat like `["1", "2", "_", "0", "0", "0"]` when trying to parse `12_000`. Using `Alias`, we can say:

```go
number = Alias(number, "number")
```
or
```go
number = number.Alias("number")
```
This means that when we parse such a number, we only get the value `"12_000"` instead of a list of tokens. This is even more powerful when combined with `Garbage`

`Garbage(*Lexer) *Lexer`

`Garbage` throws away the value kept internally by a token. It will still parse correctly, but when using something like `Alias`, will not be transfered to the next value. For example:

```go
part_of_number := Many(Digit)
underscore := Lex("_").Garbage()
number := And(part_of_number, Maybe(Many(And(underscore, Digit)))
```
Now, `number` will compile an integer like `"12_000"` and hold the value `"12000"`

### Compiling a Lexer

I've mentioned that lexing this way is nondeterministic. For example, if we have a token `"*"` and a token `"**"`, it can either be parsed as two `"*"` tokens or one `"**"` token. And so, when you run `lexer.Compile("a sample string")` you'll end up with a list of possible results. This can be advantageous sometimes, but most of the time you're looking for the specific result that correctly parses the entire string you've given it. Using `lexer.MustCompile("a string")` is the way to do this.

Once you've gotten the results from a lexer, go ahead and make a tree out of them!

## Tree-Generation Phase

Given a specific result, `AbstractFromResult` makes a basic tree from the result. Then, you can add information about the operators to finish generating the tree.

Say I want to parse "2+3*3" into an Abstract Syntax Tree, a very simple example. My lexer would look something like this:

```go
operator := OneOf(Lex("+"), Lex("*"))
lexer := And(Digit, Many(And(operator, Digit)))
```

Now I compile the lexer against a string, and generate an initial tree:

```go
result := lexer.MustCompile("2+3*3")
tree := AbstractFromResult(result)
```

And finally, I add the operations I want to use:
```go
tree.Operator("*", 1, 1)
tree.Operator("+", 1, 1)
```

The two numbers following this operator indicates that it's an infix operator, in other words it
takes one argument on the left of it and one argument on the right. Note that the string used
here corresponds to the 'name' of the token you want to parse, which is either what you passed
to `Lex`, or any `Alias` you might have used.

You'll notice that operation precedence is decided by the order in which they're declared.
This is mostly intuitive, but if I want to add two operators with the same precedence,
I need additional syntax:

```go
tree.Rule(Operator("+", 1, 1), Operator("-"), 1, 1), ...)
```

Finally, Both of these parse in a left-associative manner. If you want a right-associative
operator, use the function `ROperator`.


## TODO

* Make an `NMunch` function which has the same effect as `NMany`.
* Make a `Between` function that walks a tree like `Rule`
* Implement Right-Associative operators.
