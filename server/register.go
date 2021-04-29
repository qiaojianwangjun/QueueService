package server

import (
	"QueueService/log"
	"reflect"
	"sync"
)

var HandlerMap map[string]*methodType = make(map[string]*methodType)

type methodType struct {
	sync.Mutex
	Logic    reflect.Value
	Method   reflect.Method
	numCalls uint
}

// Register 注册处理逻辑
func Register(recvHandler interface{}) map[string]*methodType {
	typ := reflect.TypeOf(recvHandler)
	v := reflect.ValueOf(recvHandler)
	method := suitableMethods(v, typ, true)
	return method
}

// suitableMethods 遍历方法，注册到逻辑处理中
func suitableMethods(v reflect.Value, typ reflect.Type, reportErr bool) map[string]*methodType {
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// 必须是三个参数的才处理
		if mtype.NumIn() != 3 {
			if reportErr {
				log.Info("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		HandlerMap[mname] = &methodType{Logic: v, Method: method}
		log.Info("Register:", mname)
	}
	return HandlerMap
}
