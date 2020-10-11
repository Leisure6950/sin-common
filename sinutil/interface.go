package sinutil

import "reflect"

// 切片型接口转为切片接口
func Interface2SliceInterface(i interface{}) []interface{} {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice {
		panic("i is not slice type")
	}
	len := v.Len()
	si := make([]interface{}, len)
	for i := 0; i < len; i++ {
		si[i] = v.Index(i).Addr().Interface()
	}
	return si
}
