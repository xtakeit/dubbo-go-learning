package mysql

import (
	"reflect"
	"strings"
)

const defaultColTag = "db"

var colTag = defaultColTag

func SetColTag(tag string) {
	colTag = tag
}

// getcolmp 获取字段与结构体成员下标映射
func getcolmp(t reflect.Type) (colmp map[string]int) {
	colmp = make(map[string]int, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		col := strings.ToLower(t.Field(i).Name)

		if tagv := t.Field(i).Tag.Get(colTag); tagv != "" {
			col = strings.Split(tagv, ",")[0]
		}

		colmp[col] = i
	}

	return
}

// getsts 获取扫描对象指针列表
func getsts(rp reflect.Value, cols []string, colmp map[string]int) (sts []interface{}) {
	rv := rp.Elem()
	for _, col := range cols {
		fi := rv.Field(colmp[col]).Addr().Interface()
		sts = append(sts, fi)
	}

	return
}

// isstlist 扫描对象必须是指向slice的非空指针, 且slice元素类型必须是结构体的指针
func isstlist(v interface{}) (ok bool) {
	if reflect.ValueOf(v).IsNil() {
		return
	}

	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Ptr {
		return
	}

	rt = rt.Elem()
	if rt.Kind() != reflect.Slice {
		return
	}

	rt = rt.Elem()
	if rt.Kind() != reflect.Ptr {
		return
	}

	rt = rt.Elem()
	if rt.Kind() != reflect.Struct {
		return
	}

	ok = true
	return
}

// isstrecord 扫描对象必须是指向结构体的非空指针
func isstrecord(v interface{}) (ok bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return
	}

	ok = true
	return
}
