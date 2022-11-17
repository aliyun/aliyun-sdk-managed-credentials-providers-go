package utils

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/aliyun/aliyun-sdk-managed-credentials-providers-go/aliyun-sdk-managed-credentials-providers/aliyun-sdk-common-managed-credentials-provider/sdk/constants"
	"github.com/aliyun/aliyun-secretsmanager-client-go/sdk/logger"
)

func Call(instance interface{}, methodName string, params ...interface{}) ([]reflect.Value, error) {
	defer func() {
		if e := recover(); e != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(fmt.Sprintf("%v", e))
		}
	}()
	mtV := reflect.ValueOf(instance)
	parameters := make([]reflect.Value, len(params))
	for i, param := range params {
		parameters[i] = reflect.ValueOf(param)
	}
	method := mtV.MethodByName(methodName)
	if method.IsValid() {
		method.Call(parameters)
	}
	return nil, nil
}

func SetUnExportedField(ptr interface{}, fieldName string, newFieldVal interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(fmt.Sprintf("%v", e))
		}
	}()
	v := reflect.ValueOf(ptr).Elem().FieldByName(fieldName)
	v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	nv := reflect.ValueOf(newFieldVal)
	if v.IsValid() && v.CanSet() {
		v.Set(nv)
		return nil
	} else {
		return errors.New(fmt.Sprintf("fieldName %s is invalid or can not set", fieldName))
	}
}

func GetUnExportedField(ptr interface{}, fieldName string) reflect.Value {
	defer func() {
		if e := recover(); e != nil {
			logger.GetCommonLogger(constants.LoggerName).Errorf(fmt.Sprintf("%v", e))
		}
	}()
	v := reflect.ValueOf(ptr).Elem().FieldByName(fieldName)
	return v
}
