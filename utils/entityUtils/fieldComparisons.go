package entityUtils

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// set values from one struct to other struct
// both need ptr struct
func SetFormValues(from interface{}, to interface{}, skipTheseFields, onlyIncludeTheseFields []string) {
	if len(skipTheseFields) > 0 && len(onlyIncludeTheseFields) > 0 {
		panic("Please specify only ONE OF 'skipTheseFields' or 'onlyIncludeTheseFields' for call to SetFormValues")
	}

	val := reflect.ValueOf(from)
	elm := reflect.Indirect(val)

	valTo := reflect.ValueOf(to)
	elmTo := reflect.Indirect(valTo)

	dupPanicAssertStructPtr(val)
	dupPanicAssertStructPtr(valTo)

outFor:
	for i := 0; i < elmTo.NumField(); i++ {
		toF := elmTo.Field(i)
		name := elmTo.Type().Field(i).Name

		if len(onlyIncludeTheseFields) > 0 {
			mustInclude := false
			for _, includeField := range onlyIncludeTheseFields {
				if includeField == name {
					mustInclude = true
				}
			}
			if !mustInclude {
				continue outFor
			}
		} else {
			for _, skip := range skipTheseFields {
				if skip == name {
					continue outFor
				}
			}
		}

		f := elm.FieldByName(name)
		if f.Kind() != reflect.Invalid {
			// set value if type matched
			if f.Type().String() == toF.Type().String() {
				toF.Set(f)
			} else {
				fInt := false
				fFloat := false
				switch f.Interface().(type) {
				case int, int8, int16, int32, int64:
					fInt = true
				case float32, float64:
					fFloat = true
				case uint, uint8, uint16, uint32, uint64:
				default:
					continue outFor
				}
				switch toF.Interface().(type) {
				case int, int8, int16, int32, int64:
					var v int64
					if fInt {
						v = f.Int()
					} else {
						vu := f.Uint()
						if vu > math.MaxInt64 {
							continue outFor
						}
						v = int64(vu)
					}
					if toF.OverflowInt(v) {
						continue outFor
					}
					toF.SetInt(v)
				case float32, float64:
					var v float64
					if fFloat {
						v = f.Float()
					} else {
						vu := f.String()
						v2, err := strconv.ParseFloat(vu, 32)
						if err != nil {
							continue outFor
						}
						v = v2
					}
					toF.SetFloat(v)
				case uint, uint8, uint16, uint32, uint64:
					var v uint64
					if fInt {
						vu := f.Int()
						if vu < 0 {
							continue outFor
						}
						v = uint64(vu)
					} else {
						v = f.Uint()
					}
					if toF.OverflowUint(v) {
						continue outFor
					}
					toF.SetUint(v)
				}
			}
		}
	}
}

type ChangedField struct {
	FieldName string
	OldValue  string
	NewValue  string
}

// compare field values between two struct pointer
// return changed field names
func FormChanges(base interface{}, modified interface{}, skipTheseFields, onlyIncludeTheseFields []string) (fields []ChangedField) {
	if len(skipTheseFields) > 0 && len(onlyIncludeTheseFields) > 0 {
		panic("Please specify only ONE OF 'skipTheseFields' or 'onlyIncludeTheseFields' for call to SetFormValues")
	}

	val := reflect.ValueOf(base)
	elm := reflect.Indirect(val)

	valMod := reflect.ValueOf(modified)
	elmMod := reflect.Indirect(valMod)

	dupPanicAssertStructPtr(val)
	dupPanicAssertStructPtr(valMod)

outerForLoop:
	for i := 0; i < elmMod.NumField(); i++ {
		modF := elmMod.Field(i)
		name := elmMod.Type().Field(i).Name

		f := elm.FieldByName(name)
		if f.Kind() == reflect.Invalid {
			continue
		}

		if len(onlyIncludeTheseFields) > 0 {
			mustInclude := false
			for _, includeField := range onlyIncludeTheseFields {
				if includeField == name {
					mustInclude = true
				}
			}
			if !mustInclude {
				continue outerForLoop
			}
		} else {
			for _, skip := range skipTheseFields {
				if skip == name {
					continue outerForLoop
				}
			}
		}

		fTOrig, success := elm.Type().FieldByName(name)
		if success {
			for _, v := range strings.Split(fTOrig.Tag.Get("orm"), ";") {
				v = strings.TrimSpace(v)
				if v == "auto_now_add" {
					continue outerForLoop
				}
				if v == "auto_now" {
					continue outerForLoop
				}
				if v == "-" {
					continue outerForLoop
				}
			}
		}

		fTMod := elmMod.Type().Field(i)

		for _, v := range strings.Split(fTMod.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outerForLoop
			}
		}

		// compare two values use string
		oldValStr := ToStr(f.Interface())
		newValStr := ToStr(modF.Interface())
		if newValStr != oldValStr {
			fields = append(fields, ChangedField{FieldName: name, OldValue: oldValStr, NewValue: newValStr})
		}
	}

	return
}

type FieldWithValue struct {
	FieldName  string
	FieldValue string
}

type IgnoreFieldTypes struct {
	//Which field types to ignore when using GetAllFieldsOfStruct
	OrmSkipped  bool
	Created     bool
	Updated     bool
	Relationals bool //All fields that have orm "reverse(..." or "rel(..."
}

func GetAllFieldsOfStruct(structObj interface{}, ignoreSettings *IgnoreFieldTypes) (fieldNames []FieldWithValue) {
	val := reflect.ValueOf(structObj)
	elm := reflect.Indirect(val)

	dupPanicAssertStructPtr(val)

outerForLoop:
	for i := 0; i < elm.NumField(); i++ {
		fT := elm.Type().Field(i)
		name := fT.Name

		for _, v := range strings.Split(fT.Tag.Get("orm"), ";") {
			v = strings.TrimSpace(v)
			if ignoreSettings.OrmSkipped && v == "-" {
				continue outerForLoop
			}
			if ignoreSettings.Created && v == "auto_now_add" {
				continue outerForLoop
			}
			if ignoreSettings.Updated && v == "auto_now" {
				continue outerForLoop
			}
			if ignoreSettings.Relationals &&
				(strings.HasPrefix(v, "rel(") || strings.HasPrefix(v, "reverse(")) {
				continue outerForLoop
			}
		}

		f := elm.FieldByName(name)
		if f.Kind() == reflect.Invalid {
			continue outerForLoop
		}

		if f.Type().Kind() == reflect.Ptr {
			val2 := reflect.ValueOf(f.Interface())
			elm2 := reflect.Indirect(val2)
			if !elm2.IsValid() {
				continue outerForLoop
			}
			for j := 0; j < elm2.NumField(); j++ {
				fT2 := elm2.Type().Field(j)
				name2 := fT2.Name
				if name2 != "Id" {
					continue
				}
				fieldVal2 := ToStr(elm2.FieldByName(name2).Interface())
				fieldNames = append(fieldNames, FieldWithValue{FieldName: name, FieldValue: fieldVal2})
				continue outerForLoop
			}
			continue outerForLoop
		}

		fieldVal := ToStr(f.Interface())
		fieldNames = append(fieldNames, FieldWithValue{FieldName: name, FieldValue: fieldVal})
	}

	return
}

func GetAllFieldNamesOfStruct(structObj interface{}, ignoreSettings *IgnoreFieldTypes) []string {
	allFields := GetAllFieldsOfStruct(structObj, ignoreSettings)

	fieldNames := []string{}
	for _, fieldObj := range allFields {
		fieldNames = append(fieldNames, fieldObj.FieldName)
	}

	return fieldNames
}

func DoesStructContainField(structObj interface{}, fieldName string) bool {
	val := reflect.ValueOf(structObj)
	elm := reflect.Indirect(val)

	dupPanicAssertStructPtr(val)

	for i := 0; i < elm.NumField(); i++ {
		tmpFieldName := elm.Type().Field(i).Name

		if strings.EqualFold(tmpFieldName, fieldName) {
			return true
		}
	}

	return false
}

// assert an object must be a struct pointer
func dupPanicAssertStructPtr(val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		return
	}
	panic(fmt.Errorf("%s must be a struct pointer", val.Type().Name()))
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}
