package sinxml

import (
	"github.com/beevik/etree"
	"github.com/pkg/errors"
	"github.com/sin-z/sin-common/sinlog"
	"io"
	"reflect"
	"strconv"
)

func decode(node *etree.Element, out reflect.Value, parent reflect.Value) error {
	for out.Kind() == reflect.Ptr || out.Kind() == reflect.Interface {
		if out.IsNil() {
			return nil
		}
		out = out.Elem()
	}
	switch out.Kind() {
	case reflect.Map:
		if len(node.Child) <= 0 {
			sinlog.Debugw("not found children")
			return nil
		}
		v := out.MapIndex(reflect.ValueOf(node.Tag))
		if !v.IsValid() {
			sinlog.With("tag", node.Tag).Debugw("out value is not valid")
			return nil
		}
		childEles := node.ChildElements()
		if len(childEles) > 0 {
			for _, ele := range childEles {
				if err := decode(ele, v, out); err != nil {
					sinlog.Errorw("decode error", "err", err)
					return err
				}
			}
		} else {
			decode(node, v, out)
		}
	case reflect.Struct:
		if len(node.Child) <= 0 {
			sinlog.Debugw("not found children")
			return nil
		}
		v := out.FieldByName(node.Tag)
		if !v.IsValid() {
			sinlog.Debugw("out value is not valid")
			return nil
		}
		childEles := node.ChildElements()
		if len(childEles) > 0 {
			for _, ele := range childEles {
				if err := decode(ele, v, out); err != nil {
					sinlog.Errorw("decode error", "err", err)
					return err
				}
			}
		} else {
			decode(node, v, out)
		}
	case reflect.Slice, reflect.Array:
		return errors.Errorf("unsupport slice or array type")
	case reflect.Int:
		// value
		v, err := strconv.ParseInt(node.Text(), 10, 64)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(int(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(int(v)))
		}
	case reflect.Int64:
		v, err := strconv.ParseInt(node.Text(), 10, 64)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(v))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(v))
		}
	case reflect.Int32:
		v, err := strconv.ParseInt(node.Text(), 10, 32)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(int32(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(int32(v)))
		}
	case reflect.Int16:
		v, err := strconv.ParseInt(node.Text(), 10, 16)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(int16(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(int16(v)))
		}
	case reflect.Int8:
		v, err := strconv.ParseInt(node.Text(), 10, 8)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(int8(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(int8(v)))
		}
	case reflect.Uint:
		v, err := strconv.ParseUint(node.Text(), 10, 64)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(uint(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(uint(v)))
		}
	case reflect.Uint64:
		v, err := strconv.ParseUint(node.Text(), 10, 64)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(v))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(v))
		}
	case reflect.Uint32:
		v, err := strconv.ParseUint(node.Text(), 10, 32)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(uint32(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(uint32(v)))
		}
	case reflect.Uint16:
		v, err := strconv.ParseUint(node.Text(), 10, 16)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(uint16(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(uint16(v)))
		}
	case reflect.Uint8:
		v, err := strconv.ParseUint(node.Text(), 10, 8)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(uint8(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(uint8(v)))
		}
	case reflect.Float64:
		v, err := strconv.ParseFloat(node.Text(), 64)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(v))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(v))
		}
	case reflect.Float32:
		v, err := strconv.ParseFloat(node.Text(), 32)
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(float32(v)))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(float32(v)))
		}
	case reflect.Bool:
		v, err := strconv.ParseBool(node.Text())
		if err != nil {
			return err
		}
		if out.CanSet() {
			out.Set(reflect.ValueOf(v))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(v))
		}
	case reflect.String:
		if out.CanSet() {
			out.Set(reflect.ValueOf(node.Text()))
		} else {
			parent.SetMapIndex(reflect.ValueOf(node.Tag), reflect.ValueOf(node.Text()))
		}
	}
	return nil
}

func Unmarshal(r io.Reader, out interface{}) error {
	doc := etree.NewDocument()
	_, err := doc.ReadFrom(r)
	if err != nil {
		sinlog.Errorw("read xml data error", "err", err)
		return err
	}

	rootEle := doc.Root()
	if rootEle == nil {
		sinlog.Debugw("not found root element")
		return nil
	}
	outReflect := reflect.ValueOf(out)
	return decode(rootEle, outReflect, reflect.ValueOf(nil))
}
