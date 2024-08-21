package utils

import (
	"reflect"
)

// BeanCopy copies fields from source to destination
func BeanCopy(source, destination interface{}) {
	srcVal := reflect.ValueOf(source)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}
	dstVal := reflect.ValueOf(destination)
	if dstVal.Kind() != reflect.Ptr {
		panic("destination must be a pointer")
	}
	dstVal = dstVal.Elem()
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Field(i)
		if field.CanSet() {
			srcField := srcVal.FieldByName(dstVal.Type().Field(i).Name)
			if srcField.IsValid() && srcField.Type().AssignableTo(field.Type()) {
				field.Set(srcField)
			} else if srcField.IsValid() && srcField.Type().ConvertibleTo(field.Type()) {
				field.Set(srcField.Convert(field.Type()))
			}
		}
	}
}

// BeanCopyList copies fields from a list of sources to a list of destinations
func BeanCopyList(source interface{}, destination interface{}) interface{} {
	srcVal := reflect.ValueOf(source)
	dstType := reflect.TypeOf(destination).Elem()
	dstSlice := reflect.MakeSlice(reflect.SliceOf(dstType), srcVal.Len(), srcVal.Len())
	for i := 0; i < srcVal.Len(); i++ {
		dstElem := reflect.New(dstType).Elem()
		BeanCopy(srcVal.Index(i).Interface(), dstElem.Addr().Interface())
		dstSlice.Index(i).Set(dstElem)
	}
	return dstSlice.Interface()
}

// BeanCopyObject copies fields from source to destination and returns the destination
func BeanCopyObject(source, destination interface{}) interface{} {
	BeanCopy(source, destination)
	return destination
}
