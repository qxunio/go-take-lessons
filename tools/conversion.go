package tools

import (
	"errors"
	"go-take-lessons/domain/comm"
	"reflect"
	"strconv"
	"time"
)

func StringToInt64(str string) int64 {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return i
	}
	println("String To int64 err " + str)
	return 1
}

func StringToUint8(str string) uint8 {
	if i, err := strconv.ParseUint(str, 10, 64); err == nil {
		return uint8(i)
	}
	println("String To uint8 err " + str)
	return 1
}

func StringToInt(str string) int {
	if i, err := strconv.ParseInt(str, 10, 64); err == nil {
		return int(i)
	}
	println("String To int err " + str)
	return 1
}

func SliceCut(slicePtr interface{}, field string) (map[string][]map[string]string, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(slicePtr))
	if sliceValue.Kind() != reflect.Slice {
		return nil, errors.New("needs a pointer to a slice")
	}

	if sliceValue.Len() < 1 {
		return nil, nil
	}
	t := sliceValue.Index(0).Type()

	var types []string
	for i := 0; i < t.NumField(); i++ {
		types = append(types, t.Field(i).Name)
	}

	all := make(map[string][]map[string]string)

	for i := 0; i < sliceValue.Len(); i++ {
		value := sliceValue.Index(i)
		stc := make(map[string]string)
		var keyValue string
		for index := range types {
			if str, err := toString(value.FieldByName(types[index])); err == nil {
				stc[types[index]] = str
				if types[index] == field {
					keyValue = str
				}
			} else {
				return nil, errors.New("not support type for field ")
			}
		}
		all[keyValue] = append(all[keyValue], stc)
	}
	return all, nil
}

func toString(field reflect.Value) (string, error) {
	switch field.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), nil
	case reflect.Int:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Int8:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Int16:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Int32:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Uint8:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Uint16:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Uint32:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.String:
		return field.String(), nil
	case reflect.Struct:
		switch field.Type().Name() {
		case "Time":
			if t, ok := field.Interface().(time.Time); ok {
				return t.Format(comm.TimeFormatTime), nil
			}
		}
	case reflect.UnsafePointer:
	case reflect.Invalid:
	default:
	}
	return "", errors.New("not support type for field ")
}
