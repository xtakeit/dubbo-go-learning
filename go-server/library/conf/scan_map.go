package conf

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// scan 将map src中的信息扫描到结构体指针dst, 扫描时src的key与结构体成员的tag对应的值相对应,
// 结构体类型必须是非nil的结构体指针, 且结构体成员拥有指定tag时src必须存在对应的key, 否则将返回错误
func scan(src map[string]string, dst interface{}, tag string) (err error) {
	rv, err := muststptr(dst)
	if err != nil {
		err = fmt.Errorf("must struct pointer: %w", err)
		return
	}

	for i := 0; i < rv.NumField(); i++ {
		key := strings.Split(rv.Type().Field(i).Tag.Get(tag), ",")[0]

		if key == "" {
			continue
		}

		val := strings.TrimSpace(src[key])
		if val == "" {
			err = fmt.Errorf("field %s is not present is source", key)
			return
		}

		fv := rv.Field(i)

		switch fv.Kind() {
		default:
			return fmt.Errorf("unsurpported field type: %s", fv.Kind())
		case reflect.Bool:
			boolVal := strings.ToLower(val) != "false" && val != "" && val != "0"
			fv.SetBool(boolVal)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return fmt.Errorf(
					"parse int: %w", err,
				)
			}
			fv.SetInt(intVal)
		case reflect.String:
			fv.SetString(val)
		}
	}

	return nil
}

// muststptr 检测传入的接口值v是否是非空的结构体指针,
// 如果是则返回该指针指向的结构的反射值rv和nil, 否则将返回不为空的错误err
func muststptr(v interface{}) (rv reflect.Value, err error) {
	rv = reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		err = fmt.Errorf("interface's kind is non-pointer or interface is nil")
		return
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		err = fmt.Errorf("interface's elem kind is not struct")
		return
	}

	return
}
