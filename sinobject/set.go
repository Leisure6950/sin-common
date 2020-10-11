package sinobject

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"unicode"
)

func parseTag(tag string) string {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}
func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}
func canNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice:
		return true
	default:
		return false
	}
}
func getElem(v reflect.Value) (reflect.Value, bool) {
	for {
		//if canNil(v) {
		//	if v.IsNil() {
		//		return v, false
		//	}
		//	v = v.Elem()
		//} else {
		//	return v, true
		//}

		switch v.Kind() {
		case reflect.Ptr, reflect.Interface:
			if v.IsNil() {
				return v, false
			}
			v = v.Elem()
		default:
			return v, true
		}
	}
}
func getElemWithAutoNew(v reflect.Value) reflect.Value {
	v, ok := getElem(v)
	if !ok {
		v.Set(reflect.New(v.Type().Elem()))
		return getElemWithAutoNew(v)
	}
	return v
}
func print(msg string, v reflect.Value) {
	d, _ := json.Marshal(v.Interface())
	fmt.Println(msg, v.Type(), string(d))
}
func copyFromBaseType(dst reflect.Value, src reflect.Value) error {
	dstElem, ok1 := getElem(dst)
	srcElem, ok2 := getElem(src)
	if srcElem.Type() == dstElem.Type() {
		print("srcElem", srcElem)
		print("dstElem", dstElem)
		// 类型相同
		if err := copyRecursive(srcElem, dstElem); err != nil {
			return err
		}
		return nil
	}
	if !ok1 && !ok2 {
		return nil
	} else if ok1 && !ok2 {
		// dst 设置成空
		if canNil(dstElem) {
			dstElem.SetPointer(nil)
			return nil
		} else {
			return errors.Errorf("src is nil, but dst can't be nil")
		}
	} else if !ok1 && ok2 {
		dstElem = getElemWithAutoNew(dstElem)
	}
	switch dstElem.Kind() {
	case reflect.Struct:
		if srcElem.Kind() != reflect.Map || srcElem.Type().Key().Kind() != reflect.String {
			return errors.Errorf("dst is struct, src need be map, map key need be string")
		}
		for i := 0; i < dstElem.NumField(); i++ {
			f := dstElem.Type().Field(i)
			key, ok := f.Tag.Lookup("json")
			if ok {
				key = parseTag(key)
				if !isValidTag(key) {
					key = f.Name
				}
			} else {
				key = f.Name
			}
			v := srcElem.MapIndex(reflect.ValueOf(key))
			if !v.IsValid() {
				// 未找到
				continue
			}
			return copyFromBaseType(dstElem.Field(i), v)
		}
	case reflect.Slice:
		if srcElem.Kind() != reflect.Slice && srcElem.Kind() != reflect.Array {
			return errors.Errorf("need slice, but src is not slice")
		}
		if dstElem.Type().Elem() == srcElem.Type().Elem() {
			dstElem.Set(srcElem)
			return nil
		}
		l := srcElem.Len()
		dstElem.Set(reflect.MakeSlice(dstElem.Type(), l, l))
		for i := 0; i < l; i++ {
			if err := copyFromBaseType(dstElem.Index(i), src.Index(i)); err != nil {
				return err
			}
		}
	case reflect.Array:
	case reflect.Map:
		if srcElem.Kind() != reflect.Map {
			return errors.Errorf("dst is map, src need be map")
		}
		if dstElem.Type().Key() != srcElem.Type().Key() {
			return errors.Errorf("dst map key type is not equals src key type")
		}
		if dstElem.Type().Elem() == srcElem.Type().Elem() {
			dstElem.Set(srcElem)
			return nil
		}
		keys := dstElem.MapKeys()
		for _, k := range keys {
			srcValue := srcElem.MapIndex(k)
			if !srcValue.IsValid() {
				// 无效
				continue
			}
			if err := copyFromBaseType(dstElem.MapIndex(k), srcValue); err != nil {
				return err
			}
		}
	default:
		if dstElem.Type() != srcElem.Type() {
			return errors.Errorf("dst map value type is not equals src value type")
		}
		dstElem.Set(srcElem)
	}
	return nil
}
func setValue(dst reflect.Value, key reflect.Value, src reflect.Value) error {
	topElem, ok := getElem(dst)
	if !ok {
		return errors.Errorf("dst is null")
	}
	//keyStr := key.String()
	switch topElem.Kind() {
	case reflect.Struct:
		//for i := 0; i < topElem.NumField(); i++ {
		//	f := topElem.Type().Field(i)
		//	tag, ok := f.Tag.Lookup("json")
		//	if ok {
		//		tag = parseTag(tag)
		//		if !isValidTag(tag) {
		//			tag = f.Name
		//		}
		//	} else {
		//		tag = f.Name
		//	}
		//	if tag != keyStr && f.Name != keyStr {
		//		continue
		//	}
		//	return copyFromBaseType(topElem.Field(i), src)
		//}
		//return errors.Errorf("not found field from struct")
		return errors.Errorf("unsupported struct")
	case reflect.Slice:
		return errors.Errorf("unsupported slice dst")
	case reflect.Array:
		return errors.Errorf("unsupported array dst")
	case reflect.Map:
		if topElem.Type().Key() != key.Type() {
			return errors.Errorf("dst map key type is not same as src key type")
		}
		//return copyFromBaseType(topElem.MapIndex(key), src)
		topElem.SetMapIndex(key, src)
	default:
		return errors.Errorf("dst is not container type")
	}
	return nil
}
func SetValue(dst interface{}, key interface{}, src interface{}) error {
	return setValue(reflect.ValueOf(dst), reflect.ValueOf(key), reflect.ValueOf(src))
}
