package deepcopy

import (
	"fmt"
	"reflect"
	"unsafe"
)

// FromTo deep copies original and assigns the copy to the copy argument (pointer).
func FromTo(original, copy interface{}) error {
	if original == nil {
		copy = nil
		return nil
	} else if copy == nil { // TODO try to initialize it here
		return fmt.Errorf("FromTo: copy target is nil, it should be a valid pointer")
		// copyValue := reflect.New(value.Type().Elem()).Elem()
		// copy = copyValue.Interface()
	}
	copyValue := reflect.ValueOf(copy)
	if copyValue.Kind() != reflect.Ptr {
		return fmt.Errorf("FromTo: copy target type %T and not a pointer", copy)
	}
	value := reflect.ValueOf(original)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			copy = nil // TODO return typed nil
			return nil
		}
		value = value.Elem()
	}
	copyValue.Elem().Set(deepCopy(value))
	return nil
}

func deepCopy(original reflect.Value) reflect.Value {
	switch original.Kind() {
	case reflect.Slice:
		return deepCopySlice(original)
	case reflect.Map:
		return deepCopyMap(original)
	case reflect.Ptr:
		return deepCopyPointer(original)
	case reflect.Struct:
		return deepCopyStruct(original)
	case reflect.Chan:
		return deepCopyChan(original)
	case reflect.Array:
		return deepCopyArray(original)
	default:
		return forceCopyValue(original)
	}
}

func deepCopySlice(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	_copy := reflect.MakeSlice(original.Type(), 0, 0)
	for i := 0; i < original.Len(); i++ {
		elementCopy := deepCopy(original.Index(i))
		_copy = reflect.Append(_copy, elementCopy)
	}
	return _copy
}

func deepCopyArray(original reflect.Value) reflect.Value {
	if original.Len() == 0 {
		// it cannot be changed anyway, so we can return the original
		return original
	}
	elementType := original.Index(0).Type()
	arrayType := reflect.ArrayOf(original.Len(), elementType)
	newPointer := reflect.New(arrayType)
	_copy := newPointer.Elem()
	for i := 0; i < original.Len(); i++ {
		subCopy := deepCopy(original.Index(i))
		_copy.Index(i).Set(subCopy)
	}
	return _copy
}

func deepCopyMap(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	keyType := original.Type().Key()
	valueType := original.Type().Elem()
	mapType := reflect.MapOf(keyType, valueType)
	_copy := reflect.MakeMap(mapType)
	for _, key := range original.MapKeys() {
		value := deepCopy(original.MapIndex(key))
		_copy.SetMapIndex(key, value)
	}
	return _copy
}

func deepCopyPointer(original reflect.Value) reflect.Value {
	if original.IsNil() {
		return original
	}
	element := original.Elem()
	_copy := reflect.New(element.Type())
	copyElement := deepCopy(element)
	_copy.Elem().Set(copyElement)
	return _copy
}

func deepCopyStruct(original reflect.Value) reflect.Value {
	_copy := reflect.New(original.Type()).Elem()
	_copy.Set(original)
	for i := 0; i < original.NumField(); i++ {
		fieldValue := _copy.Field(i)
		fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
		fieldValue.Set(deepCopy(fieldValue))
	}
	return _copy
}

func deepCopyChan(original reflect.Value) reflect.Value {
	return reflect.MakeChan(original.Type(), original.Cap())
}

// forceCopyValue simply creates a new pointer and sets its value to the original.
func forceCopyValue(original reflect.Value) reflect.Value {
	originalType := original.Type()
	newPointer := reflect.New(originalType)
	newPointer.Elem().Set(original)
	return newPointer.Elem()
}
