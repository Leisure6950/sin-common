package sinobject

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

type dcInterface interface {
	DeepCopy() interface{}
}

func DeepClone(src interface{}) interface{} {
	if src == nil {
		return nil
	}
	original := reflect.ValueOf(src)
	cpy := reflect.New(original.Type()).Elem()
	copyRecursive(original, cpy)

	return cpy.Interface()
}
func elem(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}
func copyRecursive(src, dst reflect.Value) error {
	//fmt.Printf("src_type:%v, dst_type:%v\n", src.Type(), dst.Type())
	if src.Type() != dst.Type() {
		return errors.New("src and dst type is different")
	}
	if src.CanInterface() {
		if copier, ok := src.Interface().(dcInterface); ok {
			dst.Set(reflect.ValueOf(copier.DeepCopy()))
			return nil
		}
	}

	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			dst.Set(src)
			return nil
		}
		originalValue := src.Elem()

		if !originalValue.IsValid() {
			return errors.Errorf("src value is not valid")
		}
		if dst.IsNil() {
			dst.Set(reflect.New(originalValue.Type()))
			//fmt.Println("nil ptr")
		} else {
			//fmt.Println("tru ptr")
		}
		return copyRecursive(originalValue, dst.Elem())

	case reflect.Interface:
		if src.IsNil() {
			dst.Set(src)
			return nil
		}
		originalValue := src.Elem()
		if dst.IsNil() {
			copyValue := reflect.New(originalValue.Type()).Elem()
			err := copyRecursive(originalValue, copyValue)
			if err != nil {
				return err
			}
			dst.Set(copyValue)
			//fmt.Println("nil interface")
		} else {
			//fmt.Println("tru interface")
			return copyRecursive(originalValue, dst.Elem())
		}
	case reflect.Struct:
		t, ok := src.Interface().(time.Time)
		if ok {
			dst.Set(reflect.ValueOf(t))
			return nil
		}
		for i := 0; i < src.NumField(); i++ {
			if src.Type().Field(i).PkgPath != "" {
				continue
			}
			err := copyRecursive(src.Field(i), dst.Field(i))
			if err != nil {
				continue
			}
		}
	case reflect.Slice:
		if src.IsNil() {
			dst.Set(src)
			return nil
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			copyRecursive(src.Index(i), dst.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			dst.Set(src)
			return nil
		}
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			originalValue := src.MapIndex(key)
			copyValue := reflect.New(originalValue.Type()).Elem()
			copyRecursive(originalValue, copyValue)
			copyKey := DeepClone(key.Interface())
			dst.SetMapIndex(reflect.ValueOf(copyKey), copyValue)
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		} else {
			return errors.Errorf("unknown can set src")
		}
	}
	return nil
}

func DeepCopy(dst, src interface{}) error {
	return copyRecursive(reflect.ValueOf(src), reflect.ValueOf(dst))
}

func equalRecursive(src, dst reflect.Value, fieldName string) (bool, string) {
	//fmt.Printf("src_type:%v, dst_type:%v\n", src.Type(), dst.Type())
	if !src.IsValid() {
		if !dst.IsValid() {
			return true, fieldName
		} else {
			return false, fieldName
		}
	} else {
		if !dst.IsValid() {
			return false, fieldName
		}
	}

	if src.Type() != dst.Type() {
		return false, "different type"
	}
	if src.Kind() != dst.Kind() {
		return false, "different kind"
	}

	switch src.Kind() {
	case reflect.Ptr, reflect.Interface:
		if src.IsNil() {
			if dst.IsNil() {
				return true, fieldName
			} else {
				return false, fieldName
			}
		} else {
			if dst.IsNil() {
				return false, fieldName
			}
		}
		if src.Interface() == dst.Interface() {
			fmt.Printf("same pointer\n")
			return true, fieldName
		}
		orgValue := src.Elem()
		orgDstValue := dst.Elem()

		if !orgValue.IsValid() {
			if !orgDstValue.IsValid() {
				return true, fieldName
			} else {
				return false, fieldName
			}
		} else {
			if !orgDstValue.IsValid() {
				return false, fieldName
			}
		}
		return equalRecursive(orgValue, orgDstValue, fieldName)
	case reflect.Struct:
		t, ok := src.Interface().(time.Time)
		if ok {
			t2, ok := dst.Interface().(time.Time)
			if !ok {
				return false, fieldName
			}
			if t == t2 {
				return true, fieldName
			}
			return false, fieldName
		}
		for i := 0; i < src.NumField(); i++ {
			if src.Type().Field(i).PkgPath != "" {
				continue
			}
			equals, fn := equalRecursive(src.Field(i), dst.Field(i), fieldName+"."+src.Type().Field(i).Name)
			if !equals {
				return equals, fn
			}
		}
		return true, fieldName
	case reflect.Slice:
		if src.IsNil() {
			if dst.IsNil() {
				return true, fieldName
			} else {
				return false, fieldName
			}
		} else {
			if dst.IsNil() {
				return false, fieldName
			}
		}
		if src.Len() != dst.Len() {
			return false, fieldName
		}
		for i := 0; i < src.Len(); i++ {
			equals, fn := equalRecursive(src.Index(i), dst.Index(i), fieldName+"["+strconv.FormatInt(int64(i), 10)+"]")
			if !equals {
				return equals, fn
			}
		}
		return true, fieldName
	case reflect.Map:
		if src.IsNil() {
			if dst.IsNil() {
				return true, fieldName
			} else {
				return false, fieldName
			}
		} else {
			if dst.IsNil() {
				return false, fieldName
			}
		}
		if len(src.MapKeys()) != len(dst.MapKeys()) {
			return false, fieldName
		}
		for _, key := range src.MapKeys() {
			originalValue := src.MapIndex(key)
			orgDstValue := dst.MapIndex(key)
			equals, fn := equalRecursive(originalValue, orgDstValue, fieldName+"."+key.String())
			if !equals {
				return equals, fn
			}
		}
		return true, fieldName
	default:
		if src.Interface() == dst.Interface() {
			return true, fieldName
		}
	}
	return false, fieldName
}
func DeepEquals(v1, v2 interface{}) (bool, string) {
	return equalRecursive(reflect.ValueOf(v1), reflect.ValueOf(v2), "$")
}
