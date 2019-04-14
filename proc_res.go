package gobatis

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

type resultTypeProc = func(rows *sql.Rows, res interface{}) error

var resSetProcMap = map[ResultType]resultTypeProc{
	resultTypeMap:     resMapProc,
	resultTypeMaps:    resMapsProc,
	resultTypeSlice:   resSliceProc,
	resultTypeArray:   resSliceProc,
	resultTypeSlices:  resSlicesProc,
	resultTypeArrays:  resSlicesProc,
	resultTypeValue:   resValueProc,
	resultTypeStructs: resStructsProc,
	resultTypeStruct:  resStructProc,
}

func resStructProc(rows *sql.Rows, res interface{}) error {
	resVal := reflect.ValueOf(res)
	for resVal.Kind() == reflect.Ptr {
		if resVal.IsNil() {
			return ErrStruct
		}
		resVal = resVal.Elem()
	}
	if resVal.Kind() != reflect.Struct {
		return ErrStruct
	}
	arr, err := rowsToStructs(rows, resVal)
	if nil != err {
		return err
	}

	// fixme: 查询结果是返回错误呢, 觉得如果返回错误就会造成错误的困惑,
	//  因为这里的错误定义是用于参数以及异常校验,
	//  如果用户结果校验, 那么如果用户单单用err来判断是否存在查询对象而忽略了其它一些类似sql语句错误, 传参错误等,
	//  还是不处理好呢??? 如果有人看到这里可以提下意见|･ω･｀)
	if len(arr) > 1 {
		//return errors.New("Struct query result more than one row")
		log.Println("[WARN] Struct query result more than one row")
		resVal.Elem().Set(reflect.ValueOf(arr[0]))
	}

	// fixme: 查询结果是返回错误呢, 觉得如果返回错误就会造成错误的困惑,
	//  因为这里的错误定义是用于参数校验以及异常,
	//  如果用户结果校验, 那么如果用户单单用err来判断是否存在查询对象而忽略了其它一些类似sql语句错误, 传参错误等,
	//  还是不处理好呢??? 如果有人看到这里可以提下意见|･ω･｀)
	if len(arr) == 0 {
		//return errors.New("No result")
		log.Println("[WARN] Struct query result is nil")
	}

	if len(arr) == 1 {
		resVal.Set(reflect.ValueOf(arr[0]).Elem())
	}

	return nil
}

func resStructsProc(rows *sql.Rows, res interface{}) error {
	resVal := reflect.ValueOf(res)
	for resVal.Kind() == reflect.Ptr {
		if resVal.IsNil() {
			return ErrStruct
		}
		resVal = resVal.Elem()
	}

	slicePtr := reflect.Indirect(resVal)
	if slicePtr.Kind() != reflect.Slice && slicePtr.Kind() != reflect.Array {
		return errors.New("structs query result must be slice")
	}

	ele := slicePtr.Type().Elem()
	for ele.Kind() == reflect.Ptr {
		ele = ele.Elem()
	}
	result := reflect.New(ele).Elem()
	arr, err := rowsToStructs(rows, result)
	if nil != err {
		return err
	}

	for i := 0; i < len(arr); i++ {
		if ele.Kind() == reflect.Ptr {
			slicePtr.Set(reflect.Append(slicePtr, reflect.ValueOf(arr[i])))
		}else{
			slicePtr.Set(reflect.Append(slicePtr, reflect.Indirect(reflect.ValueOf(arr[i]))))
		}
	}
	return nil
}

func resValueProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("value query result must be ptr")
	}

	arr, err := rowsToSlices(rows)
	if nil != err {
		return err
	}

	if len(arr) > 1 {
		return errors.New("value query result more than one row")
	}

	tempResSlice := arr[0].([]interface{})
	if len(tempResSlice) > 1 {
		return errors.New("value query result more than one col")
	}

	if len(tempResSlice) > 0 {
		if nil != tempResSlice[0] {
			value := reflect.Indirect(resPtr)
			val := dataToFieldVal(tempResSlice[0], value.Type(), "val")
			value.Set(reflect.ValueOf(val))
		}
	}

	return nil
}

func resSlicesProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("slices query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("slices query result must be slice ptr")
	}

	arr, err := rowsToSlices(rows)
	if nil != err {
		return err
	}

	for i := 0; i < len(arr); i++ {
		value.Set(reflect.Append(value, reflect.ValueOf(arr[i])))
	}

	return nil
}

func resSliceProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("slice query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("slice query result must be slice ptr")
	}

	arr, err := rowsToSlices(rows)
	if nil != err {
		return err
	}

	if len(arr) > 1 {
		return errors.New("slice query result more than one row")
	}

	if len(arr) > 0 {
		tempResSlice := arr[0].([]interface{})
		value.Set(reflect.AppendSlice(value, reflect.ValueOf(tempResSlice)))
	}

	return nil
}

func resMapProc(rows *sql.Rows, res interface{}) error {
	resBean := reflect.ValueOf(res)
	if resBean.Kind() == reflect.Ptr {
		return errors.New("map query result can not be ptr")
	}

	if resBean.Kind() != reflect.Map {
		return errors.New("map query result must be map")
	}

	arr, err := rowsToMaps(rows)
	if nil != err {
		return err
	}

	if len(arr) > 1 {
		return errors.New("map query result more than one row")
	}

	if len(arr) > 0 {
		resMap := res.(map[string]interface{})
		tempResMap := arr[0].(map[string]interface{})
		for k, v := range tempResMap {
			resMap[k] = v
		}
	}

	return nil
}

func resMapsProc(rows *sql.Rows, res interface{}) error {
	resPtr := reflect.ValueOf(res)
	if resPtr.Kind() != reflect.Ptr {
		return errors.New("maps query result must be ptr")
	}

	value := reflect.Indirect(resPtr)
	if value.Kind() != reflect.Slice {
		return errors.New("maps query result must be slice ptr")
	}
	arr, err := rowsToMaps(rows)
	if nil != err {
		return err
	}

	for i := 0; i < len(arr); i++ {
		value.Set(reflect.Append(value, reflect.ValueOf(arr[i])))
	}

	return nil
}

func rowsToMaps(rows *sql.Rows) ([]interface{}, error) {
	res := make([]interface{}, 0)
	for rows.Next() {
		resMap := make(map[string]interface{})
		cols, err := rows.Columns()
		if nil != err {
			log.Println(err)
			return res, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		rows.Scan(scanArgs...)
		for i := 0; i < len(cols); i++ {
			val := vals[i]
			if nil != val {
				v := reflect.ValueOf(val)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					val = string(val.([]uint8))
				}
			}
			resMap[cols[i]] = val
		}

		res = append(res, resMap)
	}

	return res, nil
}

func rowsToSlices(rows *sql.Rows) ([]interface{}, error) {
	res := make([]interface{}, 0)
	for rows.Next() {
		resSlice := make([]interface{}, 0)
		cols, err := rows.Columns()
		if nil != err {
			log.Println(err)
			return nil, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		rows.Scan(scanArgs...)
		for i := 0; i < len(cols); i++ {
			val := vals[i]
			if nil != val {
				v := reflect.ValueOf(val)
				if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
					val = string(val.([]uint8))
				}
			}
			resSlice = append(resSlice, val)
		}

		res = append(res, resSlice)
	}

	return res, nil
}

func rowsToStructs(rows *sql.Rows, resVal reflect.Value) ([]interface{}, error) {
	fieldsMapper := make(map[string]string)
	resType := resVal.Type()
	fields := resType.NumField()
	for i := 0; i < fields; i++ {
		field := resType.Field(i)
		fieldsMapper[field.Name] = field.Name
		tag := field.Tag.Get("field")
		if tag != "" {
			fieldsMapper[tag] = field.Name
		}
	}

	res := make([]interface{}, 0)
	for rows.Next() {
		cols, err := rows.Columns()
		if nil != err {
			log.Println("rows.Columns() err:", err)
			return nil, err
		}

		vals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range vals {
			scanArgs[i] = &vals[i]
		}

		err = rows.Scan(scanArgs...)
		if nil != err {
			return nil, err
		}

		obj := reflect.New(resType).Elem()
		objPtr := reflect.Indirect(obj)
		for i := 0; i < len(cols); i++ {
			colName := cols[i]
			fieldName := fieldsMapper[colName]
			field := objPtr.FieldByName(fieldName)
			// 设置相关字段的值,并判断是否可设值
			if field.CanSet() && vals[i] != nil {
				//获取字段类型并设值
				data := dataToFieldVal(vals[i], field.Type(), fieldName)

				// 数据库返回类型与字段类型不符合的情况下通知用户
				if reflect.TypeOf(data).Name() != field.Type().Name() {
					warnInfo := "[WARN] fieldType != dataType, filedName:" + fieldName +
						" fieldType:" + field.Type().Name() +
						" dataType:" + reflect.TypeOf(data).Name()
					log.Println(warnInfo)
				}
				if nil != data {
					_ = valSet(data,field)
				}
			}
		}

		if objPtr.CanInterface() {
			res = append(res, objPtr.Addr().Interface())
		}
	}

	return res, nil
}
