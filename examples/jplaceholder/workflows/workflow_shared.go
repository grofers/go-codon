package workflows

import (
	"reflect"
)

func maybe(obj interface{}, err error) interface{} {
	if err != nil {
		return nil
	}
	return obj
}

func resolvePointers(obj interface{}) interface{} {
	if obj == nil {
		return obj
	}
	obj_val := reflect.ValueOf(obj)
	for obj_val.Kind() == reflect.Ptr {
		obj_val = obj_val.Elem()
	}
	if obj_val.CanInterface() {
		return obj_val.Interface()
	} else {
		return nil
	}
}
