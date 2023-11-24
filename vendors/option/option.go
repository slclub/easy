package option

import (
	"errors"
	"github.com/slclub/go-tips/safe"
	"reflect"
)

type Option struct {
	// the variable target must be a struct. if not will panic
	// target object must be a ptr, can addressable.
	target any
	// can be any type.
	config any
	// default function map
	defaultSequece []ItemAssignment
	// If reassign with zero, You need to consider whether to keep the default values
	// true : using default value
	// false : completely determined by assignment
	force bool
	// Its  priority is highest.
	finalSequece []ItemAssignment

	// Apply had been run.
	// default: false
	applyState bool
}

func OptionWith(config any) Assignment {
	opt := &Option{}
	opt.Config(config)
	return opt
}

func (self *Option) targetFields() []reflect.StructField {
	rtn := []reflect.StructField{}
	t := reflect.TypeOf(self.target).Elem()
	for i := 0; i < t.NumField(); i++ {
		rtn = append(rtn, (t.Field(i)))
	}
	return rtn
}

func (self *Option) mostSimilarMethods(mth reflect.Method) []any {

	targetfields := self.targetFields()
	for _, field := range targetfields {
		if mth.Name == field.Name {
			return []any{field, mth}
		}
		if mth.Name == firstUpper(field.Name) {
			return []any{field, mth}
		}
	}

	return nil
}

func (self *Option) mostSimilarFields(configField reflect.StructField, tcs []string) []any {
	tcss := safe.SliceString(tcs)
	targetfields := self.targetFields()
	for _, field := range targetfields {
		if tcss.In(field.Name) >= 0 {
			continue
		}
		if configField.Name == field.Name {
			return []any{field, configField}
		}
		if configField.Name == firstUpper(field.Name) {
			return []any{field, configField}
		}
	}
	return nil
}

func (self *Option) mostSimilarFieldString(fieldName string) *reflect.StructField {
	targetfields := self.targetFields()
	for _, field := range targetfields {

		if fieldName == field.Name {
			return &field
		}
		if fieldName == firstUpper(field.Name) {
			return &field
		}
	}
	return nil
}

func (self *Option) setTargetFieldWithValue(tt reflect.StructField, tv reflect.Value, av any) error {
	v := reflect.ValueOf(av)
	if tv.Kind().String() != "string" && v.Kind().String() == "string" {
		panic(any(errors.New("field:" + tt.Name + " is type " + tt.Type.Name() + " value: is type " + v.Kind().String())))
	}

	if tt.IsExported() {
		//tv.Set(v)
		//log.Log().Print("------ default setting", tt.Name)
		self.setReflectValue(tv, v)
		return nil
	}
	// changing unexported field
	untv := GetPtrUnExportFiled(self.target, tt.Name)
	//untv.Set(v)
	self.setReflectValue(untv, v)
	return nil
}

func (self *Option) setTargetFieldValueWithMethod(tcs []any) {
	if len(tcs) == 0 {
		return
	}
	targetValue := reflect.ValueOf(self.target).Elem()
	field, _ := tcs[0].(reflect.StructField)
	config, _ := tcs[1].(reflect.Method)
	mvalue := reflect.ValueOf(self.config).MethodByName(config.Name).Call(nil)
	if field.IsExported() {
		value := targetValue.FieldByName(field.Name)
		//value.Set(mvalue[0])
		//log.Log().Print("------ struct method", field.Name, mvalue[0])
		self.setReflectValue(value, mvalue[0])
		return
	}
	// changing unexported field
	tv := GetPtrUnExportFiled(self.target, field.Name)
	//tv.Set(mvalue[0])
	self.setReflectValue(tv, mvalue[0])
}

func (self *Option) setTargetFieldValueWithField(tcs []any) {
	if len(tcs) == 0 {
		return
	}
	targetValue := reflect.ValueOf(self.target).Elem()
	field, _ := tcs[0].(reflect.StructField)
	config, _ := tcs[1].(reflect.StructField)
	if field.IsExported() {
		value := targetValue.FieldByName(field.Name)
		//value.Set(self.configValue().FieldByName(config.Name))
		//log.Log().Print("------ struct field", field.Name)
		self.setReflectValue(value, self.configValue().FieldByName(config.Name))
		return
	}
	//  unexported field changing
	tv := GetPtrUnExportFiled(self.target, field.Name)
	//tv.Set(self.configValue().FieldByName(config.Name))
	self.setReflectValue(tv, self.configValue().FieldByName(config.Name))
}

func (self *Option) configValue() reflect.Value {
	v := reflect.ValueOf(self.config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func (self *Option) configType() reflect.Type {
	v := reflect.TypeOf(self.config)
	if v == nil {
		return v
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func (self *Option) getAllExportedConfigMethods() []reflect.Method {
	ts := reflect.TypeOf(self.config)
	methods := []reflect.Method{}
	for i := 0; i < ts.NumMethod(); i++ {
		if ts.Method(i).IsExported() {
			methods = append(methods, ts.Method(i))
		}
	}
	return methods
}

func (self *Option) getAllExportedConfigFields() []reflect.StructField {
	ts := self.configType()
	fields := []reflect.StructField{}
	for i := 0; i < ts.NumField(); i++ {
		if ts.Field(i).IsExported() {
			fields = append(fields, ts.Field(i))
		}
	}
	return fields
}

func (self *Option) Target(v any) {
	self.target = v
}

func (self *Option) Config(v any) Assignment {
	self.config = v
	return self
}

func (self *Option) Apply() {
	if self.applyState {
		return
	}
	defer func() { self.applyState = true }()
	// first assign default value please
	self.defaultApply(self.defaultSequece)
	defer self.defaultApply(self.finalSequece)
	// assign the values of user done.
	if self.configType() == nil {
		return
	}
	switch self.configType().Kind().String() {
	case "struct":
		self.structApply()
	}
}

func (self *Option) structApply() {
	// set the fields of object of target with methods of config.
	methods := self.getAllExportedConfigMethods()
	//mresults := [][]any{}
	usedFields := []string{}
	for _, m := range methods {
		result := self.mostSimilarMethods(m)
		if len(result) == 0 {
			continue
		}
		//mresults = append(mresults, result)
		f, _ := result[0].(reflect.StructField)
		usedFields = append(usedFields, f.Name)
		self.setTargetFieldValueWithMethod(result)
	}

	// set the fields of object of target with fields of config.
	configFields := self.getAllExportedConfigFields()
	//fresutls := [][]any{}
	for _, f := range configFields {
		result := self.mostSimilarFields(f, usedFields)
		if len(result) == 0 {
			continue
		}
		self.setTargetFieldValueWithField(result)
		//fresutls = append(fresutls, result)
	}
}

func (self *Option) defaultApply(sequece []ItemAssignment) {
	targetValue := reflect.ValueOf(self.target).Elem()
	targetType := reflect.TypeOf(self.target).Elem()
	for _, item := range sequece {
		k, v := item.Apply()
		similar := self.mostSimilarFieldString(k)
		if similar == nil {
			continue
		}
		value := targetValue.FieldByName(similar.Name)
		ttype, _ := targetType.FieldByName(similar.Name)
		self.setTargetFieldWithValue(ttype, value, v)
	}
}

func (self *Option) Default(assignments ...ItemAssignment) Assignment {
	if len(assignments) == 0 {
		return self
	}
	for _, handle := range assignments {
		if handle == nil {
			continue
		}
		if fn, ok := handle.(OptionFunc); ok {
			if reflect.ValueOf(fn).Pointer() == reflect.ValueOf(DEFAULT_IGNORE_ZERO).Pointer() {
				self.force = true
				continue
			}
		}
		self.defaultSequece = append(self.defaultSequece, handle)
	}
	return self
}

func (self *Option) Final(assignments ...ItemAssignment) Assignment {
	if len(assignments) == 0 {
		return self
	}
	self.finalSequece = append(self.finalSequece, assignments...)
	return self
}

func (self *Option) setReflectValue(target reflect.Value, value reflect.Value) {
	//log.Log().Print("----------- target:", target.Kind().String(), "  value-type:", value.Kind().String(), " =", value)

	switch target.Kind() {
	case reflect.String:
		if self.force && value.IsZero() {
			return
		}
		target.SetString(value.String())

	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		if self.force && value.IsZero() {
			return
		}
		target.SetInt(value.Int())

	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		if self.force && value.IsZero() {
			return
		}
		target.SetUint(value.Uint())

	case reflect.Float64, reflect.Float32:
		if self.force && value.IsZero() {
			return
		}
		target.SetFloat(value.Float())

	default:
		target.Set(value)
	}
}
