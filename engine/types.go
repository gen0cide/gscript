package engine

import (
	"bytes"
	"fmt"
	"reflect"
)

type typeBuffer struct {
	Input  interface{}
	Orig   reflect.Value
	Copy   reflect.Value
	Buf    *bytes.Buffer
	Engine *Engine
}

func (e *Engine) newTypeBuffer(i interface{}) *typeBuffer {
	if i == nil {
		return nil
	}
	orig := reflect.ValueOf(i)
	copy := reflect.New(orig.Type()).Elem()
	t := &typeBuffer{
		Input:  i,
		Buf:    new(bytes.Buffer),
		Orig:   orig,
		Copy:   copy,
		Engine: e,
	}
	return t
}

// func translate(obj interface{}) interface{} {
// 	// Wrap the original in a reflect.Value
// 	original := reflect.ValueOf(obj)

// 	copy := reflect.New(original.Type()).Elem()
// 	translateRecursive(copy, original)

// 	// Remove the reflection wrapper
// 	return copy.Interface()
// }

func (t *typeBuffer) WalkType() (string, error) {
	t.translateRecursive(t.Copy, t.Orig, 0)
	return t.Buf.String(), nil
}

func (t *typeBuffer) Tab(i int) {
	for x := 0; x < i; x++ {
		t.Buf.WriteString("  ")
	}
}

func (t *typeBuffer) translateRecursive(copy, original reflect.Value, depth int) {
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			t.Buf.WriteString("<nil>")
			return
		}

		t.Buf.WriteString("*")

		if !originalValue.CanSet() {
			t.Buf.WriteString("<hidden>")
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		t.translateRecursive(copy.Elem(), originalValue, depth)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		t.translateRecursive(copyValue, originalValue, depth)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		t.Buf.WriteString(original.Type().String())
		t.Buf.WriteString("{")
		//numberOfFields := original.NumField()
		fieldFound := false
		for i := 0; i < original.NumField(); i++ {
			if !original.Field(i).CanSet() {
				continue
			}
			t.Buf.WriteString("\n")
			t.Tab(depth + 1)
			t.Buf.WriteString(original.Type().Field(i).Name)
			t.Buf.WriteString(": ")
			t.translateRecursive(copy.Field(i), original.Field(i), depth+1)
			t.Buf.WriteString(", ")
		}
		if fieldFound {
			t.Buf.WriteString("\n")
			t.Tab(depth)
		}
		t.Buf.WriteString("}")

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		t.Buf.WriteString(original.Type().String())
		t.Buf.WriteString(" { ")
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			t.translateRecursive(copy.Index(i), original.Index(i), depth)
			t.Buf.WriteString(", ")
		}
		t.Buf.WriteString(" } ")

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		t.Buf.WriteString(original.Type().String())
		t.Buf.WriteString(" { ")
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {

			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			t.translateRecursive(copyValue, originalValue, depth)
			copy.SetMapIndex(key, copyValue)
		}
		t.Buf.WriteString(" }")
		return

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		t.Buf.WriteString("\"")
		t.Buf.WriteString(original.Interface().(string))
		t.Buf.WriteString("\"")
		return

	// And everything else will simply be taken from the original

	case reflect.Bool:
		t.Buf.WriteString(fmt.Sprintf("%v", original.Interface().(bool)))
		return
	default:
		t.Buf.WriteString(original.String())
		return
	}

}
