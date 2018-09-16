## Type of interpreter that i will build
- It's called a **tree-walk** interpreter.
- The interpreter will read the source code, parse it into a AST and then evaluate the AST.

## Lexer
- Takes source code as input and generates tokens as output that represent the source code.

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

## Making AST
- **Tokens -> AST**
- Tokens are parsed into nodes in AST.
- A statement consists of a set of tokens, parser identifies the type of statement by looking at the
    first token in the statement and and then builds the AST for that statement.
- Tokens are parsed into their corresponding nodes and put in the AST.


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
    - **Integer literals**
        - The value they produce is the integer itself.
    - **Prefix-Operators** or **Prefix-Expressions**
        - There are 2 prefix operators in Monkey:- `!` and `-`.
        - Usage:
            ```
            -5;
            !foobar;
            5 + -10;
            ```
        - They have form `<prefix operator><expression>;`
        - **Any expression can follow a prefix operator as an operand.**
    - **Infix-Expressions**
        - They have the form `<expression> <infix-operator> <expression>`
    - **Boolean literals**
    - **Grouped expressions**
    - **If expressions**
        - In Monkey, if-else conditionals are expressions.
        - They have the form `if (<condition>) <consequence> else <alternative>`
    - **Function literals**
        - The only way to define functions.
        - form - `func <parameters> <block statement>`
        - `<parameters> := (<parameter1>, <parameter2>, <parameter3>, ...)`
    - **Call expressions**.
        - form - `Identifier or Function literal(<comma seperated expressions>)`

## Fuzzy line b/w interpreters and compilers
- The line between interpreters and compilers is a blurry one. 
- The notion of an interpreter as something that doesn't leave executable artifacts behind (in contrast to a compiler, which does just that) gets fuzzy real fast when looking at the implementations of real-world and highly-optimized programming languages.

## Evaluation
- Tree-walking interpreter
    - The most obvious and classical choice of what to do with the AST is to just interpret it. 
    - Traverse the AST, visit each node and do what the node signifies: print a string, add two numbers, execute a function's body      all on the fly. 
    - Interpreters working this way are called **"tree-walking interpreters"** and are the archetype of interpreters. 
    - Sometimes their evaluation step is preceded by small optimizations that rewrite the AST (e.g. remove unused variable bindings) or convert it into another intermediate representation (IR) that's more suitable for recursive and repeated evaluation.
- Bytecode emitter type 
    - Instead of building an AST the parser emits bytecode directly. 
    - Now, are we still talking about interpreters or compilers? Isn't emitting bytecode that gets then interpreted (or should we say "executed"?) a form of compilation? I told you: the line becomes blurry. 
    - And to make it even more fuzzy, consider this: some implementations of programming languages parse the source code, build an AST and convert this AST to bytecode. 
    - But instead of executing the operations specified by the bytecode directly in a virtual machine, the virtual machine then compiles the bytecode to native machine code, right before its executed - just in time. 
    - That's called a JIT (for "just in time") interpreter/compiler. 
- AST -> JIT
    - They recursively traverse the AST but before executing a particular branch of it the node is compiled to native machine code. - And then executed. Again, "just in time".

## Tree-walking interpreter
- It take the AST our parser builds for us and interprets it "on the fly", without any preprocessing or compilation step.

## Representing Monkey values in Host language
- We need a way to represent values of Monkey in the Host language.
- We need to design a **value sysytem** or an **object system**.
- This sysytem defines what our **eval** function returns.
- There are different ways to do this.
- Some use the native types (integers, booleans) of the host language to represent values of the interpreted language not wrapped in anything.

## Foundation of our object system
- Every value in Monkey will be represented as an **Object**, an interface of our design.
- The values of monkey are are the values of the custom **Object** interface type.
- Each value will also have a specific representation by custom struct types like **Integer**, **Boolean**.
- So an integer in Monkey is a value of **Integer** struct, which implements the **Object** interface.
- On the REPL, the integer value is shown by callong the **Inspect()** method of the **Integer** struct value.

## Evaluating expressions and statements
- The function signature is `func Eval(node ast.Node) object.Object`
- Traverse the AST, and evaluate the nodes. 
- **Self-evaluating expressions**
    - **Integer literals and Boolean literals**
        - ex:- single-line input `5`. This program consists of a single statement, an expression statement with an integer          literal as it's expression.
    - **Prefix-expressions**
        - The first step to evaluating a prefix expression is to evaluate its operand and then use the result of this evaluation with the operator.
    - **Infix-expressions**
        - There are 8 infix operators in Monkey. They can be divided into 2 groups, one produces integer values
            the other produces boolean values
    - **If-Else expression**
        - In the condition of the if expression, anything other than NULL and false are evaluated to TRUE.
- **Return statements**
    - Can be used as a top-level statement in the program.
    - Statements after the it wont be evaluated.
    - We need to keep track of the return value so that we can later decide whether to stop evaluation or not.