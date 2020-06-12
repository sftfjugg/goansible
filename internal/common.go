package internal

import "reflect"

func Contain(container interface{}, item interface{})(bool){
	containerValues := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind(){
	case reflect.Slice, reflect.Array:
		for i:=0;i<containerValues.Len();i++{
			if containerValues.Index(i).Interface() == item{
				return true
			}
		}
	case reflect.Map:
		if containerValues.MapIndex(reflect.ValueOf(item)).IsValid(){
			return true
		}
	}
	return false
}
