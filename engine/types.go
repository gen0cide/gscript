package engine

// import (
// 	"bytes"
// 	"reflect"

// 	"github.com/robertkrimen/otto"
// )

// type TypeBuffer struct {
// 	Input interface{}
// 	Orig  reflect.Value
// 	Copy  reflect.Value
// 	Buf   *bytes.Buffer
// }

// func NewTypeBuffer(i interface{}) *TypeBuffer {
// 	orig := reflect.ValueOf(i)
// 	copy := reflect.New(orig.Type()).Elem()
// 	t := &TypeBuffer{
// 		Input: i,
// 		Buf:   new(bytes.Buffer),
// 		Orig:  orig,
// 		Copy:  copy,
// 	}
// }

// func (e *Engine) vmTypeChecker(call otto.FunctionCall) otto.Value {
// 	if len(call.ArgumentList) == 0 {
// 		return otto.UndefinedValue()
// 	} else if len(call.ArgumentList) == 1 {
// 		val, err := call.Argument(0).Export()
// 		if err != nil {
// 			return e.Raise("jsexport", "could not convert argument number 0")
// 		}
// 		decType := reflect.TypeOf(val)
// 		retObj, err := e.VM.ToValue(parseComplexType(decType))
// 		if err != nil {
// 			return e.Raise("jsimport", "could not convert return value")
// 		}
// 		return retObj
// 	}
// 	typeArray := []string{}
// 	for idx, a := range call.ArgumentList {
// 		val, err := a.Export()
// 		if err != nil {
// 			return e.Raise("jsexport", "could not convert argument number %d", idx)
// 		}
// 		decType := reflect.TypeOf(val)
// 		typeArray = append(typeArray, parseComplexType(decType))
// 	}
// 	retObj, err := e.VM.ToValue(typeArray)
// 	if err != nil {
// 		return e.Raise("jsimport", "could not convert return array")
// 	}
// 	return retObj
// }

// func translate(obj interface{}) interface{} {
// 	// Wrap the original in a reflect.Value
// 	original := reflect.ValueOf(obj)

// 	copy := reflect.New(original.Type()).Elem()
// 	translateRecursive(copy, original)

// 	// Remove the reflection wrapper
// 	return copy.Interface()
// }

// func (t *TypeBuffer) WalkType() (string, error) {
// 	t.translateRecursive(t.Copy, t.Orig)
// }

// func (t *TypeBuffer) translateRecursive(copy, original reflect.Value) {
// 	switch original.Kind() {
// 	// The first cases handle nested structures and translate them recursively

// 	// If it is a pointer we need to unwrap and call once again
// 	case reflect.Ptr:
// 		t.Buf.WriteString("")
// 		originalValue := original.Elem()
// 		// Check if the pointer is nil
// 		if !originalValue.IsValid() {
// 			return
// 		}
// 		// Allocate a new object and set the pointer to it
// 		copy.Set(reflect.New(originalValue.Type()))
// 		// Unwrap the newly created pointer
// 		translateRecursive(copy.Elem(), originalValue)

// 	// If it is an interface (which is very similar to a pointer), do basically the
// 	// same as for the pointer. Though a pointer is not the same as an interface so
// 	// note that we have to call Elem() after creating a new object because otherwise
// 	// we would end up with an actual pointer
// 	case reflect.Interface:
// 		// Get rid of the wrapping interface
// 		originalValue := original.Elem()
// 		// Create a new object. Now new gives us a pointer, but we want the value it
// 		// points to, so we have to call Elem() to unwrap it
// 		copyValue := reflect.New(originalValue.Type()).Elem()
// 		translateRecursive(copyValue, originalValue)
// 		copy.Set(copyValue)

// 	// If it is a struct we translate each field
// 	case reflect.Struct:
// 		for i := 0; i < original.NumField(); i += 1 {
// 			translateRecursive(copy.Field(i), original.Field(i))
// 		}

// 	// If it is a slice we create a new slice and translate each element
// 	case reflect.Slice:
// 		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
// 		for i := 0; i < original.Len(); i += 1 {
// 			translateRecursive(copy.Index(i), original.Index(i))
// 		}

// 	// If it is a map we create a new map and translate each value
// 	case reflect.Map:
// 		copy.Set(reflect.MakeMap(original.Type()))
// 		for _, key := range original.MapKeys() {
// 			originalValue := original.MapIndex(key)
// 			// New gives us a pointer, but again we want the value
// 			copyValue := reflect.New(originalValue.Type()).Elem()
// 			translateRecursive(copyValue, originalValue)
// 			copy.SetMapIndex(key, copyValue)
// 		}

// 	// Otherwise we cannot traverse anywhere so this finishes the the recursion

// 	// If it is a string translate it (yay finally we're doing what we came for)
// 	case reflect.String:
// 		translatedString := dict[original.Interface().(string)]
// 		copy.SetString(translatedString)

// 	// And everything else will simply be taken from the original
// 	default:
// 		copy.Set(original)
// 	}

// }
