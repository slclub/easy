package option

import (
	"reflect"
	"unsafe"
)

/**
 * Go 导出结构体内 非导出字段
 */

// reflect/value.go
// NewAt returns a Value representing a pointer to a value of the
// specified type, using p as that pointer.
// reflect 内部方法NewAt
/*
func NewAt(typ Type, p unsafe.Pointer) Value {
	fl := flag(Ptr)
	t := typ.(*rtype)
	return Value{t.ptrTo(), p, fl}
}
*/

func anysdfsdf() {
}

/*
	type Example struct {
	  a string
	}

var eg testData.Example
a:=GetStructPtrUnExportedField(&eg, "a").String()
*/
func GetStructPtrUnExportedField(source interface{}, fieldName string) reflect.Value {
	// 获取非导出字段反射对象
	v := reflect.ValueOf(source).Elem().FieldByName(fieldName)
	// 构建指向该字段的可寻址（addressable）反射对象
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

/**
reference:
	https://www.zhihu.com/tardis/bd/art/148231342
*/
