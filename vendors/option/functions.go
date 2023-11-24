package option

import (
	"reflect"
	"strings"
	"unsafe"
)

// -----------------------------------------------------------------
// common functions
// -----------------------------------------------------------------

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetPtrUnExportFiled(s interface{}, filed string) reflect.Value {
	v := reflect.ValueOf(s).Elem().FieldByName(filed)
	// 必须要调用 Elem()
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}

//
//func SetPtrUnExportFiled(s interface{}, filed string, val interface{}) error {
//	v := GetPtrUnExportFiled(s, filed)
//	rv := reflect.ValueOf(val)
//	if v.Kind() != v.Kind() {
//		return fmt.Errorf("invalid kind, expected kind: %v, got kind:%v", v.Kind(), rv.Kind())
//	}
//
//	v.Set(rv)
//	return nil
//}
