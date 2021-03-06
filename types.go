package gocsv

import (
	"fmt"
	"reflect"
	"strconv"
)

// --------------------------------------------------------------------------
// Conversion interfaces

// TypeMarshaller is implemented by any value that has a MarshalCSV method
// This converter is used to convert the value to it string representation
type TypeMarshaller interface {
	MarshalCSV() (string, error)
}

// Stringer is implemented by any value that has a String method
// This converter is used to convert the value to it string representation
// This converter will be used if your value does not implement TypeMarshaller
type Stringer interface {
	String() string
}

// TypeUnmarshaller is implemented by any value that has an UnmarshalCSV method
// This converter is used to convert a string to your value representation of that string
type TypeUnmarshaller interface {
	UnmarshalCSV(string) error
}

var (
	stringerType     = reflect.TypeOf((*Stringer)(nil)).Elem()
	marshallerType   = reflect.TypeOf((*TypeMarshaller)(nil)).Elem()
	unMarshallerType = reflect.TypeOf((*TypeUnmarshaller)(nil)).Elem()
)

// --------------------------------------------------------------------------
// Conversion helpers

func toString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return string(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return string(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), 64, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}

func toBool(in interface{}) (bool, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := inValue.String()
		if s == "true" || s == "yes" || s == "1" {
			return true, nil
		} else if s == "false" || s == "no" || s == "0" {
			return false, nil
		}
	case reflect.Bool:
		return inValue.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := inValue.Int()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := inValue.Uint()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		f := inValue.Float()
		if f != 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to bool")
}

func toInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseInt(inValue.String(), 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func toUint(in interface{}) (uint64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseUint(inValue.String(), 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inValue.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to uint")
}

func toFloat(in interface{}) (float64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return strconv.ParseFloat(inValue.String(), 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		s, err := toString(value)
		if err != nil {
			return err
		}
		field.SetString(s)
	case reflect.Bool:
		b, err := toBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := toInt(value)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := toUint(value)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case reflect.Float32, reflect.Float64:
		f, err := toFloat(value)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		return unmarshall(field, value)
	}
	return nil
}

func getFieldAsString(field reflect.Value) (str string, err error) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Bool:
		str, err = toString(field.Bool())
		if err != nil {
			return str, err
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		str, err = toString(field.Int())
		if err != nil {
			return str, err
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str, err = toString(field.Uint())
		if err != nil {
			return str, err
		}
	case reflect.Float32, reflect.Float64:
		str, err = toString(field.Float())
		if err != nil {
			return str, err
		}
	default:
		return marshall(field)
	}
	return str, nil
}

// --------------------------------------------------------------------------
// Un/serializations helpers

func unmarshall(field reflect.Value, value string) error {
	dupField := field
	unMarshallIt := func(finalField reflect.Value) error {
		if finalField.CanInterface() && finalField.Type().Implements(unMarshallerType) {
			if err := finalField.Interface().(TypeUnmarshaller).UnmarshalCSV(value); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("No known conversion from string to " + field.Type().String() + ", " + field.Type().String() + " does not implements TypeUnmarshaller")
	}
	for dupField.Kind() == reflect.Interface || dupField.Kind() == reflect.Ptr {
		if dupField.IsNil() {
			dupField = reflect.New(field.Type().Elem())
			field.Set(dupField)
			return unMarshallIt(dupField)
			break
		}
		dupField = dupField.Elem()
	}
	if dupField.CanAddr() {
		return unMarshallIt(dupField.Addr())
	}
	return fmt.Errorf("No known conversion from string to " + field.Type().String() + ", " + field.Type().String() + " does not implements TypeUnmarshaller")
}

func marshall(field reflect.Value) (value string, err error) {
	dupField := field
	marshallIt := func(finalField reflect.Value) (string, error) {
		if finalField.CanInterface() && finalField.Type().Implements(marshallerType) { // Use TypeMarshaller when possible
			return finalField.Interface().(TypeMarshaller).MarshalCSV()
		} else if finalField.CanInterface() && finalField.Type().Implements(stringerType) { // Otherwise try to use Stringer
			return finalField.Interface().(Stringer).String(), nil
		}
		return value, fmt.Errorf("No known conversion from " + field.Type().String() + " to string, " + field.Type().String() + " does not implements TypeMarshaller nor Stringer")
	}
	for dupField.Kind() == reflect.Interface || dupField.Kind() == reflect.Ptr {
		if dupField.IsNil() {
			return value, nil
		}
		dupField = dupField.Elem()
	}
	if dupField.CanAddr() {
		return marshallIt(dupField.Addr())
	}
	return value, fmt.Errorf("No known conversion from " + field.Type().String() + " to string, " + field.Type().String() + " does not implements TypeMarshaller nor Stringer")
}
