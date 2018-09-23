package evaluator

import (
	"github.com/shksa/yeezy/object"
)

var builtins = map[string]object.BuiltInFunction{
	"len": func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("Wrong number of arguments. want=%d, got=%d", 1, len(args))
		}

		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}

		default:
			return newError("len doesn'nt support the given argument. got=%s", args[0].Type())
		}
	},
}
