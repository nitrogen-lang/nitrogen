## Arithmetic Operators

Here's a table that show all the binary operators and what number types that can be used on

| Op | Name                  | Types                        |
|----|-----------------------|------------------------------|
| +  |  sum                  |  integers, floats, strings   |
| -  |  difference           |  integers, floats            |
| *  |  product              |  integers, floats            |
| /  |  quotient             |  integers, floats            |
| %  |  remainder            |  integers, floats            |
|    |                       |                              |
| &  |  bitwise AND          |  integers                    |
| \| |  bitwise OR           |  integers                    |
| ^  |  bitwise XOR          |  integers                    |
| &^ |  bit clear (AND NOT)  |  integers                    |
|    |                       |                              |
| << |  left shift           |  integer << unsigned integer |
| >> |  right shift          |  integer >> unsigned integer |
|    |                       |                              |
| += |  sum assign           |  integers, floats, strings   |
| -= |  difference assign    |  integers, floats            |
| *= |  product assign       |  integers, floats            |
| /= |  quotient assign      |  integers, floats            |
| %= |  remainder assign     |  integers, floats            |

## Operator Precedence

There are 5 main precedence levels for binary operators. The operators bind strongest from highest
level to lowest level. Operators on the same level are left associative and will bind left to right.

| Level | Operators          |
|:-----:|--------------------|
|   5   | `* / % >> << & &^` |
|   4   | `+ - \| ^`         |
|   3   | `< >`              |
|   2   | `== != <= >=`      |
|   1   | `and or`           |
