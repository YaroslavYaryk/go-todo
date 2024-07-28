package helper

import "reflect"

func CheckIfInt(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Int
}
