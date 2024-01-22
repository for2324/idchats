package utils

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

const (
	ignoreTagValue = "-"
	nameConnector  = "::"
)

func checkTagValidity(tagValue string) bool {
	if tagValue != "" && tagValue != ignoreTagValue {
		return true
	}
	return false
}

func isTimeType(value reflect.Value) bool {
	if _, ok := value.Interface().(time.Time); ok {
		return true
	}
	return false
}

func GetStructTag(field reflect.StructField) (tagValue string) {
	tagValue = field.Tag.Get(Tag)
	if checkTagValidity(tagValue) {
		return tagValue
	}

	tagValue = field.Tag.Get("json")
	if checkTagValidity(tagValue) {
		return strings.Split(tagValue, ",")[0]
	}

	return field.Name
}

var (
	// SkipStruct Field is returned when a struct field is skipped.
	SkipStruct = errors.New("skip struct")
)

type Typer struct {
	typ           reflect.Type
	val           reflect.Value
	fieldsMapping map[string]int
	fields        map[string]int
	name          string
	ptr           uintptr
}

var (
	Tag = "z"
)

var (
	zeroValue   = reflect.Value{}
	fieldTagMap sync.Map
	registerMap sync.Map
)

func TypeOf(obj interface{}) reflect.Type {
	return getTypElem(reflect.TypeOf(obj))
}

func ValueOf(obj interface{}) (v reflect.Value, err error) {
	v = reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		err = errors.New("not ptr")
	}
	return
}

func NewTyp(typ reflect.Type) (*Typer, error) {
	return registerValue(typ)
}

func NewVal(val reflect.Value) (*Typer, error) {
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("not ptr")
	}
	return newVal(val.Elem())
}

func newVal(val reflect.Value) (*Typer, error) {
	ot, err := registerValue(val.Type())
	if err != nil {
		return nil, err
	}
	t := *ot
	t.val = val
	return &t, nil
}

func (t *Typer) Fields() map[string]int {
	if t.fields == nil {

	}
	return t.fields
}

func (t *Typer) Field(i int) reflect.StructField {
	return t.typ.Field(i)
}

func (t *Typer) ValueOf() reflect.Value {
	return t.val
}

func (t *Typer) TypeOf() reflect.Type {
	return t.typ
}

func (t *Typer) Interface() interface{} {
	return reflect.New(t.typ).Interface()
}

func (t *Typer) CheckExistsField(name string) (int, bool) {
	if t.fieldsMapping == nil {
		return CheckExistsField(t.name, name)
	}
	i, ok := t.fieldsMapping[GetFieldTypName(t.name, name)]
	return i, ok
}

func (t *Typer) GetFieldTypName(name string) string {
	return GetFieldTypName(t.name, name)
}

var pool = sync.Pool{New: func() interface{} {
	return &Typer{}
}}

func getTyper() *Typer {
	return pool.Get().(*Typer)
}

func putTyper(t *Typer) {
	t.typ = nil
	t.fieldsMapping = nil
	t.name = ""
	pool.Put(t)
}

func (t *Typer) Name() string {
	return t.name
}

func Register(obj interface{}) error {
	typ, ok := obj.(reflect.Type)
	if !ok {
		typ = reflect.TypeOf(obj)
	}
	_, err := registerValue(typ)
	return err
}

func GetFieldTypName(typName, fieldName string) string {
	return typName + nameConnector + fieldName
}

func getTypElem(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}
func String2Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
func IsLcfirst(str string) bool {
	b := String2Bytes(str)
	if b[0] >= byte('a') && b[0] <= byte('z') {
		return true
	}
	return false
}
func registerValue(typ reflect.Type) (*Typer, error) {
	typ = getTypElem(typ)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("only registered structure")
	}
	var t *Typer
	can, typName := typ.Name() != "", ""
	if can {
		typName = typ.String()
		if v, ok := registerMap.Load(typName); ok {
			if t, ok = v.(*Typer); ok {
				return t, nil
			}
		}
	}

	t = &Typer{typ: typ, name: typName, fieldsMapping: map[string]int{}, fields: map[string]int{}}

	var register func(typ reflect.Type, name string)
	register = func(r reflect.Type, name string) {
		for i := 0; i < r.NumField(); i++ {
			field := r.Field(i)
			if IsLcfirst(field.Name) {
				continue
			}
			tag := GetStructTag(field)
			mapFieldName := GetFieldTypName(name, tag)
			typ := field.Type
			if can {
				fieldTagMap.Store(mapFieldName, i)
			}
			t.fieldsMapping[mapFieldName] = i
			t.fields[tag] = i
			if typ.Kind() == reflect.Struct {
				register(typ, mapFieldName)
			}
		}
	}
	register(typ, t.name)
	if can {
		registerMap.Store(t.name, t)
	}
	return t, nil
}
func Struct2Map(from interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(from)
	t, err := newVal(val)
	if err != nil {
		return map[string]interface{}{}, err
	}
	m := make(map[string]interface{}, len(t.fields))
	err = t.ForEachVal(func(parent []string, index int, tag string, field reflect.StructField, val reflect.Value) error {
		if len(parent) > 0 {
			return SkipStruct
		}

		m[GetStructTag(field)] = val.Interface()

		return nil
	})
	return m, err
}

func Map2Struct(from map[string]interface{}, obj interface{}) error {
	val := reflect.ValueOf(obj)
	t, err := NewVal(val)
	if err != nil {
		return err
	}

	return MapTypStruct(from, t)
}

func MapTypStruct(from map[string]interface{}, t *Typer) error {
	if t.val == zeroValue {
		return errors.New("no reflect.Value")
	}
	val := t.val
	for n := range from {
		v := from[n]
		index, exists := t.CheckExistsField(n)
		if !exists {
			continue
		}
		field := val.Field(index)
		err := SetStructFidld(t.name, n, field, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckExistsField check field is exists by name
func CheckExistsField(typeName, fieldName string) (index int, exists bool) {
	i, ok := fieldTagMap.Load(typeName + nameConnector + fieldName)
	if !ok {
		return -1, false
	}
	return i.(int), true
}

// ForEachVal For Each Struct field
func (t *Typer) ForEachVal(fn func(parent []string, index int, tag string, field reflect.StructField, val reflect.Value) error) (err error) {
	if t.val == zeroValue {
		return errors.New("no reflect.Value")
	}
	var forField func(t *Typer, v reflect.Value, parent []string) error
	forField = func(t *Typer, v reflect.Value, parent []string) error {
		for i := 0; i < t.typ.NumField(); i++ {
			field := t.Field(i)
			fieldTag := GetStructTag(field)
			if _, ok := t.CheckExistsField(fieldTag); !ok {
				continue
			}
			err = fn(parent, i, fieldTag, field, v.Field(i))
			if err == SkipStruct {
				continue
			}

			if err == nil && field.Type.Kind() == reflect.Struct {
				nt := getTyper()
				nt.typ = field.Type
				nt.fieldsMapping = t.fieldsMapping
				nt.name = t.GetFieldTypName(fieldTag)
				err = forField(nt, v.Field(i), append(parent, fieldTag))
				putTyper(nt)
			}

			if err != nil {
				return err
			}
		}
		return nil
	}

	return forField(t, t.val, []string{})
}

// ForEach For Each Struct field
func (t *Typer) ForEach(fn func(parent []string, index int, tag string, field reflect.StructField) error) (err error) {
	var forField func(t *Typer, parent []string) error
	forField = func(t *Typer, parent []string) error {
		for i := 0; i < t.typ.NumField(); i++ {
			field := t.Field(i)
			fieldTag := GetStructTag(field)
			if _, ok := t.CheckExistsField(fieldTag); !ok {
				continue
			}

			err = fn(parent, i, fieldTag, field)
			if err == SkipStruct {
				continue
			}

			if err == nil && field.Type.Kind() == reflect.Struct {
				nt := getTyper()
				nt.typ = field.Type
				nt.fieldsMapping = t.fieldsMapping
				nt.name = t.GetFieldTypName(fieldTag)
				err = forField(nt, append(parent, fieldTag))
				putTyper(nt)
			}

			if err != nil {
				return err
			}
		}
		return nil
	}

	return forField(t, []string{})
}

func SetStructFidld(typName, tag string, fValue reflect.Value, val interface{}) error {
	tp := fValue.Type()
	fkind := tp.Kind()
	if !fValue.CanSet() {
		return nil
	}
	switch fkind {
	case reflect.Bool:
		if val == nil {
			fValue.SetBool(false)
		} else if v, ok := val.(bool); ok {
			fValue.SetBool(v)
		} else {
			v := ToBool(val)
			fValue.SetBool(v)
		}
	case reflect.String:
		s, ok := val.(string)
		if !ok {
			s = ToString(val)
		}
		fValue.SetString(s)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fValue.SetInt(ToInt64(val))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fValue.SetUint(ToUint64(val))
	case reflect.Float64, reflect.Float32:
		fValue.SetFloat(ToFloat64(val))
	case reflect.Struct:
		if val == nil {
			fValue.Set(reflect.Zero(tp))
			return nil
		} else if vmap, ok := val.(map[string]interface{}); ok {
			t := &Typer{name: typName + nameConnector + tag, typ: tp, val: fValue}
			return MapTypStruct(vmap, t)
		} else if isTimeType(fValue) {
			var (
				timeString string
				timeInt    int64
			)
			switch d := val.(type) {
			case []byte:
				timeString = string(d)
			case string:
				timeString = d
			case int64:
				timeInt = d
			case int:
				timeInt = int64(d)
			}
			if timeInt > 0 {
				fValue.Set(reflect.ValueOf(time.Unix(timeInt, 0)))
			} else if timeString != "" {
				t, err := ParseTime(timeString)
				if err == nil {
					fValue.Set(reflect.ValueOf(t))
				}
			}
		} else {
			valVof := reflect.ValueOf(val)
			if valVof.Type() == tp {
				fValue.Set(valVof)
			} else if valVof.Kind() == reflect.Map {
				nv := make(map[string]interface{}, valVof.Len())
				mapKeys := valVof.MapKeys()
				for i := range mapKeys {
					m := mapKeys[i]
					nv[m.String()] = valVof.MapIndex(m).Interface()
				}
				t := &Typer{name: typName + nameConnector + tag, typ: tp, val: fValue}
				return MapTypStruct(nv, t)
			}
			// return errors.new("not support " + fkind.String())
		}
	default:
		v := reflect.ValueOf(val)
		if v.Type() == tp {
			fValue.Set(v)
		}
	}
	return nil
}
