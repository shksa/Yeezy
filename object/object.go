package object

import (
	"fmt"
)

/* Why Object?
- Object is an abstraction for "data" in the lang.
- An object will have a type and a value.
- All the primitive and composite types of Monkey are the various "kinds" of the objects.
- An Integer type in Monkey is an Object with a particular type name called "INTEGER" and can take all 64-bit integer values
	of Golang.
- A Boolean type in Monkey is an Object with a particular type name called "BOOLEAN" can take either the true or false boolean
	values of Golang.
*/

// Object is an interface which is implemented by all types in  Monkey lang.
type Object interface {
	Inspect() string // returns value of the object
	Type() string    // returns type of the object
}

// list of all the types in Monkey
const (
	INTEGER = "INTEGER"
	BOOLEAN = "BOOLEAN"
	NULL    = "NULL"
)

/* IMPORTANT
- Every value has a different representation in the Host lang, so each type of value is represented by a struct.
- Whenever interger literals are encountered in the source, they are turned into an ast.IntegerLiteral and then
	when evaluating the AST node, we turn it into an value.Integer, saving the value inside the struct.
*/

// Integer is type for representing all integer literal values in the Monkey lang.
type Integer struct {
	Value int64
}

// Inspect returns the value in string format
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns the type's name
func (i *Integer) Type() string { return INTEGER }

// Boolean is type for representing all boolean literal values in the Monkey lang.
type Boolean struct {
	Value bool
}

// Inspect returns the value in string format
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Type returns the type's name
func (b *Boolean) Type() string { return BOOLEAN }

// Null is type for representing the absence of values in Monkey
type Null struct{}

// Inspect returns the value in string format
func (n *Null) Inspect() string { return "null" }

// Type returns the type's name
func (n *Null) Type() string { return NULL }
