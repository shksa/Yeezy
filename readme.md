## Type of interpreter that i will build
- It's called a **tree-walk** interpreter.
- The interpreter will read the source code, parse it into a AST and then evaluate the AST.

## Lexer
- Takes source code as input and generates tokens as output that represent the source code.

## Why semi-colon?
- semicolan represents a character which is neither a number nor a letter.
- recognizing numbers or letter names involves reading successive chars continuously untill a char that is neither a digit or a letter is encountered respectively.
- If such a char is encountered, the mactching loop breaks and the number or word is recognized.
- So we NEED a char to represent neither a digit nor a letter so that the matching loop breaks.

## Parser
- A parser takes source code as input (either as text or tokens) and produce a data structure which represents this source code.
- In most interpreters and compilers the data structure used for the internal representation of the source code is called a "syntax tree" or an "abstract syntax tree" (AST for short). 
- The "abstract" is based on the fact that certain details visible in the source code are omitted in the AST. 
- Semicolons, newlines, whitespace, comments, braces, bracket and parentheses -- depending on the language and the parser, these details are not represented in the AST, but merely guide the parser when constructing it.

## Nodes in AST
- Nodes in AST of Monkey can be of 2 types.
    - **Statement**
    - **Expression**

## Statements and Expressions.
- ***Programs in Monkey are a series of statements.***
- A **let** statement has the form ```let <identifier> = <expression>```
- A **return** statement has the form ```return <expression>```
- ***There are only 2 types of statements in Monkey, a let statement and a return statement***.
- ***The rest of the language consists of expressions***.
- **Statements vs Expressions**
    - Expressions produce values. ```ex:- 5, 10, add(5, 10)```
    - Statements do not produce values. ```ex:- return 5, let x = 5```
    - Depends on the language though.
- **A lot of things in monkey are expressions including function literals**.

## Expression Statements
- ***An expression statement is a top-level statement that consists solely of one expression.***
- We need it because it's totally legal in Monkey to write the following
    ```
    let x = 5;
    x * 5 + 10;
    ```
- The first line is a "let" statement.
- The second line is a "expression" statement.
- Most scripting languages have expression statements.
- **They make it possible to have one line consisting only of an expression.**
- **So this type of statement should be represented as it's own type of node in the AST.**


## Parsing expressions
- **Top-Down Operator Precedence Parsing or Pratt parsing**
    - The main idea of a Pratt's parser is that parsing functions are associated with token types.
    - Whenever a particular token type is encountered, the associated parsing functions are called to parse the expression
        and return an AST node that represents it.
    - Each token type can have upto 2 parsing functions depending on whether the token is found in a infix or a prefix
        position.
- Types of expressions
    - **Identifiers** 
      - identifiers are expressions just like 1 + 2.
      - identifiers produce values just like other expressions.
      - identifiers evaluate to the value they are bound to.