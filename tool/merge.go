package tool

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// Deprecated: 已废弃
func MergeByReflection(from interface{}, to interface{}) error {
	// 检查 from 和 to 是否为指针类型
	if reflect.TypeOf(from).Kind() != reflect.Ptr {
		return errors.New("'from' must be a pointer")
	}
	if reflect.TypeOf(to).Kind() != reflect.Ptr {
		return errors.New("'to' must be a pointer")
	}

	fromValue := reflect.ValueOf(from).Elem()
	toValue := reflect.ValueOf(to).Elem()

	// 检查 from 和 to 是否为结构体类型
	if fromValue.Kind() != reflect.Struct || toValue.Kind() != reflect.Struct {
		return errors.New("'from' and 'to' must be struct")
	}

	fromType := fromValue.Type()

	// 遍历 from 的字段并复制到 to
	for i := 0; i < fromValue.NumField(); i++ {
		fromField := fromValue.Field(i)
		toField := toValue.FieldByName(fromType.Field(i).Name)

		if toField.IsValid() && toField.CanSet() {
			toField.Set(fromField)
		}
	}

	return nil
}

func MergeByJson(from interface{}, to interface{}) {
	fromJson, _ := json.Marshal(from)
	_ = json.Unmarshal(fromJson, to)
}

func MarshalWithoutErr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func Marshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}

func Unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
