package object

import (
	"fmt"
)

/* Why Object?
- Object is an abstraction for "data" in the lang.
- An object will have a type and a value.
- All the primitive and composite types of Monkey are the various "kinds" of the objects.
- The "types" of Monkey are defined by various custom structure types in Go that implement the Object interface.
- An Integer type in Monkey is an Object with a particular type name called "INTEGER" and can take all 64-bit integer values
	of Golang.
- A Boolean type in Monkey is an Object with a particular type name called "BOOLEAN" can take either the true or false boolean
	values of Golang.
*/

// Object is an interface which is implemented by all types in  Monkey lang.
type Object interface {
	Inspect() string // returns value of the object in string format
	Type() string    // returns type of the object
}

// list of all the types of Objects in Monkey
const (
	INTEGER   = "INTEGER"
	BOOLEAN   = "BOOLEAN"
	NULL      = "NULL"
	RETURNOBJ = "RETURN_VALUE"
	ERROROBJ  = "ERROR"
)

/* Types in Monkey
- Every value has a different representation in the Host lang, so each type of value is represented by a struct.
- Whenever interger literals are encountered in the source, they are turned into an ast.IntegerLiteral and then
	when evaluating the AST node, we turn it into an object.Integer, saving the value inside the struct.
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
type Null struct{} // Does not use Golang's nil to represent null values

// Inspect returns the value in string format
func (n *Null) Inspect() string { return "null" }

// Type returns the type's name
func (n *Null) Type() string { return NULL }

// ReturnValue is a type that wraps an "Object" value
type ReturnValue struct {
	Value Object
}

// Type returns the type's name
func (rv *ReturnValue) Type() string { return RETURNOBJ }

// Inspect returns the value in string format
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error is a type for representing all errors in Monkey lang.
type Error struct {
	Message string
}

// Type returns the type's name
func (e *Error) Type() string { return ERROROBJ }

// Inspect returns the value in string format
func (e *Error) Inspect() string { return "Error: " + e.Message }

// Environment is a type for representing the interpreter's environment.
type Environment struct {
	store map[string]Object
}

// NewEnvironment returns a pointer to a newly created Environment value.
func NewEnvironment() *Environment {
	e := &Environment{store: make(map[string]Object)}
	return e
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
