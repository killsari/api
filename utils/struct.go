package utils

import (
	"reflect"
	"strings"
)

func Struct2MapWithLowerKey(obj interface{}) *map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return &data
}

func Struct2Map(obj interface{}) *map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return &data
}
