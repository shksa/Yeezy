package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/shksa/yeezy/ast"
)

/* Why Object?
- Object is an abstraction for "data" in the lang.
- An object will have a type and a value.
- All the primitive and composite types of yeezy are the various "kinds" of the objects.
- The "types" of yeezy are defined by various custom structure types in Go that implement the Object interface.
- An Integer type in yeezy is an Object with a particular type name called "INTEGER" and can take all 64-bit integer values
	of Golang.
- A Boolean type in yeezy is an Object with a particular type name called "BOOLEAN" can take either the true or false boolean
	values of Golang.
*/

// Object is an interface which is implemented by all objects in yeezy lang.
type Object interface {
	Inspect() string // returns value of the object in string format
	Type() string    // returns type of the object
}

// list of all the types of Objects in yeezy
const (
	INTEGER         = "INTEGER"
	BOOLEAN         = "BOOLEAN"
	STRING          = "STRING"
	NULL            = "NULL"
	RETURNOBJ       = "RETURN_VALUE"
	ERROROBJ        = "ERROR"
	FUNCTION        = "FUNCTION"
	BUILTINFUNCTION = "BUILTIN_FUNCTION"
)

/* Types in yeezy
- Every value has a different representation in the Host lang, so each type of value is represented by a struct.
- Whenever interger literals are encountered in the source, they are turned into an ast.IntegerLiteral and then
	when evaluating the AST node, we turn it into an object.Integer, saving the value inside the struct.
*/

// Integer is type for representing all integer literal objects in the yeezy lang.
type Integer struct {
	Value int64
}

// Inspect returns the value in string format
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns the type's name
func (i *Integer) Type() string { return INTEGER }

// Boolean is type for representing all boolean literal objects in the yeezy lang.
type Boolean struct {
	Value bool
}

// Inspect returns the value in string format
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Type returns the type's name
func (b *Boolean) Type() string { return BOOLEAN }

// String is a type for representing all string literal objects in the yeezy lang.
type String struct {
	Value string
}

// Inspect returns the value in string format
func (s *String) Inspect() string { return s.Value }

// Type returns the type's name
func (s *String) Type() string { return STRING }

// Null is type for representing the absence of values in yeezy
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

// Error is a type for representing all errors in yeezy lang.
type Error struct {
	Message string
}

// Type returns the type's name
func (e *Error) Type() string { return ERROROBJ }

// Inspect returns the value in string format
func (e *Error) Inspect() string { return "Error: " + e.Message }

// Environment is a type for representing the interpreter's environment.
type Environment struct {
	store    map[string]Object
	outerEnv *Environment
}

// NewEnvironment returns a pointer to a newly created Environment value.
func NewEnvironment() *Environment {
	e := &Environment{store: make(map[string]Object), outerEnv: nil}
	return e
}

// Get returns the object mapped to a identifier in the current environemnt or in any of the enclosing environments.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outerEnv != nil { // If the current env has an outer env, that outer env is checked for the binding.
		obj, ok = e.outerEnv.Get(name) // Recursively checks all the enclosing envs of the current env to find the binding.
	}
	return obj, ok
}

// Set maps a identifier name to an object
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// Function is a type for representing all the function literal values in yeezy.
type Function struct {
	Parameters []*ast.IdentifierNode
	Body       *ast.BlockStatementNode
	Env        *Environment // functions carry their environment with them
}

// Type returns the type's name
func (f *Function) Type() string { return FUNCTION }

// Inspect returns the value in string format
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(f.Body.String())

	return out.String()
}

// NewEnclosedEnvironment return a new environment that extends the current enclosing environment.
func NewEnclosedEnvironment(outerEnv *Environment) *Environment {
	newEnv := NewEnvironment()
	newEnv.outerEnv = outerEnv
	return newEnv
}

// BuiltInFunction is a type for all representing all the built-in functions, it is an object that is exposed to the users of yeezy
type BuiltInFunction func(...Object) Object

// Inspect returns the value in string format
func (bf BuiltInFunction) Inspect() string { return "built-in function" }

// Type returns the type's name
func (bf BuiltInFunction) Type() string { return BUILTINFUNCTION }
