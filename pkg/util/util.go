package util

import (
	"reflect"
)

func StructToStrMap(i interface{}, len int) map[string]string {
	mapStr := make(map[string]string, len)
	iVal := reflect.ValueOf(i).Elem()
	typ := iVal.Type()
	for i := 0; i < iVal.NumField(); i++ {
		f := iVal.Field(i)
		mapStr[typ.Field(i).Name] = f.String()
	}
	return mapStr
}
