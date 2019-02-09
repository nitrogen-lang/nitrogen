# Semicolons

Nitrogen's formal grammar uses semicolons to denote the end of a statement or expression.
However, Nitrogen uses automatic semicolon insertion where possible.
Semicolons are inserted after an identifier, literal, nil, a closing parenthesis,
curly or square bracket, and after the keywords `return`, `break`, and `continue`.
This requires the programmer to keep in mind how they're formatting code.
For example, the following if statement is invalid:

```
if (var == 3) // A semicolon is inserted here which will fail
{
    print("It's a 3")
} // Same here
else // And here
{
    print("No, it's not 3")
}
```

But the following is syntactically valid (though not recommended):

```
if (var == 3) { print("It's a 3") } else { print("No, it's not 3") }
```

The parser will catch any such errors and will warn the programmer.
