package gobatis

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
)



func stringToVal(data interface{}, tp reflect.Type) interface{} {
	str := data.(string)
	switch tp.Kind() {
	case reflect.Bool:
		data = false
		if str == "1" {
			data = true
		}
	case reflect.Int:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int(i)
	case reflect.Int8:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int8(i)
	case reflect.Int16:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int16(i)
	case reflect.Int32:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Int64:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int64(i)
	case reflect.Uint:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Uint8:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint8(ui)
	case reflect.Uint16:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint16(ui)
	case reflect.Uint32:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint32(ui)
	case reflect.Uint64:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint64(ui)
	case reflect.Uintptr:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uintptr(ui)
	case reflect.Float32:
		f64, _ := strconv.ParseFloat(str, 64)
		data = float32(f64)
	case reflect.Float64:
		f64, _ := strconv.ParseFloat(str, 64)
		data = f64
	case reflect.Complex64:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex64
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	case reflect.Complex128:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex128
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	}

	return data
}

func valSet(val interface{}, field reflect.Value) error{
	switch field.Kind() {
	case reflect.Bool:
		switch val.(type) {
		case bool:
			field.SetBool(val.(bool))
		case string:
			b, e := strconv.ParseBool(val.(string))
			if nil == e {
				field.SetBool(b)
			}
		case int, int8, int16, int32, int64:
			if dataVal(val).(int64) == 1 {
				field.SetBool(true)
			}
		case uint, uint8,uint16,uint32,uint64: {
			if dataVal(val).(uint64) == 1 {
				field.SetBool(true)
			}

		}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if isSameType(reflect.ValueOf(val), field) {
			field.SetInt(dataVal(val).(int64))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if isSameType(reflect.ValueOf(val), field) {
			field.SetUint(dataVal(val).(uint64))
		}
	case reflect.Float32, reflect.Float64:
		if isSameType(reflect.ValueOf(val), field) {
			field.SetFloat(dataVal(val).(float64))
		}
	case reflect.Complex64,reflect.Complex128:
		if isSameType(reflect.ValueOf(val), field) {
			field.SetComplex(dataVal(val).(complex128))
		}
	default:
		if isSameType(reflect.ValueOf(val), field) {
			field.SetString(val.(string))
		}
	}
	return nil
}

const (
	intType = iota
	uintType
	floatType
	complexType
	otherType
)

func isSameType(v1 reflect.Value, v2 reflect.Value) bool {
	i := valType(v1)
	i2 := valType(v2)
	if i != otherType && i2 != otherType {
		return i == i2
	}
    return v1.Kind() == v2.Kind()
}


func valType(val reflect.Value ) int {
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intType
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintType
	case reflect.Float32, reflect.Float64:
		return floatType
	case reflect.Complex64,reflect.Complex128:
		return complexType
	}
	return otherType
}


func dataVal(data interface{}) interface{}{
	t := reflect.ValueOf(data)

	switch t.Kind() {
	case reflect.Int:
		return int64(data.(int))
	case reflect.Int8:
		return int64(data.(int8))
	case reflect.Int16:
		return int64(data.(int16))
	case reflect.Int32:
		return int64(data.(int32))
	case reflect.Int64:
		return data.(int64)
	case reflect.Uint:
		return uint64(data.(uint))
	case reflect.Uint8:
		return uint64(data.(uint8))
	case reflect.Uint16:
		return uint64(data.(uint16))
	case reflect.Uint32:
		return uint64(data.(uint32))
	case reflect.Uint64:
		return data.(uint64)
	case reflect.Uintptr:
		return uint64(data.(uintptr))
	case reflect.Float32:
		return float64(data.(float32))
	case reflect.Float64:
		return data.(float64)
	case reflect.Complex64:
		return complex128(data.(complex64))
	case reflect.Complex128:
		return data.(complex128)
	}
	return data
}


func bytesToVal(data interface{}, tp reflect.Type) interface{} {
	str := string(data.([]uint8))
	switch tp.Kind() {
	case reflect.Bool:
		data = false
		if str == "1" {
			data = true
		}
	case reflect.Int:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int(i)
	case reflect.Int8:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int8(i)
	case reflect.Int16:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int16(i)
	case reflect.Int32:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Int64:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int64(i)
	case reflect.Uint:
		i, _ := strconv.ParseInt(str, 10, 64)
		data = int32(i)
	case reflect.Uint8:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint8(ui)
	case reflect.Uint16:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint16(ui)
	case reflect.Uint32:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint32(ui)
	case reflect.Uint64:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uint64(ui)
	case reflect.Uintptr:
		ui, _ := strconv.ParseUint(str, 0, 64)
		data = uintptr(ui)
	case reflect.Float32:
		f64, _ := strconv.ParseFloat(str, 64)
		data = float32(f64)
	case reflect.Float64:
		f64, _ := strconv.ParseFloat(str, 64)
		data = f64
	case reflect.Complex64:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex64
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	case reflect.Complex128:
		binBuf := bytes.NewBuffer(data.([]uint8))
		var x complex128
		_ = binary.Read(binBuf, binary.BigEndian, &x)
		data = x
	}

	return data
}

func valToString(data interface{}) string {
	tp := reflect.TypeOf(data)
	s := ""
	switch tp.Kind() {
	case reflect.Bool:
		s = strconv.FormatBool(data.(bool))
	case reflect.Int:
		s = strconv.FormatInt(int64(data.(int)), 10)
	case reflect.Int8:
		s = strconv.FormatInt(int64(data.(int8)), 10)
	case reflect.Int16:
		s = strconv.FormatInt(int64(data.(int16)), 10)
	case reflect.Int32:
		s = strconv.FormatInt(int64(data.(int32)), 10)
	case reflect.Int64:
		s = strconv.FormatInt(int64(data.(int64)), 10)
	case reflect.Uint:
		s = strconv.FormatUint(uint64(data.(uint)), 10)
	case reflect.Uint8:
		s = strconv.FormatUint(uint64(data.(uint8)), 10)
	case reflect.Uint16:
		s = strconv.FormatUint(uint64(data.(uint16)), 10)
	case reflect.Uint32:
		s = strconv.FormatUint(uint64(data.(uint32)), 10)
	case reflect.Uint64:
		s = strconv.FormatUint(uint64(data.(uint64)), 10)
	case reflect.Uintptr:
		s = fmt.Sprint(data.(uintptr))
	case reflect.Float32:
		s = strconv.FormatFloat(float64(data.(float32)), 'f', -1, 64)
	case reflect.Float64:
		s = strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case reflect.Complex64:
		s = fmt.Sprint(data.(complex64))
	case reflect.Complex128:
		s = fmt.Sprint(data.(complex128))
	default:
		log.Println("[WARN]no process for type:" + tp.Name())
	}
	return s
}

func dataToFieldVal(data interface{}, tp reflect.Type, fieldName string) interface{} {
	defer func() {
		if err := recover(); nil != err {
			log.Println("[WARN] data to field val panic, fieldName:", fieldName, " err:", err)
		}
	}()

	typeName := tp.Name()
	switch {
	case typeName == "bool" ||
		typeName == "int" ||
		typeName == "int8" ||
		typeName == "int16" ||
		typeName == "int32" ||
		typeName == "int64" ||
		typeName == "uint" ||
		typeName == "uint8" ||
		typeName == "uint16" ||
		typeName == "uint32" ||
		typeName == "uint64" ||
		typeName == "uintptr" ||
		typeName == "float32" ||
		typeName == "float64" ||
		typeName == "complex64" ||
		typeName == "complex128":
		if nil != data {
			dataTp := reflect.TypeOf(data)
			if dataTp.Kind() == reflect.Slice ||
				dataTp.Kind() == reflect.Array {
				data = bytesToVal(data, tp)
			}

			dataTp = reflect.TypeOf(data)
			if dataTp.Kind() == reflect.String {
				data = stringToVal(data, tp)
			}

			return data
		}
	case typeName == "string":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				return string(data.([]byte))
			}

			data = valToString(data)
			return string(data.(string))
		}
	case typeName == "Time":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			tm, err := time.Parse("2006-01-02 15:04:05", data.(string))
			if err != nil {
				panic("time.Parse err:" + err.Error())
			}
			return tm
		}
	case typeName == "NullString":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}
			return NullString{String: data.(string), Valid: true}
		}
	case typeName == "NullInt64":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			i, err := strconv.ParseInt(data.(string), 10, 64)
			if err != nil {
				panic("ParseInt err:" + err.Error())
			}
			return NullInt64{Int64: i, Valid: true}
		}
	case typeName == "NullBool":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}
			if data.(string) == "true" {
				return NullBool{Bool: true, Valid: true}
			}
			return NullBool{Bool: false, Valid: true}
		}
	case typeName == "NullFloat64":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			f64, err := strconv.ParseFloat(data.(string), 64)
			if err != nil {
				panic("ParseFloat err:" + err.Error())
			}

			return NullFloat64{Float64: f64, Valid: true}
		}
	case typeName == "NullTime":
		if nil != data {
			if reflect.TypeOf(data).Kind() == reflect.Slice ||
				reflect.TypeOf(data).Kind() == reflect.Array {
				data = string(data.([]byte))
			} else {
				data = valToString(data)
			}

			tm, err := time.Parse("2006-01-02 15:04:05", data.(string))
			if err != nil {
				panic("time.Parse err:" + err.Error())
			}

			return NullTime{Time: tm, Valid: true}
		}
	}

	return nil
}
